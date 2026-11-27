package curriculum

import (
	"context"
	"strings"

	"hexletbasics/ent"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courseversion"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/landingpage"
	"hexletbasics/internal/store"
)

// Directory is the read side of the lesson graph: the administrators' list of
// lessons, in the locale it is asked for (the admin's interface language, which
// is what legacy filtered the page by). The Lesson Progress records are listed
// by progress and the AI lesson reviews by lessonreview, each with its concept.
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
				courseversion.HasCurrentCoursesWith(landingpage.Completed()),
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

	return store.Paged(ctx, query, page)
}
