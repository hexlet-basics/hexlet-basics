// Package exerciserunner runs a learner's submission against a Lesson's tests
// in a Docker container, and is the only Docker-daemon consumer in the Go app
// (ADR-0013).
//
// It implements progress.ExerciseRunner. The progress module decides whether a
// submission may run and what its outcome means; everything here is about
// running untrusted code and getting its output back without letting it harm
// the host.
package exerciserunner

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"

	"hexletbasics/internal/config"
	"hexletbasics/internal/progress"
)

// timeoutExitStatus is what GNU coreutils' `timeout` exits with when it kills
// the command: "124 if command times out, and --preserve-status is not
// specified". It is the whole basis for telling an infinite loop apart from a
// failing test, which is why the timeout runs INSIDE the container: the outer
// clock can only say "I gave up", never why the process stopped.
//
// It is an undocumented contract with the exercise images — a future minimal
// image without `timeout` would silently turn timeouts into ordinary failures.
const timeoutExitStatus = 124

// Options are the knobs the deployment can turn. Their defaults tighten what
// the legacy runner allowed, so any one of them can be relaxed from config if a
// course turns out to need it (ADR-0013).
type Options struct {
	// Timeout is the budget given to the in-container `timeout`, and is what
	// classifies a submission as an infinite loop.
	Timeout time.Duration

	// GraceTimeout is added to Timeout for the outer deadline. It fires only
	// when the container never reported at all, which is a daemon problem rather
	// than a learner's, so it surfaces as an error, not as a failed check.
	GraceTimeout time.Duration

	MemoryBytes int64

	// SwapBytes is the container's memory+swap ceiling. Equal to MemoryBytes
	// means no swap at all; -1 restores the legacy behaviour of unlimited swap,
	// which is not a limit but the absence of one.
	SwapBytes int64

	// PidsLimit caps processes. Legacy set none, so a fork bomb — a one-line
	// submission in most of these languages — was bounded by nothing.
	PidsLimit int64

	// NanoCPUs caps CPU, in billionths of a core.
	NanoCPUs int64

	// ReadonlyRootfs is off by default: compiled-language exercises write build
	// artifacts next to their sources, so turning it on needs a run of every
	// course's reference solutions first.
	ReadonlyRootfs bool

	// User runs the command as somebody other than the image's own user. Empty
	// keeps the image's, for the same reason ReadonlyRootfs is off.
	User string

	// MaxOutputBytes caps captured output. Legacy had no cap, so a program
	// printing continuously right up to the timeout returned everything it
	// printed.
	MaxOutputBytes int

	// Concurrency bounds how many submissions run at once, because the real
	// ceiling is the daemon's capacity rather than this process's.
	Concurrency int

	// ImageTag is the moving registry tag a course publishes its latest build
	// under. It is pulled once per Course Version and then frozen under a pinned
	// reference, so it decides what a NEW version is graded by, never what an
	// existing one is graded by.
	ImageTag string
}

// OptionsFrom adapts the deployment's configuration. The two representations
// are kept apart on purpose: config speaks environment variables and seconds,
// this package speaks durations and Docker fields.
func OptionsFrom(cfg config.ExerciseRunnerConfig) Options {
	return Options{
		Timeout:        time.Duration(cfg.TimeoutSeconds) * time.Second,
		GraceTimeout:   time.Duration(cfg.GraceSeconds) * time.Second,
		MemoryBytes:    cfg.MemoryBytes,
		SwapBytes:      cfg.SwapBytes,
		PidsLimit:      cfg.PidsLimit,
		NanoCPUs:       cfg.NanoCPUs,
		ReadonlyRootfs: cfg.ReadonlyRootfs,
		User:           cfg.RunAsUser,
		MaxOutputBytes: cfg.MaxOutputBytes,
		Concurrency:    cfg.Concurrency,
		ImageTag:       cfg.ImageTag,
	}
}

// Docker is the production runner.
type Docker struct {
	docker *client.Client
	opts   Options

	// slots bounds concurrent runs. A submission that cannot get one waits
	// rather than being refused: the learner already pressed the button, and a
	// queue of a few seconds beats an error they cannot act on.
	slots chan struct{}
}

