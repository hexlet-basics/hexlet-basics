package handlers

import (
	"context"
	"errors"

	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/progress"
)

// StartLesson marks a lesson started. What starting means lives in the progress
// module; this handler only translates its outcome to the contract's responses.
//
// A guest reaches it too, and nothing is stored for them: their state is the
// signed cookie, and only a check moves it.
func (s *Server) StartLesson(ctx context.Context, params api.StartLessonParams) error {
	_, err := s.progress.StartLesson(
		ctx,
		currentLearner(ctx),
		int(params.ID),
		s.i18n.Locale(ctx),
	)
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return errLessonNotFound
	case err != nil:
		// Anything else is a server fault, reported by the central handler.
		return err
	}

	return nil
}

// CheckLesson runs a submitted solution and reports what it changed.
//
// Public, and deliberately so: a guest must be able to submit. Which storage
// answers is the progress module's business — the handler passes on who is
// asking, and hands back the cookie the module produced when there is one,
// so the response body is identical either way.
func (s *Server) CheckLesson(
	ctx context.Context,
	req *api.CheckLessonInput,
	params api.CheckLessonParams,
) (*api.LessonCheckingResponseHeaders, error) {
	result, err := s.progress.CheckSolution(ctx, progress.Check{
		Learner:   currentLearner(ctx),
		LessonID:  int(params.ID),
		VersionID: int(req.VersionId),
		Code:      req.Code,
		Locale:    s.i18n.Locale(ctx),
	})
	switch {
	case errors.Is(err, domain.ErrNotFound):
		// The lesson, or the version the code was written against, is not part
		// of the course's current build.
		return nil, errLessonNotFound
	case err != nil:
		return nil, err
	}

	response := &api.LessonCheckingResponseHeaders{
		Response: api.LessonCheckingResponse{
			Passed:                result.Outcome.Passed,
			Output:                result.Outcome.Output,
			Result:                checkResultOf(result.Outcome),
			Status:                apiconv.Int32FromInt(result.Outcome.Status),
			LessonHasBeenFinished: result.LessonFinished,
			CourseHasBeenFinished: result.CourseFinished,
		},
	}
	if result.Guest != nil {
		cookie, err := s.guests.set(*result.Guest)
		if err != nil {
			return nil, err
		}

		response.SetCookie = []string{cookie}
	}

	return response, nil
}

// currentLearner is whose progress a request is about: the signed-in user when
// the session cookie identified one, and otherwise the visitor's cookie
// progress — empty for a first visit, which is exactly the position of someone
// who has finished nothing.
func currentLearner(ctx context.Context) progress.Learner {
	if u, ok := AuthenticatedUser(ctx); ok {
		return progress.Learner{
			UserID: u.ID,
			Guest: progress.GuestProgress{
				Entries: nil,
			},
		}
	}

	guest, _ := GuestProgress(ctx)

	return progress.Learner{
		Guest:  guest,
		UserID: 0,
	}
}

// checkResultOf maps the runner's classification onto the contract's enum. An
// unrecognized classification reads as a plain failure: the one thing it must
// not do is report a pass the runner did not give.
func checkResultOf(outcome progress.Outcome) api.LessonCheckingResponseResult {
	switch outcome.Result {
	case progress.ResultPassed:
		return api.LessonCheckingResponseResultPassed
	case progress.ResultFailedInfinity:
		return api.LessonCheckingResponseResultFailedInfinity
	default:
		return api.LessonCheckingResponseResultFailed
	}
}
