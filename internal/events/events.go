// Package events owns domain-event contracts and their Watermill transport.
//
// Domain callers publish immutable facts through the transactional publisher;
// Watermill owns routing, consumer offsets, retries, and fan-out (ADR-0004).
package events

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

const (
	// Topic is the single durable PostgreSQL stream for domain facts.
	Topic = "domain_events"

	userSignedUpName    = "user_signed_up"
	userSignedInName    = "user_signed_in"
	bookRequestedName   = "book_requested"
	courseStartedName   = "course_started"
	courseFinishedName  = "course_finished"
	lessonStartedName   = "lesson_started"
	lessonFinishedName  = "lesson_finished"
	solutionCheckedName = "solution_checked"
	emailConfirmedName  = "email_confirmed"
	leadCreatedName     = "lead_created"
	staffReadName       = "staff_read_personal_data"
	schemaVersion       = "1"
)

// Event is closed to this package so arbitrary transport payloads cannot be
// published as domain facts without first defining their stable contract here.
type Event interface {
	eventName() string
}

// UserSignedUp is the immutable account snapshot captured when public
// registration first succeeds. Its JSON shape mirrors the legacy event while
// adding the occurrence time needed by asynchronous consumers.
type UserSignedUp struct {
	UserID     int       `json:"user_id"`
	Email      *string   `json:"email"`
	FirstName  *string   `json:"first_name"`
	LastName   *string   `json:"last_name"`
	Locale     string    `json:"locale"`
	OccurredAt time.Time `json:"occurred_at"`
}

func (UserSignedUp) eventName() string { return userSignedUpName }

// UserSignedIn records a successful authentication.
type UserSignedIn struct {
	UserID          int       `json:"user_id"`
	OccurrenceCount int       `json:"occurrence_count"`
	Email           *string   `json:"email"`
	Locale          string    `json:"locale"`
	OccurredAt      time.Time `json:"occurred_at"`
}

func (UserSignedIn) eventName() string { return userSignedInName }

// BookRequested records the first request for the downloadable book.
type BookRequested struct {
	Locale     string    `json:"locale"`
	OccurredAt time.Time `json:"occurred_at"`
}

func (BookRequested) eventName() string { return bookRequestedName }

// CourseStarted records the first enrollment transition for a course.
type CourseStarted struct {
	OccurrenceCount int       `json:"occurrence_count"`
	Slug            string    `json:"slug"`
	Locale          string    `json:"locale"`
	OccurredAt      time.Time `json:"occurred_at"`
}

func (CourseStarted) eventName() string { return courseStartedName }

// CourseFinished records a course completion transition.
type CourseFinished struct {
	OccurrenceCount int       `json:"occurrence_count"`
	Slug            string    `json:"slug"`
	Locale          string    `json:"locale"`
	OccurredAt      time.Time `json:"occurred_at"`
}

func (CourseFinished) eventName() string { return courseFinishedName }

// LessonStarted records the first start of a lesson within a course.
type LessonStarted struct {
	OccurrenceCount int       `json:"occurrence_count"`
	LessonSlug      string    `json:"lesson_slug"`
	CourseSlug      string    `json:"course_slug"`
	Locale          string    `json:"locale"`
	OccurredAt      time.Time `json:"occurred_at"`
}

func (LessonStarted) eventName() string { return lessonStartedName }

// LessonFinished records a lesson completion transition.
type LessonFinished struct {
	OccurrenceCount int       `json:"occurrence_count"`
	LessonSlug      string    `json:"lesson_slug"`
	CourseSlug      string    `json:"course_slug"`
	Locale          string    `json:"locale"`
	OccurredAt      time.Time `json:"occurred_at"`
}

func (LessonFinished) eventName() string { return lessonFinishedName }

// SolutionChecked records every exercise check, including guest checks.
type SolutionChecked struct {
	LessonSlug string    `json:"lesson_slug"`
	CourseSlug string    `json:"course_slug"`
	Locale     string    `json:"locale"`
	Passed     bool      `json:"passed"`
	OccurredAt time.Time `json:"occurred_at"`
}

func (SolutionChecked) eventName() string { return solutionCheckedName }

// EmailConfirmed records the confirmation transition.
type EmailConfirmed struct {
	OccurredAt time.Time `json:"occurred_at"`
}

func (EmailConfirmed) eventName() string { return emailConfirmedName }

// LeadCreated is the integration snapshot consumed by the amoCRM handler.
type LeadCreated struct {
	LeadID      int     `json:"lead_id"`
	UserID      int     `json:"user_id"`
	UserName    string  `json:"user_name"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	YMClientID  *string `json:"ym_client_id"`
	UTMSource   *string `json:"utm_source"`
	UTMMedium   *string `json:"utm_medium"`
	UTMCampaign *string `json:"utm_campaign"`
	UTMTerm     *string `json:"utm_term"`
	UTMContent  *string `json:"utm_content"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Telegram    *string `json:"telegram"`
	WhatsApp    *string `json:"whatsapp"`
	// LandingPage, Referrer and IP are what the amoCRM payload used to read off
	// ahoy's visit (form page, referer, IP, and the click ids parsed from the
	// landing URL). The first two come from the browser's first-visit cookie,
	// the IP from the request (ADR-0015). IP is the raw client address; the
	// amoCRM adapter decides what it can send.
	LandingPage *string   `json:"landing_page"`
	Referrer    *string   `json:"referrer"`
	IP          *string   `json:"ip"`
	OccurredAt  time.Time `json:"occurred_at"`
}

func (LeadCreated) eventName() string { return leadCreatedName }

// StaffReadPersonalData is the fact that a staff member read personal data
// through the admin surface or over MCP (ADR-0019). The request publishes it and
// a worker consumer writes the Audit Record, so recording never slows the read.
//
// It carries identifiers and a count, never a field value: EntityIDs are the ids
// of the rows the response held, up to a cap (the whole count is in Count), so
// the record can say which learners were looked at without repeating what was
// seen. The rest is who and from where, copied from the request because the
// consumer has no request.
type StaffReadPersonalData struct {
	ActorUserID    *int      `json:"actor_user_id"`
	Channel        string    `json:"channel"`
	ConnectedAppID *int      `json:"connected_app_id"`
	TokenID        *string   `json:"token_id"`
	Operation      string    `json:"operation"`
	EntityType     string    `json:"entity_type"`
	EntityIDs      []string  `json:"entity_ids"`
	Count          int       `json:"count"`
	RequestID      string    `json:"request_id"`
	Address        string    `json:"address"`
	UserAgent      string    `json:"user_agent"`
	Locale         string    `json:"locale"`
	OccurredAt     time.Time `json:"occurred_at"`
}

func (StaffReadPersonalData) eventName() string { return staffReadName }

// Marshaler encodes events into outbox messages and decodes them back; the
// stable event name travels in the message metadata. The publisher writes with
// it and tests read the outbox with it.
var Marshaler = cqrs.JSONMarshaler{GenerateName: nameOf}

func nameOf(value any) string {
	if event, ok := value.(Event); ok {
		return event.eventName()
	}
	return ""
}
