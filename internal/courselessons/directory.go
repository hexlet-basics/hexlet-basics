// Package courselessons reads the lesson graph for the administrators' index
// pages: the lessons of the current course versions, how far learners have got
// with each, and the AI reviews of them.
package courselessons

import (
	"context"
	"strings"

	"hexletbasics/ent"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessonreview"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courseversion"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/landingpages"
)

// Directory is the read side of the lesson graph. Every list is in the locale it
// is asked for: the admin's interface language, which is what legacy filtered
// these pages by.
type Directory struct {
	db *ent.Client
}

// NewDirectory builds the directory over an ent client.
func NewDirectory(db *ent.Client) *Directory {
	return &Directory{db: db}
}

// Lessons lists lessons as the index shows them: one row per lesson for the
// locale, drawn from the current version's translation (legacy
// `Language::Lesson::Version::Info.current.with_locale`). The name and the
// description come from the translation row; the slug is joined from the lesson.
//
// A non-blank search keeps the rows whose name or description, or whose
// lesson's slug, contains it, case-insensitively. It narrows the rows the list
// already shows and never adds one the list hides.
func (d *Directory) Lessons(ctx context.Context, locale, search string, page domain.Page) ([]*ent.CourseLessonTranslation, int, error) {
	query := d.db.CourseLessonTranslation.Query().
		Where(
			courselessontranslation.LocaleEQ(locale),
			// Legacy `.current`: the translation's course version must be
			// selected as current by a completed course.
			courselessontranslation.HasCourseVersionWith(
				courseversion.HasCurrentCoursesWith(landingpages.Completed()),
			),
		).
		WithLesson().
		// Legacy orders by the translation id (ransack `sf=id`), not by the
		// lesson id this row reports as its own, so the order is the
		// translation's field.
		Order(ent.Desc(courselessontranslation.FieldID))
	if term := strings.TrimSpace(search); term != "" {
		query = query.Where(courselessontranslation.Or(
			courselessontranslation.NameContainsFold(term),
			courselessontranslation.DescriptionContainsFold(term),
			courselessontranslation.HasLessonWith(courselesson.SlugContainsFold(term)),
		))
	}
	return domain.Paged(ctx, query, page)
}

// Progress lists learners' per-lesson progress rows, newest first (legacy
// `admin/lesson_progress`). The course and the lesson, with its translation in
// the locale, come with each row.
func (d *Directory) Progress(ctx context.Context, locale string, page domain.Page) ([]*ent.LessonProgress, int, error) {
	return domain.Paged(ctx,
		d.db.LessonProgress.Query().
			WithCourse().
			WithLesson(func(q *ent.CourseLessonQuery) {
				q.WithInfos(func(q *ent.CourseLessonTranslationQuery) {
					q.Where(courselessontranslation.LocaleEQ(locale)).
						Order(ent.Asc(courselessontranslation.FieldID))
				})
			}).
			Order(ent.Desc(lessonprogress.FieldID)),
		page)
}

// Reviews lists the AI lesson reviews of the locale that have a summary, newest
// first (legacy `admin/course_lesson_reviews`, `where(locale:).with_summary`).
// A review with an empty summary marks a lesson no student asked about; legacy
// skips it, and it must not inflate the total, so the filter is on the one query
// the count and the page share.
func (d *Directory) Reviews(ctx context.Context, locale string, page domain.Page) ([]*ent.CourseLessonReview, int, error) {
	return domain.Paged(ctx,
		d.db.CourseLessonReview.Query().
			Where(
				courselessonreview.LocaleEQ(locale),
				courselessonreview.SummaryNEQ(""),
			).
			WithCourse().
			WithLesson().
			Order(ent.Desc(courselessonreview.FieldID)),
		page)
}
