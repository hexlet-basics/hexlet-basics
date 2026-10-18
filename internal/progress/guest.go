package progress

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/events"
)

// GuestCookieName carries a visitor's progress. It is httpOnly: the client
// renders what the server derives from it and never reads it, so there is one
// place that interprets this state.
const GuestCookieName = "guest_progress"

// guestCookieMaxBytes keeps the encoded value clear of the 4096-byte browser
// limit. A visitor who exceeds it loses their least recently touched courses
// rather than the whole cookie — silently dropping everything would be the one
// failure mode a learner would notice.
const guestCookieMaxBytes = 3500

// guestCookieTTL matches the legacy guest window.
const guestCookieTTL = 365 * 24 * time.Hour

// ErrGuestCookieInvalid reports a cookie whose signature does not verify, or
// which is not decodable at all. Both are treated the same way: the visitor is
// handled as if they had no cookie.
var ErrGuestCookieInvalid = errors.New("progress: guest cookie is not valid")

// GuestProgress is a visitor's whole state: per Course, the slugs of the
// Lessons they finished.
//
// Slugs rather than positions, because a position is meaningless across
// Versions — promoting a Version that inserts or removes Lessons renumbers
// everything, and a stored number would come to denote a different Lesson. A
// slug is resolved against the current Version at read time, and one it no
// longer contains simply counts for nothing.
//
// Lessons may be taken in any order (ADR-0016), so the finished set is stored as
// a set: a single "furthest" Lesson cannot say that the third was passed and the
// first was not. Entries are ordered least-recently-touched first, which is the
// order they are evicted in.
type GuestProgress struct {
	Entries []GuestEntry `json:"entries"`
}

// GuestEntry is one Course's finished Lessons.
type GuestEntry struct {
	CourseSlug string   `json:"course"`
	Lessons    []string `json:"lessons"`
}

// Has reports whether the visitor finished a Lesson of a Course.
func (g GuestProgress) Has(courseSlug, lessonSlug string) bool {
	for _, entry := range g.Entries {
		if entry.CourseSlug == courseSlug {
			return slices.Contains(entry.Lessons, lessonSlug)
		}
	}
	return false
}

// Record adds a finished Lesson to its Course's set and marks the Course most
// recently touched. Adding one already in the set changes nothing but the
// Course's recency.
func (g GuestProgress) Record(courseSlug, lessonSlug string) GuestProgress {
	var lessons []string
	kept := make([]GuestEntry, 0, len(g.Entries)+1)
	for _, entry := range g.Entries {
		if entry.CourseSlug == courseSlug {
			lessons = entry.Lessons
			continue
		}
		kept = append(kept, entry)
	}
	if !slices.Contains(lessons, lessonSlug) {
		lessons = append(slices.Clone(lessons), lessonSlug)
	}
	return GuestProgress{Entries: append(kept, GuestEntry{CourseSlug: courseSlug, Lessons: lessons})}
}

// GuestCodec signs and verifies the cookie with the authentication secret.
//
// Signing is not optional: this state becomes database rows and domain events
// when the visitor signs up, so an unsigned cookie would let anyone mint
// fabricated completions on their own account.
type GuestCodec struct {
	secret []byte
}

// NewGuestCodec builds the codec around the existing authentication secret.
func NewGuestCodec(secret string) *GuestCodec {
	return &GuestCodec{secret: []byte(secret)}
}

// Encode renders the progress as a signed cookie value, evicting the least
// recently touched Courses until it fits.
func (c *GuestCodec) Encode(progress GuestProgress) (string, error) {
	for {
		value, err := c.sign(progress)
		if err != nil {
			return "", err
		}
		if len(value) <= guestCookieMaxBytes || len(progress.Entries) == 0 {
			return value, nil
		}
		progress = GuestProgress{Entries: progress.Entries[1:]}
	}
}

