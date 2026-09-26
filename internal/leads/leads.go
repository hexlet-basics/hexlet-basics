// Package leads owns the site's lead capture: storing a contact request and
// raising the LeadCreated fact the amoCRM integration consumes.
package leads

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/ent/tag"
	"hexletbasics/ent/tagging"
	"hexletbasics/internal/events"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/store"
)

// Contact methods a lead can leave. Each names the `leads` column the contact
// value is written to, exactly as legacy's `write_attribute(contact_method)`.
const (
	ContactTelegram = "telegram"
	ContactPhone    = "phone"
	ContactWhatsApp = "whatsapp"
)

// FirstVisit is where the visitor first came from, as the browser recorded it.
// Legacy took these from ahoy's visit; ahoy is not ported (ADR-0015).
type FirstVisit struct {
	UTMSource   *string
	UTMMedium   *string
	UTMCampaign *string
	UTMContent  *string
	UTMTerm     *string
	LandingPage *string
	Referrer    *string
}

// Submission is one validated lead form.
type Submission struct {
	UserID        int
	ContactMethod string
	ContactValue  string
	YMClientID    *string
	FirstVisit    FirstVisit
	// IP is the client address of the submitting request, nil when unknown.
	IP *string
}

// Creator is the seam the handler depends on.
type Creator interface {
	Create(ctx context.Context, submission Submission) (*ent.Lead, error)
}

// Recorder atomically stores a lead and its LeadCreated outbox record, so a
// lead never exists without reaching amoCRM and amoCRM never hears of a lead
// that was rolled back.
type Recorder struct {
	store     store.Transactor
	publisher events.TxPublisher
	now       func() time.Time
}

// NewRecorder builds the production lead recorder.
func NewRecorder(txStore store.Transactor, publisher events.TxPublisher) *Recorder {
	return &Recorder{store: txStore, publisher: publisher, now: time.Now}
}

// courseData is one row of the legacy `courses_data` snapshot.
type courseData struct {
	Slug                 *string `json:"slug"`
	LessonsFinishedCount int     `json:"lessons_finished_count"`
}

// emptySurveyAnswers is what legacy wrote on this path: surveys are not ported
// (ADR-0015), so there are never answers to attach.
const emptySurveyAnswers = "[]"

