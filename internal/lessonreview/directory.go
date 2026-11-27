package lessonreview

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/courselessonreview"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/store"
)

// Directory is the read side of the AI lesson reviews, for the
// administrators' list.
type Directory struct {
	db *ent.Client
}

// NewDirectory builds the directory over an ent client.
func NewDirectory(db *ent.Client) *Directory {
	return &Directory{db: db}
}

// List lists the reviews of the locale that have a summary, newest first
// (legacy `admin/course_lesson_reviews`, `where(locale:).with_summary`). A
// review with an empty summary marks a lesson no student asked about; legacy
// skips it, and it must not inflate the total, so the filter is on the one
// query the count and the page share.
func (d *Directory) List(ctx context.Context, locale string, page domain.Page) ([]*ent.CourseLessonReview, int, error) {
	return store.Paged(ctx,
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
