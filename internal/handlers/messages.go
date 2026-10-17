package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/domain"
)

// AdminListMessages returns a page of in-lesson assistant messages, newest
// first (legacy `admin/messages#index` over `ai_messages`). The admin surface
// is read-only (index only).
func (s *Server) AdminListMessages(ctx context.Context, params api.AdminListMessagesParams) (api.AdminListMessagesRes, error) {
	locale := s.i18n.Locale(ctx)
	return listFrom(ctx, params.Page, params.PerPage,
		func(ctx context.Context, page domain.Page) ([]*ent.AiMessage, int, error) {
			return s.messageDirectory.Messages(ctx, locale, params.Q.Value, page)
		},
		s.conv.ToLessonAssistantMessages,
		func(items []api.LessonAssistantMessage, total, page, perPage int32) *api.LessonAssistantMessagePage {
			return &api.LessonAssistantMessagePage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}
