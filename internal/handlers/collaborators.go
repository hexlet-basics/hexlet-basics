package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/dashboard"
	"hexletbasics/internal/leads"
	"hexletbasics/internal/progress"
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

// PermissionWriter writes what a role may do: its grants.
type PermissionWriter interface {
	Apply(ctx context.Context, roleID int, grants []string) error
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
