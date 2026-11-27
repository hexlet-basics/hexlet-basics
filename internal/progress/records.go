package progress

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/store"
)

// Records lists the learners' Lesson Progress records, newest first (legacy
// `admin/lesson_progress`), for the administrators' list. The course and the
// lesson, with its translation in locale, come with each row.
func (p *Progress) Records(ctx context.Context, locale string, page domain.Page) ([]*ent.LessonProgress, int, error) {
	return store.Paged(ctx,
		p.db.LessonProgress.Query().
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