// Create stores the lead and publishes LeadCreated in the same transaction.
//
// `courses_data` and `survey_answers_data` are written as JSON, where legacy's
// YAML column coder wrote YAML. JSON is valid YAML, so legacy still loads these
// rows during the rollback window (with string keys, which only the never-sent
// amoCRM note would have read by symbol), and the admin list shows them as is.
func (r *Recorder) Create(ctx context.Context, submission Submission) (*ent.Lead, error) {
	var created *ent.Lead
	err := r.store.WithinTx(ctx, func(tx *sql.Tx, db *ent.Client) error {
		u, err := db.User.Get(ctx, submission.UserID)
		if err != nil {
			return fmt.Errorf("load lead user %d: %w", submission.UserID, err)
		}

		courses, err := coursesData(ctx, db, submission.UserID)
		if err != nil {
			return err
		}

		create := db.Lead.Create().
			SetUserID(u.ID).
			SetNillableYmClientID(submission.YMClientID).
			SetCoursesData(courses).
			SetSurveyAnswersData(emptySurveyAnswers)
		switch submission.ContactMethod {
		case ContactTelegram:
			create.SetTelegram(submission.ContactValue)
		case ContactPhone:
			create.SetPhone(submission.ContactValue)
		case ContactWhatsApp:
			create.SetWhatsapp(submission.ContactValue)
		default:
			// The contract's enum rejects anything else before it gets here.
			return fmt.Errorf("unknown contact method %q", submission.ContactMethod)
		}
		created, err = create.Save(ctx)
		if err != nil {
			return fmt.Errorf("create lead: %w", err)
		}

		if err := removeShouldBeLead(ctx, db, u.ID); err != nil {
			return err
		}

		visit := submission.FirstVisit
		if err := r.publisher.Publish(ctx, tx, events.LeadCreated{
			LeadID:      created.ID,
			UserID:      u.ID,
			UserName:    userName(u),
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			YMClientID:  created.YmClientID,
			UTMSource:   visit.UTMSource,
			UTMMedium:   visit.UTMMedium,
			UTMCampaign: visit.UTMCampaign,
			UTMTerm:     visit.UTMTerm,
			UTMContent:  visit.UTMContent,
			Email:       u.Email,
			Phone:       created.Phone,
			Telegram:    created.Telegram,
			WhatsApp:    created.Whatsapp,
			LandingPage: visit.LandingPage,
			Referrer:    visit.Referrer,
			IP:          submission.IP,
			OccurredAt:  r.now().UTC(),
		}); err != nil {
			return fmt.Errorf("publish lead created: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

// shouldBeLeadTag is the acts_as_taggable_on tag that marks a user the sales
// team wants a contact from; legacy's layout asked such a user for a contact
// method, and a submitted lead answers that ask.
const shouldBeLeadTag = "should_be_lead"

// removeShouldBeLead is legacy's `user.tag_list.remove("should_be_lead")` +
// `save!`: the user's tagging goes, and the tag's counter cache drops with it
// (acts_as_taggable_on keeps an unused tag row, as `remove_unused_tags` is
// off). It runs in the lead's transaction, so a rolled-back lead keeps the tag.
// A user who never had the tag, or a database without it, is no change.
func removeShouldBeLead(ctx context.Context, db *ent.Client, userID int) error {
	t, err := db.Tag.Query().Where(tag.Name(shouldBeLeadTag)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("load %s tag: %w", shouldBeLeadTag, err)
	}

	removed, err := db.Tagging.Delete().
		Where(
			tagging.TagID(t.ID),
			tagging.TaggableType(userTaggableType),
			tagging.TaggableID(userID),
			tagging.Context(tagListContext),
		).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("remove %s from user %d: %w", shouldBeLeadTag, userID, err)
	}
	if removed == 0 {
		return nil
	}
	if err := db.Tag.UpdateOneID(t.ID).AddTaggingsCount(-removed).Exec(ctx); err != nil {
		return fmt.Errorf("count %s taggings: %w", shouldBeLeadTag, err)
	}
	return nil
}

// The polymorphic keys acts_as_taggable_on writes for a User's `tag_list`.
const (
	userTaggableType = "User"
	tagListContext   = "tags"
)

// coursesData snapshots every Course the user is enrolled in, in any state,
// with how many of its Lessons they finished — what the sales team reads to
// know who they are calling.
func coursesData(ctx context.Context, db *ent.Client, userID int) (string, error) {
	enrollments, err := db.Enrollment.Query().
		Where(enrollment.UserID(userID)).
		WithCourse().
		WithLessonProgress(func(q *ent.LessonProgressQuery) {
			q.Where(lessonprogress.StateEQ(progress.StateFinished))
		}).
		Order(ent.Asc(enrollment.FieldID)).
		All(ctx)
	if err != nil {
		return "", fmt.Errorf("load lead enrollments: %w", err)
	}

	data := lo.Map(enrollments, func(e *ent.Enrollment, _ int) courseData {
		var slug *string
		if e.Edges.Course != nil {
			slug = e.Edges.Course.Slug
		}
		return courseData{
			Slug:                 slug,
			LessonsFinishedCount: len(e.Edges.LessonProgress),
		}
	})
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("encode courses data: %w", err)
	}
	return string(encoded), nil
}

// userName mirrors legacy `User#to_s`, down to its "first last" spacing. Phone
// sign-in never worked in production (ADR-0015), so the phone step is moot.
func userName(u *ent.User) string {
	if lo.FromPtr(u.FirstName) != "" || lo.FromPtr(u.LastName) != "" {
		return lo.FromPtr(u.FirstName) + " " + lo.FromPtr(u.LastName)
	}
	if email := lo.FromPtr(u.Email); email != "" {
		return email
	}
	return fmt.Sprintf("User #%d", u.ID)
}