// Decode verifies the signature and returns the progress it carries.
func (c *GuestCodec) Decode(raw string) (GuestProgress, error) {
	payload, signature, found := strings.Cut(raw, ".")
	if !found {
		return GuestProgress{}, ErrGuestCookieInvalid
	}

	body, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return GuestProgress{}, ErrGuestCookieInvalid
	}
	got, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return GuestProgress{}, ErrGuestCookieInvalid
	}
	if !hmac.Equal(got, c.mac(body)) {
		return GuestProgress{}, ErrGuestCookieInvalid
	}

	var progress GuestProgress
	if err := json.Unmarshal(body, &progress); err != nil {
		return GuestProgress{}, ErrGuestCookieInvalid
	}
	return progress, nil
}

// Cookie is the response cookie carrying the progress: httpOnly so only the
// server interprets it, Lax so it survives the navigation back from an OAuth
// provider, and a year long.
func (c *GuestCodec) Cookie(progress GuestProgress, secure bool) (*http.Cookie, error) {
	value, err := c.Encode(progress)
	if err != nil {
		return nil, err
	}
	//nolint:gosec // G124: HttpOnly and SameSite are set; Secure is the caller's flag, false only for plain-HTTP local development
	return &http.Cookie{
		Name:     GuestCookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   int(guestCookieTTL.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}, nil
}

// ClearedCookie is the response cookie that ends the visitor's progress cookie,
// with the same attributes Cookie sets so the browser matches it to the one it
// holds. It is sent once the progress has been credited to an account.
func (c *GuestCodec) ClearedCookie(secure bool) *http.Cookie {
	//nolint:gosec // G124: same attributes as Cookie, Secure is the caller's flag
	return &http.Cookie{
		Name:     GuestCookieName,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func (c *GuestCodec) sign(progress GuestProgress) (string, error) {
	body, err := json.Marshal(progress)
	if err != nil {
		return "", errs.Wrapf(err, "encode guest progress")
	}
	return base64.RawURLEncoding.EncodeToString(body) + "." +
		base64.RawURLEncoding.EncodeToString(c.mac(body)), nil
}

func (c *GuestCodec) mac(body []byte) []byte {
	sum := hmac.New(sha256.New, c.secret)
	sum.Write(body)
	return sum.Sum(nil)
}

// MergeGuest credits a visitor's cookie progress to their account.
//
// The transferred set is exactly the Lessons the cookie names, resolved against
// the current Version; the ones the account already finished are left as they
// are, so the merge is a union and never takes anything away. A stored slug the
// current Version no longer contains resolves to nothing and is skipped — a
// bounded loss for someone without an account, and one that cannot corrupt the
// account.
func (p *Progress) MergeGuest(ctx context.Context, userID int, guest GuestProgress, locale string) error {
	if len(guest.Entries) == 0 {
		return nil
	}

	return p.store.WithinTx(ctx, func(tx *sql.Tx, db *ent.Client) error {
		for _, entry := range guest.Entries {
			if err := p.mergeCourse(ctx, tx, db, userID, entry, locale); err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *Progress) mergeCourse(ctx context.Context, tx *sql.Tx, db *ent.Client, userID int, entry GuestEntry, locale string) error {
	crs, err := db.Course.Query().Where(course.SlugEQ(entry.CourseSlug)).Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil
	case err != nil:
		return errs.With("course_slug", entry.CourseSlug).Wrapf(err, "load course")
	}

	lessons, err := currentLessons(ctx, db, crs)
	if err != nil {
		return err
	}

	finished, err := finishedLessonIDs(ctx, db, userID, crs.ID)
	if err != nil {
		return err
	}

	// In course order, which is the order the facts are published in.
	credited := lo.Filter(lessons, func(lesson currentLesson, _ int) bool {
		return slices.Contains(entry.Lessons, lesson.slug) && !finished[lesson.lessonID]
	})
	if len(credited) == 0 {
		// Nothing the account does not already have — or nothing the current
		// Version still contains, in which case this Course resets for the guest.
		return nil
	}

	enrolled, createdEnrollment, err := p.enroll(ctx, db, userID, crs.ID)
	if err != nil {
		return err
	}

	occurredAt := p.now()
	if createdEnrollment {
		count, err := db.Enrollment.Query().
			Where(enrollment.UserID(userID), enrollment.StateEQ(StateStarted)).
			Count(ctx)
		if err != nil {
			return errs.Wrapf(err, "count started enrollments")
		}
		if err := p.publisher.Publish(ctx, tx, events.CourseStarted{
			OccurrenceCount: count,
			Slug:            entry.CourseSlug,
			Locale:          locale,
			OccurredAt:      occurredAt,
		}); err != nil {
			return errs.Wrapf(err, "publish course started")
		}
	}

	return p.creditLessons(ctx, tx, db, userID, crs, enrolled, credited, locale, occurredAt)
}

// creditLessons writes the completions of a guest's finished Lessons and the pair
// of facts each one would have produced had a learner passed it.
//
// The rows go in three statements however many Lessons there are: one read of what
// the account already holds, one update finishing the started ones, and one
// insert of the missing ones. The occurrence count each fact carries is the
// number of Lesson Progress rows on the Enrollment as of that Lesson, so it is
// counted once before the writes and advanced by every row this merge creates,
// in Position order, as the per-Lesson writes used to leave it.
func (p *Progress) creditLessons(
	ctx context.Context,
	tx *sql.Tx,
	db *ent.Client,
	userID int,
	crs *ent.Course,
	enrolled *ent.Enrollment,
	lessons []currentLesson,
	locale string,
	occurredAt time.Time,
) error {
	if len(lessons) == 0 {
		return nil
	}

	lessonIDs := lo.Map(lessons, func(l currentLesson, _ int) int { return l.lessonID })
	held, err := db.LessonProgress.Query().
		Where(lessonprogress.UserID(userID), lessonprogress.LessonIDIn(lessonIDs...)).
		Select(lessonprogress.FieldLessonID).
		All(ctx)
	if err != nil {
		return errs.Wrapf(err, "load lesson progress")
	}
	isHeld := lo.SliceToMap(held, func(row *ent.LessonProgress) (int, bool) { return row.LessonID, true })

	count, err := lessonProgressCount(ctx, db, enrolled.ID)
	if err != nil {
		return err
	}

	// Started but not finished: the guest's cookie says otherwise.
	if len(held) > 0 {
		if _, err := db.LessonProgress.Update().
			Where(lessonprogress.UserID(userID), lessonprogress.LessonIDIn(lo.Keys(isHeld)...)).
			SetState(StateFinished).
			Save(ctx); err != nil {
			return errs.Wrapf(err, "finish lesson progress")
		}
	}

	missing := lo.Filter(lessons, func(l currentLesson, _ int) bool { return !isHeld[l.lessonID] })
	if len(missing) > 0 {
		builders := lo.Map(missing, func(l currentLesson, _ int) *ent.LessonProgressCreate {
			return db.LessonProgress.Create().
				SetUserID(userID).
				SetCourseID(crs.ID).
				SetEnrollmentID(enrolled.ID).
				SetLessonID(l.lessonID).
				SetState(StateFinished)
		})
		if _, err := db.LessonProgress.CreateBulk(builders...).Save(ctx); err != nil {
			return errs.Wrapf(err, "create finished lesson progress")
		}
	}

	for _, lesson := range lessons {
		if !isHeld[lesson.lessonID] {
			count++
		}
		for _, event := range []events.Event{
			events.LessonStarted{
				OccurrenceCount: count,
				LessonSlug:      lesson.slug,
				CourseSlug:      lo.FromPtr(crs.Slug),
				Locale:          locale,
				OccurredAt:      occurredAt,
			},
			events.LessonFinished{
				OccurrenceCount: count,
				LessonSlug:      lesson.slug,
				CourseSlug:      lo.FromPtr(crs.Slug),
				Locale:          locale,
				OccurredAt:      occurredAt,
			},
		} {
			if err := p.publisher.Publish(ctx, tx, event); err != nil {
				return errs.Wrapf(err, "publish merge fact")
			}
		}
	}
	return nil
}