// New builds the runner over a Docker client. The client resolves the daemon
// from the environment (DOCKER_HOST and friends), so where the daemon lives is
// a deployment decision rather than a code one, and it connects lazily — an
// unreachable daemon surfaces on the first check, not at boot.
func New(opts Options) (*Docker, error) {
	docker, err := client.New(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("new docker client: %w", err)
	}
	if opts.Concurrency < 1 {
		opts.Concurrency = 1
	}
	return &Docker{
		docker: docker,
		opts:   opts,
		slots:  make(chan struct{}, opts.Concurrency),
	}, nil
}

// Close releases the daemon connection.
func (d *Docker) Close() error { return d.docker.Close() }

// Run grades one submission and returns what the container reported.
//
// An error means the submission could not be run at all — an unreachable
// daemon, a missing image, a container that never reported. It is never used to
// report a failing test: a failure is an outcome, and the learner sees it.
func (d *Docker) Run(ctx context.Context, submission progress.Submission) (progress.Outcome, error) {
	select {
	case d.slots <- struct{}{}:
		defer func() { <-d.slots }()
	case <-ctx.Done():
		return progress.Outcome{}, ctx.Err()
	}

	image, err := d.ensureImage(ctx, submission)
	if err != nil {
		return progress.Outcome{}, err
	}

	status, output, err := d.run(ctx, submission.WithImage(image), d.command(submission))
	if err != nil {
		return progress.Outcome{}, err
	}
	return outcomeOf(status, output), nil
}

// ensureImage resolves the image this submission is graded in, pinned to the
// Course Version that produced it.
//
// The pin is what keeps grading honest across a rebuild: the registry tag a
// course publishes moves, and a learner working through a promoted version must
// keep being graded by the build that version was promoted as — not by whatever
// was pushed since. Legacy did exactly this, and this is the port of it: look
// for the pinned reference, and when it is missing pull the moving tag once and
// freeze it under the pin.
//
// The freeze is local to the daemon, so the first check after a promotion pays
// for the pull. Pre-pulling on promotion is the follow-up that removes that
// cost; correctness does not depend on it.
func (d *Docker) ensureImage(ctx context.Context, submission progress.Submission) (string, error) {
	pinned, source := d.imageRefs(submission)
	if pinned == source {
		// The Course Version named an explicit tag; it is the pin.
		return pinned, nil
	}

	if _, err := d.docker.ImageInspect(ctx, pinned); err == nil {
		return pinned, nil
	}

	pulled, err := d.docker.ImagePull(ctx, source, client.ImagePullOptions{})
	if err != nil {
		return "", fmt.Errorf("pull %q: %w", source, err)
	}
	defer func() { _ = pulled.Close() }()
	// The pull returns as soon as the daemon accepts it; starting a container
	// against a half-present image is the confusing failure that follows from
	// not waiting.
	if err := pulled.Wait(ctx); err != nil {
		return "", fmt.Errorf("pull %q: %w", source, err)
	}

	if _, err := d.docker.ImageTag(ctx, client.ImageTagOptions{Source: source, Target: pinned}); err != nil {
		return "", fmt.Errorf("pin %q as %q: %w", source, pinned, err)
	}
	return pinned, nil
}

// imageRefs are the two references a submission needs: the pinned one it is
// graded in, and the moving one that is pulled to create it. They are equal
// when the Course Version already named a tag, which is its own pin.
func (d *Docker) imageRefs(submission progress.Submission) (pinned, source string) {
	repository := submission.Image
	if strings.Contains(path.Base(repository), ":") {
		return repository, repository
	}
	return fmt.Sprintf("%s:lv%d", repository, submission.VersionID),
		fmt.Sprintf("%s:%s", repository, d.opts.ImageTag)
}

// outcomeOf maps an exit status onto the contract's classification, exactly as
// the legacy runner did: 0 passes, 124 is the timeout, anything else failed.
//
// The output of a timed-out run is dropped rather than returned. A spinning
// loop's output is megabytes of the same line, and it tells the learner nothing
// the classification does not.
func outcomeOf(status int, output string) progress.Outcome {
	switch status {
	case 0:
		return progress.Outcome{Passed: true, Result: progress.ResultPassed, Output: output, Status: status}
	case timeoutExitStatus:
		return progress.Outcome{Result: progress.ResultFailedInfinity, Output: "", Status: status}
	default:
		return progress.Outcome{Result: progress.ResultFailed, Output: output, Status: status}
	}
}

