package handlers

import (
	"context"
	"errors"
	"net/http"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/progress"
)

// ListAssistantMessages returns the learner's in-lesson chat and whether they
// may still ask today, so the panel renders its disabled state on first paint.
func (s *Server) ListAssistantMessages(
	ctx context.Context,
	params api.ListAssistantMessagesParams,
) (api.ListAssistantMessagesRes, error) {
	u, ok := AuthenticatedUser(ctx)
	if !ok {
		return nil, errUnauthenticated
	}

	history, err := s.assistant.History(ctx, u.ID, int(params.LessonId))
	switch {
	case ent.IsNotFound(err):
		return s.lessonNotFound(ctx), nil
	case err != nil:
		return nil, err
	}
	return &api.LessonAssistantChat{
		Messages:      s.conv.ToLessonAssistantMessages(history.Messages),
		QuotaExceeded: history.QuotaExceeded,
	}, nil
}

// CreateAssistantMessage asks the assistant and answers with its reply as a
// plain text stream. Every refusal is decided inside Ask, before the first
// byte, so it still reaches the learner as a typed status; the router flushes
// each chunk of the body the generated encoder copies (see StreamResponses).
func (s *Server) CreateAssistantMessage(
	ctx context.Context,
	req *api.AssistantMessageInput,
	params api.CreateAssistantMessageParams,
) (api.CreateAssistantMessageRes, error) {
	u, ok := AuthenticatedUser(ctx)
	if !ok {
		return nil, errUnauthenticated
	}

	answer, err := s.assistant.Ask(ctx, assistant.Question{
		UserID:   u.ID,
		LessonID: int(params.LessonId),
		Message:  req.Message,
		// Null and empty mean the same to the model: nothing written, nothing run.
		UserCode: req.UserCode.Value,
		Output:   req.Output.Value,
	})
	switch {
	case errors.Is(err, assistant.ErrQuotaExceeded):
		// A declared outcome, not a failure: it never reaches the reporting path.
		problem := api.CreateAssistantMessageTooManyRequests(s.errors.problem(ctx, http.StatusTooManyRequests))
		return &problem, nil
	case errors.Is(err, progress.ErrLessonNotAvailable):
		return (*api.CreateAssistantMessageConflict)(s.errors.LessonNotAvailable(ctx)), nil
	case ent.IsNotFound(err):
		return s.lessonNotFound(ctx), nil
	case err != nil:
		return nil, err
	}
	return &api.CreateAssistantMessageOK{Data: answer}, nil
}
