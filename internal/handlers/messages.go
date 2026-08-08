package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/aimessage"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/internal/api"
)

// AdminListMessages returns a page of in-lesson assistant messages, newest
// first (legacy `admin/messages#index` over `ai_messages`). The admin surface
// is read-only (index only). Course and lesson identity is enriched through
// chat → member → course/lesson; the lesson's infos are locale-filtered and
// id-ordered so the converter's "first info" name lookup matches the member
// list's convention.
func (s *Server) AdminListMessages(ctx context.Context, params api.AdminListMessagesParams) (api.AdminListMessagesRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.AiMessageQuery {
			return s.db.AiMessage.Query().
				WithChat(func(q *ent.AiChatQuery) {
					q.WithMember(func(q *ent.LessonProgressQuery) {
						q.WithCourse().WithLesson(func(q *ent.CourseLessonQuery) {
							q.WithInfos(func(q *ent.CourseLessonTranslationQuery) {
								q.Where(courselessontranslation.LocaleEQ(defaultAdminLocale)).
									Order(ent.Asc(courselessontranslation.FieldID))
							})
						})
					})
				}).
				Order(ent.Desc(aimessage.FieldID))
		},
		s.conv.ToLessonAssistantMessages,
		func(items []api.LessonAssistantMessage, total, page, perPage int32) *api.LessonAssistantMessagePage {
			return &api.LessonAssistantMessagePage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}
