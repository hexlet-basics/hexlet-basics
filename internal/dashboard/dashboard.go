// Package dashboard assembles the learner's course dashboard.
package dashboard

import (
	"context"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/enrollment"
	"hexletbasics/internal/landingpages"
	"hexletbasics/internal/progress"
)

// Course is one course on the dashboard: the enrollment, its card and its
// progress travel together, so a caller cannot pair one course's card with
// another's state.
type Course struct {
	Enrollment *ent.Enrollment
	Card       *ent.CourseLandingPage
	State      *progress.CourseState
}

// Courses is the dashboard: what the learner is working on and what they have
// finished, each newest enrollment first.
type Courses struct {
	Started  []Course
	Finished []Course
}

// CourseLister is the dashboard seam the handlers depend on.
type CourseLister interface {
	Courses(ctx context.Context, userID int, locale string) (Courses, error)
}

// Dashboard reads the learner's courses.
type Dashboard struct {
	db       *ent.Client
	catalog  *landingpages.Catalog
	progress progress.Tracker
}

// New builds the production dashboard over the catalog and the progress module,
// so the dashboard and the course page report the same progress.
func New(db *ent.Client, catalog *landingpages.Catalog, tracker progress.Tracker) *Dashboard {
	return &Dashboard{db: db, catalog: catalog, progress: tracker}
}

// Courses returns the user's enrollments, newest first, each with its card and
// progress. Only courses with a card appear, matching the legacy scope: the
// dashboard renders course cards, and a course with none would be a row the
// client cannot draw.
//
// Newest first because a learner returning after an absence finds what they
// most recently touched at the top, and the order is total. The card is the
// course's main landing page in locale, legacy `LandingPage.web.where(main:
// true)`: the dashboard scope, which unlike the public lists does not require
// `listed`. One query serves every enrollment.
func (d *Dashboard) Courses(ctx context.Context, userID int, locale string) (Courses, error) {
	enrollments, err := d.db.Enrollment.Query().
		Where(enrollment.UserID(userID)).
		Order(ent.Desc(enrollment.FieldID)).
		All(ctx)
	if err != nil {
		return Courses{}, err
	}

	courseIDs := lo.Map(enrollments, func(e *ent.Enrollment, _ int) int { return e.CourseID })
	cards, err := d.catalog.MainPages(ctx, locale, courseIDs...)
	if err != nil {
		return Courses{}, err
	}
	states, err := d.progress.CourseStates(ctx, userID, courseIDs)
	if err != nil {
		return Courses{}, err
	}

	var dashboard Courses
	for _, e := range enrollments {
		card, ok := cards[e.CourseID]
		if !ok {
			continue
		}
		course := Course{Enrollment: e, Card: card, State: states[e.CourseID]}
		if course.State.State == progress.StateFinished {
			dashboard.Finished = append(dashboard.Finished, course)
		} else {
			dashboard.Started = append(dashboard.Started, course)
		}
	}
	return dashboard, nil
}
