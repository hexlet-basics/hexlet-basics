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

// Submission is one solution to run, resolved down to what running it needs:
// the image the Course Version was built as, where the Lesson's tests live
// inside it, and the path the submitted code has to land at.
//
// It is resolved here rather than by the runner so the runner never touches the
// database — it is a Docker adapter, and which columns hold an image name is
// not its business.
type Submission struct {
	// Image is the exercise image, from the Course Version's docker_image.
	Image string

	// TestDir is the Lesson's directory inside that image; its Makefile is what
	// gets run. From the Lesson Version's path_to_code.
	TestDir string

	// ExerciseFile is where the submitted code must land — the Lesson's
	// directory plus the Course Version's exercise filename. Replacing the file
	// the reference solution ships is exactly how the tests find the submission.
	ExerciseFile string

	// UserID is 0 for a guest. Carried for logs and container labels, so a run
	// can be traced back to who asked for it.
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
