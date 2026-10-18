// Package dashboard assembles the learner's course dashboard.
package dashboard

import (
	"context"
	"log/slog"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courseversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/internal/landingpage"
	"hexletbasics/internal/progress"
)

// Course is one course on the dashboard: the enrollment, its card and its
// progress travel together, so a caller cannot pair one course's card with
// another's state.
type Course struct {
	Enrollment *ent.Enrollment
	Card       *ent.CourseLandingPage
	State      *progress.CourseState

	// NextLesson is the Next Lesson as the request locale knows it: its id and
	// what it is called. Nil when there is none, or it has no translation in the
	// locale and so no page there; the state knows the lesson only by its slug.
	NextLesson *NextLesson
}

// NextLesson is a lesson to resume at, addressed and named.
type NextLesson struct {
	ID   int
	Name string
}

// Courses is the dashboard: what the learner is working on and what they have
// finished, each newest enrollment first.
type Courses struct {
	Started  []Course
	Finished []Course
}

// Dashboard reads the learner's courses.
type Dashboard struct {
	db       *ent.Client
	catalog  *landingpage.Catalog
	progress *progress.Progress
	log      *slog.Logger
}

// New builds the production dashboard over the catalog and the progress module,
// so the dashboard and the course page report the same progress. The logger
// reports the courses left off for lack of a card.
func New(db *ent.Client, catalog *landingpage.Catalog, tracker *progress.Progress, logger *slog.Logger) *Dashboard {
	return &Dashboard{db: db, catalog: catalog, progress: tracker, log: logger}
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

	next, err := d.nextLessons(ctx, locale, states)
	if err != nil {
		return Courses{}, err
	}

	var dashboard Courses
	for _, e := range enrollments {
		card, ok := cards[e.CourseID]
		if !ok {
			// Skipped, not failed: the learner's page is fine without it, but
			// the course is a content gap somebody has to hear about.
			d.log.WarnContext(ctx, "dashboard course has no catalog card",
				"course_id", e.CourseID, "user_id", userID, "locale", locale)
			continue
		}
		state := states[e.CourseID]
		listed := Course{Enrollment: e, Card: card, State: state}
		if state.NextLessonSlug != nil {
			if lesson, ok := next[lessonKey{courseID: e.CourseID, slug: *state.NextLessonSlug}]; ok {
				listed.NextLesson = &lesson
			}
		}
		if state.State == progress.StateFinished {
			dashboard.Finished = append(dashboard.Finished, listed)
		} else {
			dashboard.Started = append(dashboard.Started, listed)
		}
	}
	return dashboard, nil
}

// lessonKey names a lesson across courses: slugs are unique only within one.
type lessonKey struct {
	courseID int
	slug     string
}

// nextLessons resolves each course's Next Lesson in locale, in one query for
// every course. A lesson that has no translation in the locale has no entry: the
// dashboard then shows the course without a resume link rather than a slug where
// a name belongs.
func (d *Dashboard) nextLessons(
	ctx context.Context,
	locale string,
	states map[int]*progress.CourseState,
) (map[lessonKey]NextLesson, error) {
	var courseIDs []int
	var slugs []string
	for courseID, state := range states {
		if state.NextLessonSlug != nil {
			courseIDs = append(courseIDs, courseID)
			slugs = append(slugs, *state.NextLessonSlug)
		}
	}
	if len(slugs) == 0 {
		return map[lessonKey]NextLesson{}, nil
	}

	rows, err := d.db.CourseLessonTranslation.Query().
		Where(
			courselessontranslation.LocaleEQ(locale),
			courselessontranslation.HasCourseVersionWith(
				courseversion.HasCurrentCoursesWith(course.IDIn(courseIDs...)),
			),
			courselessontranslation.HasLessonWith(courselesson.SlugIn(slugs...)),
		).
		WithLesson().
		All(ctx)
	if err != nil {
		return nil, err
	}

	lessons := make(map[lessonKey]NextLesson, len(rows))
	for _, row := range rows {
		if row.Edges.Lesson == nil || row.Edges.Lesson.Slug == nil || row.Name == nil {
			continue
		}
		lessons[lessonKey{courseID: row.CourseID, slug: *row.Edges.Lesson.Slug}] = NextLesson{
			ID:   row.CourseLessonID,
			Name: *row.Name,
		}
	}
	return lessons, nil
}
