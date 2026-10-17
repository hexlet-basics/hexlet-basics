package handlers

import (
	"context"
	"time"

	"hexletbasics/ent"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/blogposts"
	"hexletbasics/internal/coursecategories"
	"hexletbasics/internal/courses"
	"hexletbasics/internal/dashboard"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/leads"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/qnaitems"
	"hexletbasics/internal/selfprotection"
	"hexletbasics/internal/sitemap"
)

// The interfaces below are what the handlers need from the use-case modules,
// declared here, with the caller, and not by the modules: each names only the
// methods a handler calls, so a test (or another adapter) fills in just those.
// The modules return concrete types; production hands the handlers
// *progress.Progress, *lead.Recorder and the rest.

// ProgressTracker is where handlers read and move a learner's progress.
type ProgressTracker interface {
	StartLesson(ctx context.Context, learner progress.Learner, lessonID int, locale string) (*progress.CourseState, error)
	CheckSolution(ctx context.Context, check progress.Check) (*progress.CheckResult, error)
	CourseState(ctx context.Context, learner progress.Learner, courseID int) (*progress.CourseState, error)
	LearnerCourse(ctx context.Context, userID int, courseSlug string) (*progress.LearnerCourse, error)
	MergeGuest(ctx context.Context, userID int, guest progress.GuestProgress, locale string) error
}

// SitemapReader reads the rows behind the sitemap page.
type SitemapReader interface {
	Read(ctx context.Context) (*sitemap.Page, error)
}

// PermissionWriter writes what a role may do: its grants. Check judges a
// submission (the role is there, every grant is declared) before Apply writes it.
type PermissionWriter interface {
	Check(ctx context.Context, roleID int, grants []string) error
	Apply(ctx context.Context, roleID int, grants []string) error
}

// SelfProtection judges an actor's change to staff management: an
// administrator-only flag, one's own Staff Member record, role grants and
// account (the selfprotection module).
type SelfProtection interface {
	StaffRecordChange(ctx context.Context, actor selfprotection.Actor, recordID, newOwner int) error
	OwnRoleGrants(ctx context.Context, actor selfprotection.Actor, roleID int) error
	AdministratorFlag(actor selfprotection.Actor, was, will bool) error
	OwnAccountDeletion(actor selfprotection.Actor, userID int) error
}

// CourseLister assembles the learner's dashboard courses.
type CourseLister interface {
	Courses(ctx context.Context, userID int, locale string) (dashboard.Courses, error)
}

// LeadCreator stores a lead and raises LeadCreated with it.
type LeadCreator interface {
	Create(ctx context.Context, submission leads.Submission) (*ent.Lead, error)
}

// BookRequester is the book funnel, from the request to the download.
type BookRequester interface {
	Request(ctx context.Context, userID int, locale string) error
	// Requested reports whether the user has requested the book, in any state.
	Requested(ctx context.Context, userID int) (bool, error)
	// MarkDownloaded records a download and reports whether the user had
	// requested the book; one who had not has nothing to mark.
	MarkDownloaded(ctx context.Context, userID int) (bool, error)
}

// AccountRemover deletes an account.
type AccountRemover interface {
	Remove(ctx context.Context, userID int) error
}

// UserRegistrar creates an account and raises UserSignedUp with it.
type UserRegistrar interface {
	Register(ctx context.Context, input accounts.Registration) (*ent.User, error)
}

// PostRemover deletes a blog post with its dependents.
type PostRemover interface {
	Remove(ctx context.Context, postID int) error
}

// UserEditor is how an administrator creates, changes and deletes user records.
type UserEditor interface {
	Create(ctx context.Context, fields accounts.UserFields) (*ent.User, error)
	Update(ctx context.Context, id int, fields accounts.UserFields) (*ent.User, error)
	Delete(ctx context.Context, id int) error
}

// UserDirectory is how an administrator reads user records.
type UserDirectory interface {
	List(ctx context.Context, page domain.Page) ([]*ent.User, int, error)
	Search(ctx context.Context, term string) ([]*ent.User, error)
	Get(ctx context.Context, id int) (*ent.User, error)
}