// run performs the container lifecycle: create, deliver the code, attach, wait,
// start, read.
//
// The order matters. The wait is issued BEFORE the start so a container that
// exits immediately cannot be gone before anyone asked about it, and the attach
// happens before the start for the same reason — with AutoRemove there is no
// container left to read logs from afterwards.
//
// The command is a parameter rather than derived here so the lifecycle can be
// exercised against an image that is not a course.
func (d *Docker) run(ctx context.Context, submission progress.Submission, command []string) (int, string, error) {
	// The deadline bounds the whole lifecycle, not just one call in it: it is the
	// backstop for a container that never reports at all, while the in-container
	// timeout is what decides an infinite loop.
	ctx, cancel := context.WithTimeout(ctx, d.opts.Timeout+d.opts.GraceTimeout)
	defer cancel()

	created, err := d.docker.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image: submission.Image,
			Cmd:   command,
			User:  d.opts.User,
			// Not a TTY, so stdout and stderr stay distinguishable in the stream.
			Tty:             false,
			AttachStdout:    true,
			AttachStderr:    true,
			NetworkDisabled: true,
			Labels: map[string]string{
				"com.hexletbasics.role": "exercise-check",
				"com.hexletbasics.user": strconv.Itoa(submission.UserID),
			},
		},
		HostConfig: d.hostConfig(),
	})
	if err != nil {
		return 0, "", fmt.Errorf("create container for %q: %w", submission.Image, err)
	}

	// AutoRemove cleans up after a normal exit; this covers the paths where the
	// container never starts, and is harmless once it is already gone.
	defer func() {
		_, _ = d.docker.ContainerRemove(context.WithoutCancel(ctx), created.ID, client.ContainerRemoveOptions{Force: true})
	}()

	archive, err := tarFile(submission.ExerciseFile, submission.Code)
	if err != nil {
		return 0, "", err
	}
	// Copied in rather than bind-mounted: a mount ties the runner to host paths
	// and cannot work against a daemon that is not local.
	if _, err := d.docker.CopyToContainer(ctx, created.ID, client.CopyToContainerOptions{
		DestinationPath: "/",
		Content:         archive,
	}); err != nil {
		return 0, "", fmt.Errorf("deliver submission to container: %w", err)
	}

	wait := d.docker.ContainerWait(ctx, created.ID, client.ContainerWaitOptions{
		Condition: container.WaitConditionNextExit,
	})

	attached, err := d.docker.ContainerAttach(ctx, created.ID, client.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return 0, "", fmt.Errorf("attach to container: %w", err)
	}
	defer attached.Close()

	if _, err := d.docker.ContainerStart(ctx, created.ID, client.ContainerStartOptions{}); err != nil {
		return 0, "", fmt.Errorf("start container: %w", err)
	}

	// Read on its own goroutine. A hijacked connection does not honour context
	// cancellation, so reading inline would mean a container that never exits —
	// an image whose `timeout` is missing, a process that ignores the signal —
	// blocks here forever: the deadline below would never be reached, the
	// container would never be removed, and the concurrency slot would never be
	// released. A few of those and every check hangs.
	//
	// The deferred close and force-remove are what end the read once this
	// function returns, so the goroutine cannot outlive the run for long.
	reading := make(chan readResult, 1)
	go func() {
		output, err := readOutput(attached.Reader, d.opts.MaxOutputBytes)
		reading <- readResult{output: output, err: err}
	}()

	var output string
	select {
	case read := <-reading:
		if read.err != nil {
			return 0, "", read.err
		}
		output = read.output
	case <-ctx.Done():
		return 0, "", d.gaveUp(ctx)
	}

	select {
	case result := <-wait.Result:
		if result.Error != nil && result.Error.Message != "" {
			return 0, "", fmt.Errorf("container wait failed: %s", result.Error.Message)
		}
		return int(result.StatusCode), output, nil
	case err := <-wait.Error:
		return 0, "", fmt.Errorf("wait for container: %w", err)
	case <-ctx.Done():
		return 0, "", d.gaveUp(ctx)
	}
}

