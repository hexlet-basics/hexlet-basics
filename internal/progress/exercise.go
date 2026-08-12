package progress

import (
	"context"
	"errors"
)

// Result classifications, as the legacy runner produced them and as the
// contract still declares them. They are derived from the exit status of the
// process that ran the tests: 0 passes, 124 is the status the timeout kills a
// non-terminating submission with, and anything else is an ordinary failure.
const (
	ResultPassed         = "passed"
	ResultFailed         = "failed"
	ResultFailedInfinity = "failed-infinity"
)

// Submission is one solution to run. Code is only ever meaningful relative to a
// Lesson Version — the Version carries the tests it runs against, the image it
// runs in and the paths inside that image — so the pair travels together.
type Submission struct {
	LessonVersionID int
	CourseVersionID int

	// UserID is 0 for a guest. The legacy runner used it to name the scratch
	// directory it wrote the submitted code into, so who submitted stays part of
	// the request rather than something the runner has to ask for separately.
	UserID int

	Code string
}

// Outcome is what running one Submission produced: whether it passed, the
// captured output, the classification above, and the raw exit status.
type Outcome struct {
	Passed bool
	Output string
	Result string
	Status int
}

// ExerciseRunner runs a Submission against its Lesson Version's tests.
//
// It is an interface because the production implementation shells out to
// Docker, and a check is the trigger for every write in this module: without
// the seam, no progress transition could be tested without a container. The
// Docker implementation is separate work; the module defines and consumes the
// interface only.
type ExerciseRunner interface {
	Run(ctx context.Context, submission Submission) (Outcome, error)
}

// ErrRunnerUnavailable reports a process wired with no exercise runner.
var ErrRunnerUnavailable = errors.New("progress: no exercise runner is configured")

// UnavailableRunner stands in until the Docker runner lands. It fails loudly
// rather than silently reporting a pass or a failure: a check that cannot run
// has no outcome, and inventing one would write progress a learner did not earn.
type UnavailableRunner struct{}

// Run always fails with ErrRunnerUnavailable.
func (UnavailableRunner) Run(context.Context, Submission) (Outcome, error) {
	return Outcome{}, ErrRunnerUnavailable
}