// LessonDirectory is the administrators' read side of the lesson graph.
type LessonDirectory interface {
	Lessons(ctx context.Context, locale string, page domain.Page) ([]*ent.CourseLessonTranslation, int, error)
	Progress(ctx context.Context, locale string, page domain.Page) ([]*ent.LessonProgress, int, error)
	Reviews(ctx context.Context, locale string, page domain.Page) ([]*ent.CourseLessonReview, int, error)
}

// BannerDirectory reads the banner a page shows.
type BannerDirectory interface {
	Active(ctx context.Context, locale string, now time.Time) (*ent.Banner, error)
}

// CategoryDirectory reads the course categories a visitor browses.
type CategoryDirectory interface {
	Index(ctx context.Context, locale string) ([]*ent.CourseCategory, error)
	View(ctx context.Context, locale, slug string) (*coursecategories.View, error)
}

// MessageDirectory is the administrators' read side of the assistant's messages.
type MessageDirectory interface {
	Messages(ctx context.Context, locale string, page domain.Page) ([]*ent.AiMessage, int, error)
}

// LeadDirectory is the administrators' read side of the leads.
type LeadDirectory interface {
	List(ctx context.Context, page domain.Page) ([]*ent.Lead, int, error)
}

// CourseDirectory is the administrators' read side of the courses.
type CourseDirectory interface {
	Options(ctx context.Context) ([]courses.Option, error)
}

// ReviewDirectory reads the reviews a visitor may see.
type ReviewDirectory interface {
	Public(ctx context.Context, locale string, page domain.Page) ([]*ent.Review, int, error)
}

// StateMover moves the Publication State of a resource: each move is an
// operation of its own, and a refused one is a *transitions.Illegal.
type StateMover interface {
	Publish(ctx context.Context, id int) error
	Unpublish(ctx context.Context, id int) error
	Archive(ctx context.Context, id int) error
}

// ReadinessMarker moves the Readiness of a course: each value is an operation of
// its own, and a refused one is a *transitions.Illegal.
type ReadinessMarker interface {
	MarkDraft(ctx context.Context, id int) error
	MarkInDevelopment(ctx context.Context, id int) error
	MarkCompleted(ctx context.Context, id int) error
}

// PostEditor creates and changes blog posts, and moves their Publication State.
type PostEditor interface {
	StateMover
	Create(ctx context.Context, creatorID int, locale string, fields blogposts.PostFields) (*ent.BlogPost, error)
	Update(ctx context.Context, id int, fields blogposts.PostFields) (*ent.BlogPost, error)
}

// LessonReviewScheduler asks for AI reviews of a course's or a lesson's
// translations.
type LessonReviewScheduler interface {
	ReviewCourse(ctx context.Context, courseID int) error
	ReviewLesson(ctx context.Context, lessonID int) error
}

// RelatedCoursesScheduler asks for the AI pick of a post's related courses.
type RelatedCoursesScheduler interface {
	Suggest(ctx context.Context, postID int) error
}

// CategoryQna edits the questions and answers of a course category.
type CategoryQna interface {
	List(ctx context.Context, categoryID int) ([]*ent.CategoryQnaItem, error)
	Create(ctx context.Context, categoryID int, fields qnaitems.Fields) (*ent.CategoryQnaItem, error)
	Update(ctx context.Context, categoryID, id int, fields qnaitems.Fields) (*ent.CategoryQnaItem, error)
	Delete(ctx context.Context, categoryID, id int) error
}

// LandingPageQna edits the questions and answers of a course landing page.
type LandingPageQna interface {
	List(ctx context.Context, landingPageID int) ([]*ent.LandingPageQnaItem, error)
	Create(ctx context.Context, landingPageID int, fields qnaitems.Fields) (*ent.LandingPageQnaItem, error)
	Update(ctx context.Context, landingPageID, id int, fields qnaitems.Fields) (*ent.LandingPageQnaItem, error)
	Delete(ctx context.Context, landingPageID, id int) error
}