// readResult carries what the output goroutine produced.
type readResult struct {
	output string
	err    error
}

// gaveUp reports a container that never finished within the budget plus its
// grace. That is infrastructure failing, not a learner writing a slow program —
// the in-container timeout classifies that one, and this cannot tell why the
// process stopped, only that nothing was heard.
func (d *Docker) gaveUp(ctx context.Context) error {
	return fmt.Errorf("container did not report within %s: %w",
		d.opts.Timeout+d.opts.GraceTimeout, ctx.Err())
}

// command is the legacy command, unchanged: run the lesson's Makefile under a
// timeout. `make` is the interface every exercise image already exposes.
func (d *Docker) command(submission progress.Submission) []string {
	seconds := strconv.Itoa(int(d.opts.Timeout.Seconds()))
	return []string{"timeout", seconds, "make", "--silent", "-C", submission.TestDir, "test"}
}

func (d *Docker) hostConfig() *container.HostConfig {
	// Copied, because the field is a pointer and the options are shared by every
	// run this runner performs.
	pidsLimit := d.opts.PidsLimit
	tmpfs := map[string]string{}
	if d.opts.ReadonlyRootfs {
		// Something has to be writable, or nothing compiles at all.
		tmpfs["/tmp"] = "rw,exec,size=64m"
	}

	return &container.HostConfig{
		AutoRemove:     true,
		NetworkMode:    "none",
		ReadonlyRootfs: d.opts.ReadonlyRootfs,
		Tmpfs:          tmpfs,
		CapDrop:        []string{"ALL"},
		SecurityOpt:    []string{"no-new-privileges"},
		Resources: container.Resources{
			Memory:     d.opts.MemoryBytes,
			MemorySwap: d.opts.SwapBytes,
			PidsLimit:  &pidsLimit,
			NanoCPUs:   d.opts.NanoCPUs,
		},
	}
}

// tarFile wraps the submitted code as the tar archive CopyToContainer takes,
// with the destination path baked in so the copy lands exactly on the file the
// reference solution ships.
func tarFile(destination, content string) (io.Reader, error) {
	var buf bytes.Buffer
	archive := tar.NewWriter(&buf)

	if err := archive.WriteHeader(&tar.Header{
		Name: strings.TrimPrefix(path.Clean(destination), "/"),
		Mode: 0o644,
		Size: int64(len(content)),
	}); err != nil {
		return nil, fmt.Errorf("write submission header: %w", err)
	}
	if _, err := archive.Write([]byte(content)); err != nil {
		return nil, fmt.Errorf("write submission: %w", err)
	}
	if err := archive.Close(); err != nil {
		return nil, fmt.Errorf("close submission archive: %w", err)
	}
	return &buf, nil
}

// readOutput demultiplexes the attached stream and returns what the exercise
// printed, capped and safe to put in JSON.
//
// stdout and stderr are merged deliberately: a learner reads one console, and
// the tests write their failures to whichever they please.
func readOutput(stream io.Reader, limit int) (string, error) {
	capped := &cappedWriter{limit: limit}
	if _, err := stdcopy.StdCopy(capped, capped, stream); err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read container output: %w", err)
	}

	// Scrubbed rather than rejected: exercise output is whatever the learner's
	// program printed, and a stray byte must not fail the response encoding.
	// Legacy escaped and base64-encoded it for the same reason.
	return strings.ToValidUTF8(capped.buf.String(), ""), nil
}

// cappedWriter keeps a bounded prefix of what the exercise printed. The cap has
// to live in the writer rather than in a limited reader: cutting the
// multiplexed stream mid-frame would corrupt the demultiplexing itself.
type cappedWriter struct {
	buf   bytes.Buffer
	limit int
}

func (w *cappedWriter) Write(p []byte) (int, error) {
	room := w.limit - w.buf.Len()
	if room > 0 {
		w.buf.Write(p[:min(room, len(p))])
	}
	// The full length is reported so the demultiplexer keeps draining: stopping
	// early would block the container on a full pipe rather than letting it
	// finish and report its status.
	return len(p), nil
}
