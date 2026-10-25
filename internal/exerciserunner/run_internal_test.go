//go:build docker

package exerciserunner

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/progress"
)

// Run is what the progress module calls. This drives it end to end against the
// real daemon with an image that is not a course: it has `timeout` but no
// `make`, so the grading command fails the way a broken exercise image would.
func TestRunReportsAnImageThatCannotGradeAsAFailedSubmission(t *testing.T) {
	t.Parallel()
	runner := dockerRunner(t)

	outcome, err := runner.Run(t.Context(), progress.Submission{
		// An explicit tag is its own pin, so nothing is pulled.
		Image:        runnerImage,
		TestDir:      "/submission",
		ExerciseFile: "/submission/solution.txt",
		Code:         "code\n",
		VersionID:    0,
		UserID:       0,
	})

	// A failing run is an outcome, never an error: the learner sees it.
	require.NoError(t, err)
	assert.False(t, outcome.Passed)
	assert.Equal(t, progress.ResultFailed, outcome.Result)
	assert.NotZero(t, outcome.Status)
}

// A submission that cannot get a run slot waits for one, and gives up when its
// caller does, instead of holding the request open behind a full queue.
func TestRunStopsWaitingForASlotWhenTheCallerGivesUp(t *testing.T) {
	t.Parallel()
	runner := dockerRunner(t)
	require.True(t, runner.slots.TryAcquire(1)) // the only slot is taken
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := runner.Run(ctx, progress.Submission{
		Image:        runnerImage,
		VersionID:    0,
		TestDir:      "",
		ExerciseFile: "",
		UserID:       0,
		Code:         "",
	})

	require.ErrorIs(t, err, context.Canceled)
}

// An image nobody can pull is an error, not a failing test: the submission was
// never run, and it must not be reported as the learner's mistake.
func TestRunErrorsWhenTheCourseImageCannotBePulled(t *testing.T) {
	t.Parallel()
	runner := dockerRunner(t)
	runner.opts.ImageTag = "latest"

	// Nothing listens on this registry address, so the pull cannot succeed and
	// the test needs no network. The daemon retries a refused registry for a
	// while; the caller's deadline is what ends the wait.
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()

	_, err := runner.Run(ctx, progress.Submission{
		Image:        "localhost:1/exercises-none",
		VersionID:    7,
		TestDir:      "",
		ExerciseFile: "",
		UserID:       0,
		Code:         "",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "pull")
}

func TestNewRunsOneSubmissionAtATimeUnlessToldOtherwise(t *testing.T) {
	t.Parallel()

	runner, err := New(Options{
		Timeout:        0,
		GraceTimeout:   0,
		MemoryBytes:    0,
		SwapBytes:      0,
		PidsLimit:      0,
		NanoCPUs:       0,
		ReadonlyRootfs: false,
		User:           "",
		MaxOutputBytes: 0,
		Concurrency:    0,
		ImageTag:       "",
	})
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, runner.Close()) })

	assert.True(t, runner.slots.TryAcquire(1))
	assert.False(t, runner.slots.TryAcquire(1), "the second run has to wait")
}
