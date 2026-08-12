package progress

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

// GuestProgress is a visitor's whole state: per Course, the slug of the
// furthest Lesson they finished.
//
// A slug rather than a position, because a position is meaningless across
// Versions — promoting a Version that inserts or removes Lessons renumbers
// everything, and a stored number would come to denote a different Lesson. The
// Position is resolved against the current Version at read time.
//
// One entry per Course is exact rather than lossy: sequential progression makes
// a guest's gaps unrepresentable, so the furthest Lesson implies the rest.
// Entries are ordered least-recently-touched first, which is the order they are
// evicted in.
type GuestProgress struct {
	Entries []GuestEntry `json:"entries"`
}

// GuestEntry is one Course's furthest finished Lesson.
type GuestEntry struct {
	CourseSlug string `json:"course"`
	LessonSlug string `json:"lesson"`
}

// Furthest reports the furthest Lesson the visitor finished in a Course.
func (g GuestProgress) Furthest(courseSlug string) (string, bool) {
	for _, entry := range g.Entries {
		if entry.CourseSlug == courseSlug {
			return entry.LessonSlug, true
		}
	}
	return "", false
}

// Record stores a Course's furthest finished Lesson and marks it most recently
// touched. The caller decides what "furthest" means — this type stores, it does
// not compare positions, because only the current Version can order them.
func (g GuestProgress) Record(courseSlug, lessonSlug string) GuestProgress {
	kept := make([]GuestEntry, 0, len(g.Entries)+1)
	for _, entry := range g.Entries {
		if entry.CourseSlug != courseSlug {
			kept = append(kept, entry)
		}
	}
	return GuestProgress{Entries: append(kept, GuestEntry{CourseSlug: courseSlug, LessonSlug: lessonSlug})}
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

func (c *GuestCodec) sign(progress GuestProgress) (string, error) {
	body, err := json.Marshal(progress)
	if err != nil {
		return "", fmt.Errorf("encode guest progress: %w", err)
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
// The transferred set is the PREFIX up to the guest's furthest finished Lesson,
// not an explicit list: sequential progression means the prefix *is* the set,
// because a guest could not have finished the fifth Lesson without finishing
// the fourth. This is the one place rows are written for Lessons whose
// completion is implied by the rule rather than individually observed, and a
// deliberate departure from the legacy merge, which carried an explicit list of
// lesson ids and could therefore transfer a non-contiguous selection.
//
// Both positions are resolved against the current Version before comparing, so
// an account already further along is left untouched. A stored slug that the
// current Version no longer contains resolves to nothing and that Course is
// skipped — the guest's position resets, which is a bounded loss for someone
// without an account and one that cannot corrupt the account.
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
		return fmt.Errorf("load course %q: %w", entry.CourseSlug, err)
	}

	lessons, err := currentLessons(ctx, db, crs)
	if err != nil {
		return err
	}

	guestPosition := positionOfSlug(lessons, entry.LessonSlug)
	if guestPosition == 0 {
		// The stored Lesson is not in the current Version: there is nothing to
		// resolve the position against, so this Course resets to the beginning.
		return nil
	}

	accountPosition, err := furthestFinishedPosition(ctx, db, userID, crs.ID, positionsOf(lessons))
	if err != nil {
		return err
	}
	if accountPosition >= guestPosition {
		// The account is already further along; the merge takes the higher of
		// the two positions, which is the one already stored.
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
			return fmt.Errorf("count started enrollments: %w", err)
		}
		if err := p.publisher.Publish(ctx, tx, events.CourseStarted{
			OccurrenceCount: count,
			Slug:            entry.CourseSlug,
			Locale:          locale,
			OccurredAt:      occurredAt,
		}); err != nil {
			return fmt.Errorf("publish course started: %w", err)
		}
	}

	finished, err := finishedLessonIDs(ctx, db, userID, crs.ID)
	if err != nil {
		return err
	}

	for _, lesson := range lessons {
		if lesson.position > guestPosition || finished[lesson.lessonID] {
			continue
		}
		if err := p.creditLesson(ctx, tx, db, userID, crs, enrolled, lesson, locale, occurredAt); err != nil {
			return err
		}
	}
	return nil
}

// creditLesson writes one implied completion and the pair of facts that a
// learner passing it would have produced.
func (p *Progress) creditLesson(
	ctx context.Context,
	tx *sql.Tx,
	db *ent.Client,
	userID int,
	crs *ent.Course,
	enrolled *ent.Enrollment,
	lesson currentLesson,
	locale string,
	occurredAt time.Time,
) error {
	existing, err := db.LessonProgress.Query().
		Where(lessonprogress.UserID(userID), lessonprogress.LessonID(lesson.lessonID)).
		Only(ctx)
	switch {
	case err == nil:
		// Started but not finished: the guest's prefix says otherwise.
		if _, err := existing.Update().SetState(StateFinished).Save(ctx); err != nil {
			return fmt.Errorf("finish lesson progress: %w", err)
		}
	case ent.IsNotFound(err):
		if _, err := db.LessonProgress.Create().
			SetUserID(userID).
			SetCourseID(crs.ID).
			SetEnrollmentID(enrolled.ID).
			SetLessonID(lesson.lessonID).
			SetState(StateFinished).
			Save(ctx); err != nil {
			return fmt.Errorf("create finished lesson progress: %w", err)
		}
	default:
		return fmt.Errorf("load lesson progress: %w", err)
	}

	count, err := lessonProgressCount(ctx, db, enrolled.ID)
	if err != nil {
		return err
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
			return fmt.Errorf("publish merge fact: %w", err)
		}
	}
	return nil
}
