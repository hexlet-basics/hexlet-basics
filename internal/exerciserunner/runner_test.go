package exerciserunner

import (
	"archive/tar"
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/config"
	"hexletbasics/internal/progress"
)

// The exit status is the whole grading contract with the exercise images, and
// it is the one thing that must not drift from legacy: 0 passes, 124 is the
// in-container timeout, everything else is an ordinary failure.
func TestOutcomeOfMapsTheExitStatus(t *testing.T) {
	passed := outcomeOf(0, "1 test, 0 failures")
	assert.True(t, passed.Passed)
	assert.Equal(t, progress.ResultPassed, passed.Result)
	assert.Equal(t, "1 test, 0 failures", passed.Output)

	failed := outcomeOf(1, "expected 3, got 2")
	assert.False(t, failed.Passed)
	assert.Equal(t, progress.ResultFailed, failed.Result)
	assert.Equal(t, "expected 3, got 2", failed.Output, "a learner needs to read why it failed")
	assert.Equal(t, 1, failed.Status)

	// 127 is "command not found" — an image without `make`, say. It is still an
	// ordinary failure rather than anything the learner can distinguish.
	assert.Equal(t, progress.ResultFailed, outcomeOf(127, "").Result)
}

// A timed-out run returns no output: a spinning loop prints megabytes of the
// same line, and the classification already says everything.
func TestOutcomeOfDropsTheOutputOfATimedOutRun(t *testing.T) {
	timedOut := outcomeOf(timeoutExitStatus, strings.Repeat("looping\n", 1000))

	assert.False(t, timedOut.Passed)
	assert.Equal(t, progress.ResultFailedInfinity, timedOut.Result)
	assert.Empty(t, timedOut.Output)
	assert.Equal(t, timeoutExitStatus, timedOut.Status)
}

// Output is capped, and the cap is applied while still draining the stream: a
// reader that stopped early would block the container on a full pipe instead of
// letting it finish and report its status.
func TestReadOutputCapsWithoutStalling(t *testing.T) {
	stream := multiplexed(t, strings.Repeat("a", 5000), "")

	output, err := readOutput(stream, 100)
	require.NoError(t, err)
	assert.Len(t, output, 100)
}

// stdout and stderr are merged: a learner reads one console, and the test
// harnesses write their failures to whichever they please.
func TestReadOutputMergesBothStreams(t *testing.T) {
	stream := multiplexed(t, "on stdout\n", "on stderr\n")

	output, err := readOutput(stream, 1024)
	require.NoError(t, err)
	assert.Contains(t, output, "on stdout")
	assert.Contains(t, output, "on stderr")
}

// Output is whatever the learner's program printed, so a stray byte must not be
// able to break the response encoding.
func TestReadOutputScrubsInvalidUTF8(t *testing.T) {
	stream := multiplexed(t, "ok \xff\xfe done", "")

	output, err := readOutput(stream, 1024)
	require.NoError(t, err)
	assert.Equal(t, "ok  done", output)
}

// The command is the legacy one, unchanged: the lesson's Makefile under a
// timeout, which is the interface every exercise image already exposes.
func TestCommandRunsTheLessonsMakefileUnderTheTimeout(t *testing.T) {
	runner := &Docker{opts: OptionsFrom(testConfig())}

	command := runner.command(progress.Submission{TestDir: "/exercises-go/modules/10-basics/20-vars"})

	assert.Equal(t,
		[]string{"timeout", "6", "make", "--silent", "-C", "/exercises-go/modules/10-basics/20-vars", "test"},
		command)
}

// The submission is delivered as a tar rooted at the container filesystem, so
// it lands exactly on the file the reference solution ships.
func TestTarFileTargetsTheExerciseFile(t *testing.T) {
	archive, err := tarFile("/exercises-go/modules/10-basics/20-vars/index.go", "package main\n")
	require.NoError(t, err)

	reader := tar.NewReader(archive)
	header, err := reader.Next()
	require.NoError(t, err)

	assert.Equal(t, "exercises-go/modules/10-basics/20-vars/index.go", header.Name,
		"relative to the copy destination, which is the container root")

	content, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "package main\n", string(content))
}

