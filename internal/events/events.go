// Package events owns domain-event contracts and their Watermill transport.
//
// Domain callers publish immutable facts through the transactional publisher;
// Watermill owns routing, consumer offsets, retries, and fan-out (ADR-0004).
package events

import "time"

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
	LeadID      int       `json:"lead_id"`
	UserID      int       `json:"user_id"`
	UserName    string    `json:"user_name"`
	FirstName   *string   `json:"first_name"`
	LastName    *string   `json:"last_name"`
	YMClientID  *string   `json:"ym_client_id"`
	UTMSource   *string   `json:"utm_source"`
	UTMMedium   *string   `json:"utm_medium"`
	UTMCampaign *string   `json:"utm_campaign"`
	UTMTerm     *string   `json:"utm_term"`
	UTMContent  *string   `json:"utm_content"`
	Email       *string   `json:"email"`
	Phone       *string   `json:"phone"`
	Telegram    *string   `json:"telegram"`
	WhatsApp    *string   `json:"whatsapp"`
	OccurredAt  time.Time `json:"occurred_at"`
}

func (LeadCreated) eventName() string { return leadCreatedName }

func nameOf(value any) string {
	if event, ok := value.(Event); ok {
		return event.eventName()
	}
	return ""
}
