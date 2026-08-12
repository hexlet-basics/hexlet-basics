# ADR-0013: How untrusted exercise code is executed

**Status:** Accepted

## Context

Checking a solution means running code a stranger wrote, on our hardware, and
telling them whether it passed. It is the product's core loop and its largest
attack surface, and the Go rewrite reached parity on everything around it —
sequential progression, enrollment, the check operation itself (ADR-0012) —
with the run behind an interface and no implementation.

The legacy runner is one `docker run` line
(`legacy/app/lib/docker_exercise_client.rb:35`) plus an exit-code mapping
(`legacy/app/lib/lesson_tester.rb:43-54`). It sets a memory cap, disables the
network, and grants **unlimited swap**. It sets no process limit, no CPU cap,
drops no capabilities, and bind-mounts the submitted file from a host path —
which the file itself flags as a problem under Docker-in-Docker.

The research behind this decision is `docs/research/2026-08-04-untrusted-code-runner-in-go.md`,
which cites the Docker client's own source for every API claim.

## Decision

**The client is `github.com/moby/moby/client`.** It is already in the build
graph through testcontainers, and `github.com/docker/docker` is absent
entirely, so promoting it to a direct dependency adds no new transitive weight.
Its API is options-struct based and is not the positional API of the older
client.

**The check stays synchronous.** The contract models one blocking POST, and
that is what ships. The pain that makes a blocking six-second call bad in Rails
is thread-pool exhaustion, which Go's goroutine-per-request model does not
have; the real ceiling is Docker daemon capacity, which is identical whether
the run sits in the request or behind a job. It is bounded by a semaphore
inside the runner instead. Moving to submit-plus-stream later is a contract
change worth making only if incremental output during a run becomes a feature,
not as a way to avoid blocking.

**The in-container `timeout` stays authoritative.** GNU coreutils documents
`timeout` as exiting 124 when it kills the command, and the whole classification
rests on it: 0 passes, 124 is `failed-infinity`, anything else failed — byte for
byte what legacy returned. An outer context deadline at the budget plus a grace
period is a *backstop*, and a container that trips it reports an infrastructure
error rather than a failed check: the outer clock can say "I gave up", never
why the process stopped. This makes a working `timeout` an image contract; an
image without one would silently turn timeouts into ordinary failures.

**The submission is copied in, not mounted.** A bind mount ties the runner to
host paths and cannot work against a daemon that is not local. A one-file tar
into the created-but-unstarted container lands on exactly the file the
reference solution ships, which is how the tests find it.

**Output is captured by attaching before start, and bounded.** With
`AutoRemove` the container can be gone before anything could read its logs, so
the attach and the wait are both issued before the start. The stream is not a
TTY, so stdout and stderr stay distinguishable and are demultiplexed with
`stdcopy` — then merged, because a learner reads one console. A byte ceiling
lives in the writer rather than in a limited reader, since cutting the
multiplexed stream mid-frame would corrupt the demultiplexing; the writer keeps
draining past the cap so a chatty program cannot block on a full pipe. Invalid
UTF-8 is scrubbed rather than rejected, which is what legacy's escape-plus-base64
was for. A timed-out run returns no output at all, as legacy did: a spinning
loop prints megabytes of one line and the classification already says
everything.

**Isolation is tightened past legacy, and every limit is configurable.** The
baseline adds what legacy never set — `PidsLimit`, `NanoCPUs`, `CapDrop: ALL`,
`no-new-privileges` — and **drops unlimited swap**: the ceiling now equals the
memory limit, so there is a limit rather than the absence of one. `PidsLimit` is
the one that matters most: neither a memory cap nor an isolated network does
anything against a fork bomb, which is a one-line submission in most of these
languages.

Two limits are deliberately **off by default**: a read-only root filesystem and
a non-root user. Compiled-language exercises write build artifacts beside their
sources, so both can turn a currently-green course red, and neither can be
judged from here. They ship as configuration so a canary — every course's
reference solutions run under the new limits — can turn them on without a
release. Every other limit is configurable for the same reason.

**The daemon comes from the environment.** The client resolves `DOCKER_HOST`
and friends, and connects lazily. Where the daemon lives — a worker-local
socket, a remote daemon, docker-in-docker — stays a deployment decision, and a
deployment with none reachable fails on the first check rather than refusing to
boot, so every other surface keeps serving. Note that mounting the host's
Docker socket into a process is effectively granting it host root; that is an
infrastructure decision this ADR deliberately does not make.

**Images are expected to be present.** The runner runs the image named by the
Course Version's `docker_image` and does not pull. A first-check-after-deploy
pull is seconds to minutes inside a six-second grading budget, which reads to a
learner as a broken lesson. Pre-pulling on version promotion and warming
current versions at worker startup is the follow-up this leaves open.

## Consequences

- The runner is the **only** Docker consumer in the Go app: content loading
  moved to git (`internal/courseloader/fetch.go`). That makes the daemon a
  single isolable dependency, and also means nothing else proves daemon access
  works in a given deployment.
- Grading behaviour is unchanged for a passing or failing submission, and
  unchanged for a timeout, so the funnel's numbers carry across the cutover.
- Any exercise that relied on swap beyond its memory limit, on more than the
  configured processes, or on more CPU than the cap will now fail where it used
  to pass. That is the intended direction, and the reason every limit is a
  configuration value rather than a constant.
- Running is separated from grading: the progress module decides whether a
  submission may run and what its outcome means, this decides how it runs. An
  error from the runner therefore never means "the tests failed" — it means the
  submission could not be run at all, and it surfaces as a server error rather
  than as a red check.
- The `timeout` binary becomes a documented expectation of every exercise image.