// The hardened limits are the point of this package: a memory cap without a
// pids limit does nothing against a fork bomb, which is a one-line submission
// in most of these languages.
func TestHostConfigCarriesTheHardenedLimits(t *testing.T) {
	runner := &Docker{opts: OptionsFrom(testConfig())}

	host := runner.hostConfig()

	assert.Equal(t, "none", string(host.NetworkMode), "untrusted code gets no network")
	assert.True(t, host.AutoRemove)
	assert.Equal(t, []string{"ALL"}, host.CapDrop)
	assert.Equal(t, []string{"no-new-privileges"}, host.SecurityOpt)
	require.NotNil(t, host.PidsLimit)
	assert.Positive(t, *host.PidsLimit)
	assert.Positive(t, host.NanoCPUs)
	assert.Equal(t, host.Memory, host.MemorySwap,
		"swap equal to memory means no swap; legacy's -1 meant unlimited")
}

// Read-only rootfs is off by default — compiled languages write build artifacts
// beside their sources — but when a deployment turns it on, something has to
// stay writable or nothing compiles at all.
func TestReadonlyRootfsBringsWritableScratch(t *testing.T) {
	cfg := testConfig()
	cfg.ReadonlyRootfs = true
	runner := &Docker{opts: OptionsFrom(cfg)}

	host := runner.hostConfig()

	assert.True(t, host.ReadonlyRootfs)
	assert.Contains(t, host.Tmpfs, "/tmp")
	assert.Empty(t, (&Docker{opts: OptionsFrom(testConfig())}).hostConfig().Tmpfs,
		"and no scratch mount when the rootfs is writable anyway")
}

// multiplexed frames content the way the Docker daemon frames an attached
// stream — an 8-byte header naming the stream and the payload length, then the
// payload — so the demultiplexing under test is the real one.
func multiplexed(t *testing.T, stdout, stderr string) io.Reader {
	t.Helper()
	var buf bytes.Buffer

	for _, frame := range []struct {
		stream  stdcopy.StdType
		payload string
	}{{stdcopy.Stdout, stdout}, {stdcopy.Stderr, stderr}} {
		if frame.payload == "" {
			continue
		}
		header := make([]byte, 8)
		header[0] = byte(frame.stream)
		binary.BigEndian.PutUint32(header[4:], uint32(len(frame.payload)))
		buf.Write(header)
		buf.WriteString(frame.payload)
	}
	return &buf
}

// testConfig is the shipped default budget, decoded from its own env tags, so
// these tests assert what a deployment actually runs under rather than values
// invented here.
func testConfig() config.ExerciseRunnerConfig {
	var cfg config.ExerciseRunnerConfig
	// Decoded against an empty environment on purpose: what is asserted is the
	// shipped default, not whatever the machine running the tests exports.
	if err := env.ParseWithOptions(&cfg, env.Options{Environment: map[string]string{}}); err != nil {
		panic(err)
	}
	return cfg
}

// Grading is pinned to the course version that produced the image: a rebuild
// moves the registry tag, and a learner working through an already-promoted
// version must keep being graded by the build it was promoted as.
func TestImageRefsPinToTheCourseVersion(t *testing.T) {
	runner := &Docker{opts: OptionsFrom(testConfig())}

	pinned, source := runner.imageRefs(progress.Submission{
		Image:     "hexletbasics/exercises-javascript",
		VersionID: 965227298,
	})

	assert.Equal(t, "hexletbasics/exercises-javascript:lv965227298", pinned)
	assert.Equal(t, "hexletbasics/exercises-javascript:latest", source,
		"the moving tag is only ever the pull source")
}

// A course version that already named a tag has said what it wants; that is its
// own pin, and nothing is pulled or re-tagged behind it.
func TestImageRefsRespectAnExplicitTag(t *testing.T) {
	runner := &Docker{opts: OptionsFrom(testConfig())}

	pinned, source := runner.imageRefs(progress.Submission{
		Image:     "ghcr.io/hexlet-basics/exercises-go:2026-08",
		VersionID: 1,
	})

	assert.Equal(t, "ghcr.io/hexlet-basics/exercises-go:2026-08", pinned)
	assert.Equal(t, pinned, source)
}

// A registry host carries a port, which is not a tag.
func TestImageRefsPinAnImageOnAPortedRegistry(t *testing.T) {
	runner := &Docker{opts: OptionsFrom(testConfig())}

	pinned, _ := runner.imageRefs(progress.Submission{
		Image:     "registry.example.com:5000/exercises-go",
		VersionID: 7,
	})

	assert.Equal(t, "registry.example.com:5000/exercises-go:lv7", pinned)
}
