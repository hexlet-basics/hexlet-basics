package exerciserunner

import (
	"context"
	"testing"
	"time"

	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/progress"
)

// runnerImage is a stand-in for an exercise image: the tests below need a
// container that starts, prints, and exits with a status of their choosing, not
// a real course. `sh` is the only thing they ask of it.
const runnerImage = "alpine:3"

// TestRunGradesAContainerItReallyStarts drives the whole Docker lifecycle —
// create, deliver the submission, attach, start, read, wait — against a real
// daemon. The rest of this package's tests are about the mapping; this is the
// one that would catch the sequencing being wrong, which is the part that
// cannot be reasoned about from the client's types alone.
//
// The submitted "code" is a shell script and the command is overridden to run
// it, because a real exercise image would mean pulling a course.
func TestRunGradesAContainerItReallyStarts(t *testing.T) {
	runner := dockerRunner(t)

	cases := []struct {
		name   string
		script string
		result string
		passed bool
		status int
		output string
	}{
		{
			name:   "a passing submission",
			script: "echo '1 test, 0 failures'; exit 0",
			result: progress.ResultPassed,
			passed: true,
			output: "1 test, 0 failures",
		},
		{
			name:   "a failing submission, with what it printed",
			script: "echo 'expected 3, got 2' >&2; exit 1",
			result: progress.ResultFailed,
			status: 1,
			output: "expected 3, got 2",
		},
		{
			// What GNU `timeout` reports when it kills a non-terminating program.
			name:   "a submission that never terminates",
			script: "echo 'looping'; exit 124",
			result: progress.ResultFailedInfinity,
			status: 124,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outcome, err := runScript(t, runner, tc.script)
			require.NoError(t, err)

			assert.Equal(t, tc.result, outcome.Result)
			assert.Equal(t, tc.passed, outcome.Passed)
			assert.Equal(t, tc.status, outcome.Status)
			if tc.output == "" {
				assert.Empty(t, outcome.Output)
			} else {
				assert.Contains(t, outcome.Output, tc.output)
			}
		})
	}
}

// The submitted code is delivered into the container rather than mounted, so it
// reaches a daemon that is not local. This proves the file arrives with the
// bytes the learner submitted, at the path the tests read.
func TestRunDeliversTheSubmittedCode(t *testing.T) {
	runner := dockerRunner(t)

	outcome, err := runScript(t, runner, "cat /submission/solution.txt")
	require.NoError(t, err)

	assert.True(t, outcome.Passed)
	assert.Contains(t, outcome.Output, "the learner's code")
}

// Output is bounded: a program that prints continuously right up to the
// deadline must not be able to hand the process megabytes to hold.
func TestRunCapsRunawayOutput(t *testing.T) {
	runner := dockerRunner(t)
	runner.opts.MaxOutputBytes = 512

	outcome, err := runScript(t, runner, "yes 'noise' | head -c 100000")
	require.NoError(t, err)

	assert.LessOrEqual(t, len(outcome.Output), 512)
}

// dockerRunner builds a runner against the local daemon, and skips the test
// when there is none — the rest of the suite needs Docker for its database
// containers, so this only skips where Docker is genuinely absent.
func dockerRunner(t *testing.T) *Docker {
	t.Helper()

	runner, err := New(OptionsFrom(testConfig()))
	if err != nil {
		t.Skipf("no docker client: %v", err)
	}
	t.Cleanup(func() { _ = runner.Close() })

	if _, err := runner.docker.Ping(context.Background(), client.PingOptions{}); err != nil {
		t.Skipf("no docker daemon: %v", err)
	}
	pullImage(t, runner)
	return runner
}

func pullImage(t *testing.T, runner *Docker) {
	t.Helper()
	ctx := t.Context()

	if _, err := runner.docker.ImageInspect(ctx, runnerImage); err == nil {
		return
	}
	pulled, err := runner.docker.ImagePull(ctx, runnerImage, client.ImagePullOptions{})
	if err != nil {
		// A daemon with no registry access is an environment without the image,
		// not a failing runner.
		t.Skipf("cannot pull %s: %v", runnerImage, err)
	}
	defer func() { _ = pulled.Close() }()
	// The pull returns before it finishes; starting a container against a
	// half-present image is the classic way to get a confusing failure.
	if err := pulled.Wait(ctx); err != nil {
		t.Skipf("cannot pull %s: %v", runnerImage, err)
	}
}

// runScript drives the real container lifecycle with a shell script in place of
// a course's Makefile: the command a real image answers is
// `timeout N make -C <dir> test`, and alpine has neither make nor a course.
// What is under test here is the lifecycle and the mapping, not the command,
// which has a test of its own.
func runScript(t *testing.T, runner *Docker, script string) (progress.Outcome, error) {
	t.Helper()

	submission := progress.Submission{
		Image:        runnerImage,
		TestDir:      "/submission",
		ExerciseFile: "/submission/solution.txt",
		Code:         "the learner's code\n",
	}
	status, output, err := runner.run(t.Context(), submission, []string{"sh", "-c", script})
	if err != nil {
		return progress.Outcome{}, err
	}
	return outcomeOf(status, output), nil
}

// A container that never exits must not hold the process: the outer deadline is
// the backstop for an image whose `timeout` is missing or ignored.
func TestRunGivesUpOnAContainerThatNeverReports(t *testing.T) {
	runner := dockerRunner(t)
	// Comfortably more than the daemon needs to create and start a container,
	// and far short of the sleep below.
	runner.opts.Timeout = 5 * time.Second
	runner.opts.GraceTimeout = 10 * time.Second

	done := make(chan error, 1)
	go func() {
		_, err := runScript(t, runner, "sleep 60")
		done <- err
	}()

	select {
	case err := <-done:
		require.Error(t, err)
		// Whichever call the budget runs out on, the run ends rather than
		// holding its concurrency slot until the process dies.
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(45 * time.Second):
		t.Fatal("the read never gave up: a container that never exits wedges the runner")
	}
}

// The pin is created by pulling the moving tag once and freezing it, so a
// course whose image has not been graded on this daemon before still runs — and
// the next check for that version finds the pin rather than pulling again.
func TestEnsureImagePinsToTheCourseVersion(t *testing.T) {
	runner := dockerRunner(t)
	runner.opts.ImageTag = "3"
	submission := progress.Submission{Image: "alpine", VersionID: 20260812}
	pinned, _ := runner.imageRefs(submission)
	t.Cleanup(func() {
		_, _ = runner.docker.ImageRemove(context.Background(), pinned, client.ImageRemoveOptions{})
	})

	resolved, err := runner.ensureImage(t.Context(), submission)
	require.NoError(t, err)
	assert.Equal(t, pinned, resolved)

	_, err = runner.docker.ImageInspect(t.Context(), pinned)
	assert.NoError(t, err, "the pin exists on the daemon now")

	// Again, from the pin this time: the image is already frozen, so nothing is
	// pulled and the answer is the same.
	resolved, err = runner.ensureImage(t.Context(), submission)
	require.NoError(t, err)
	assert.Equal(t, pinned, resolved)
}
