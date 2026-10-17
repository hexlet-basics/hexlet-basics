package assistant

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/aimessage"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/internal/domain"
)

// Directory is the administrators' read side of the assistant's conversations.
type Directory struct {
	db *ent.Client
}

// NewDirectory builds the directory over an ent client.
func NewDirectory(db *ent.Client) *Directory {
	return &Directory{db: db}
}

// Messages lists the in-lesson assistant messages, newest first (legacy
// `admin/messages#index` over `ai_messages`). Course and lesson identity comes
// through chat, lesson progress and the course and the lesson; the lesson's
// translations are narrowed to the locale and ordered by id, so the "first
// translation" name lookup matches the lesson progress list's convention.
func (d *Directory) Messages(ctx context.Context, locale string, page domain.Page) ([]*ent.AiMessage, int, error) {
	return domain.Paged(ctx,
		d.db.AiMessage.Query().
			WithChat(func(q *ent.AiChatQuery) {
				q.WithLessonProgress(func(q *ent.LessonProgressQuery) {
					q.WithCourse().WithLesson(func(q *ent.CourseLessonQuery) {
						q.WithInfos(func(q *ent.CourseLessonTranslationQuery) {
							q.Where(courselessontranslation.LocaleEQ(locale)).
								Order(ent.Asc(courselessontranslation.FieldID))
						})
					})
				})
			}).
			Order(ent.Desc(aimessage.FieldID)),
		page)
}
