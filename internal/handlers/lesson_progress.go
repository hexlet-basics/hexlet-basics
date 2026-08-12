package handlers

import (
	"context"
	"errors"
	"net/http"

	"hexletbasics/internal/api"
	"hexletbasics/internal/progress"
)

// StartLesson marks a lesson started for the signed-in learner, enrolling them
// in its course if needed. The rule — which lessons a learner may take — lives
// in the progress module; this handler only translates its outcome to the
// contract's responses.
//
// 409 rather than 403 for a gated lesson: the caller's permissions are fine,
// their progress is not, and a client that cannot tell those apart cannot
// respond to either sensibly.
func (s *Server) StartLesson(ctx context.Context, params api.StartLessonParams) (api.StartLessonRes, error) {
	u, ok := AuthenticatedUser(ctx)
	if !ok {
		// Defensive: the generated security handler rejects unauthenticated
		// requests before this runs. Without the status the central writer would
		// report a missing user as a server fault.
		return nil, withHTTPStatus(http.StatusUnauthorized, errUnauthenticated)
	}

	err := s.progress.StartLesson(ctx, u.ID, int(params.ID), s.i18n.Locale(ctx))
	switch {
	case errors.Is(err, progress.ErrLessonNotAvailable):
		return s.errors.LessonNotAvailable(ctx), nil
	case err != nil:
		// ent's not-found reaches the central handler as 404; anything else is a
		// server fault and is reported there.
		return nil, err
	}
	return &api.StartLessonNoContent{}, nil
}
