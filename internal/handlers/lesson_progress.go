package handlers

import (
	"context"
	"errors"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/progress"
)

// StartLesson marks a lesson started and answers with where the learner now
// stands. The rule — which lessons a learner may take — lives in the progress
// module; this handler only translates its outcome to the contract's responses.
//
// A guest reaches it too, and gets the same payload with nothing stored: their
// position is the signed cookie, and only a check moves it.
//
// 409 rather than 403 for a gated lesson: the caller's permissions are fine,
// their progress is not, and a client that cannot tell those apart cannot
// respond to either sensibly.
func (s *Server) StartLesson(ctx context.Context, params api.StartLessonParams) (api.StartLessonRes, error) {
	state, err := s.progress.StartLesson(
		ctx,
		currentLearner(ctx),
		int(params.ID),
		s.i18n.Locale(ctx),
	)
	switch {
	case errors.Is(err, progress.ErrLessonNotAvailable):
		return s.errors.LessonNotAvailable(ctx), nil
	case ent.IsNotFound(err):
		return s.lessonNotFound(ctx), nil
	case err != nil:
		// Anything else is a server fault, reported by the central handler.
		return nil, err
	}

	progressView := apiconv.ToCourseProgress(state)
	return &progressView, nil
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
) (api.CheckLessonRes, error) {
	result, err := s.progress.CheckSolution(ctx, progress.Check{
		Learner:   currentLearner(ctx),
		LessonID:  int(params.ID),
		VersionID: int(req.VersionId),
		Code:      req.Code,
		Locale:    s.i18n.Locale(ctx),
	})
	switch {
	case errors.Is(err, progress.ErrLessonNotAvailable):
		// Refused before the submission ran: nothing was executed and nothing
		// was written.
		return s.errors.LessonNotAvailable(ctx), nil
	case ent.IsNotFound(err):
		// The lesson, or the version the code was written against, is not part
		// of the course's current build.
		return s.lessonNotFound(ctx), nil
	case err != nil:
		return nil, err
	}

	response := &api.LessonCheckingResponseHeaders{
		Response: api.LessonCheckingResponse{
			Passed:                result.Outcome.Passed,
			Output:                result.Outcome.Output,
			Result:                checkResultOf(result.Outcome),
			Status:                int32(result.Outcome.Status),
			LessonHasBeenFinished: result.LessonFinished,
			CourseHasBeenFinished: result.CourseFinished,
		},
	}
	if result.Guest != nil {
		cookie, err := s.auth.GuestCookie(*result.Guest)
		if err != nil {
			return nil, err
		}
		response.SetCookie = []string{cookie}
	}
	return response, nil
}

// lessonNotFound is the contract's typed 404 for both lesson operations. They
// declare it, so it has to be returned as such: the central handler's problem
// document carries a different content type, which a client generated from this
// contract cannot decode as this operation's 404.
func (s *Server) lessonNotFound(ctx context.Context) *api.NotFoundError {
	return &api.NotFoundError{Message: s.i18n.Text(ctx, localization.LessonNotFound)}
}

// currentLearner is whose progress a request is about: the signed-in user when
// the session cookie identified one, and otherwise the visitor's cookie
// progress — empty for a first visit, which is exactly the position of someone
// who has finished nothing.
func currentLearner(ctx context.Context) progress.Learner {
	if u, ok := AuthenticatedUser(ctx); ok {
		return progress.Learner{UserID: u.ID}
	}
	guest, _ := GuestProgress(ctx)
	return progress.Learner{Guest: guest}
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
