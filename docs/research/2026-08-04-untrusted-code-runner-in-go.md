# Running untrusted exercise code in Docker from Go

**Question.** How should the Go rewrite run untrusted user code in Docker from a
river worker, replacing the legacy synchronous `POST /api/lessons/:id/check`?

**Date.** 2026-08-04. **Status.** Research only — no code written. Needs an ADR
(see [Decisions that need a human](#decisions-that-need-a-human)).

**Sources.** Docker/moby claims are cited to the module source already in this
repo's build graph (`$GOMODCACHE/...`, which is the published module content) or
to docs.docker.com. Go library claims are cited to the vendored source.
No secondary write-ups were used as the basis of any claim.

Legacy baseline being replaced:
`legacy/app/lib/docker_exercise_client.rb:35` (the `docker run` line) and
`legacy/app/lib/lesson_tester.rb:43-54` (exit-code → result mapping).

---

## Recommendation

1. **Client: `github.com/moby/moby/client`.** It is *already* in this repo's
   build graph at v0.5.1 (pulled by testcontainers-go), and
   `github.com/docker/docker` is absent entirely — so this choice adds zero new
   transitive dependency weight. Its API is options-struct based and differs from
   the older `docker/docker/client`; don't copy signatures from older examples.
2. **Keep the in-container `timeout` and the exit-124 mapping.** It is the only
   signal that distinguishes "infinite loop" from "tests failed" *without*
   racing your own wall clock, and it preserves the legacy contract byte for
   byte. Add an outer `context` deadline + kill as a backstop for a container
   that never reports, not as the primary mechanism.
3. **Deliver the submitted code with `CopyToContainer`, not a bind mount.** The
   legacy bind mount couples the runner to host paths and already carries a
   `# FIXME docker in docker volume` (`docker_exercise_client.rb:26`). A tar
   `CopyToContainer` into the created-but-not-started container works
   identically against a remote daemon.
4. **Attach before start; demux with `stdcopy`; cap the output.** `Tty: false`
   so stdout/stderr stay distinguishable, and a hard byte ceiling so a chatty
   program can't exhaust worker memory.
5. **Harden past legacy.** Legacy sets only memory, `--network none` and
   `--rm`. Add `PidsLimit`, `NanoCPUs`, `ReadonlyRootfs` + a small `Tmpfs`,
   `CapDrop: ["ALL"]`, `SecurityOpt: ["no-new-privileges"]`, and a non-root
   `User`. Also **drop `--memory-swap=-1`**: legacy explicitly grants *unlimited
   swap*, which is weaker than Docker's default (§2).
6. **Stay synchronous in the request for v1 — do not put this behind river yet.**
   The pain that makes a blocking 6-second call bad in Rails is thread-pool
   exhaustion; Go's goroutine-per-request model doesn't have it. The real
   ceiling is Docker daemon capacity, which is identical in both shapes, so
   bound it with a semaphore inside the runner. This keeps `checkLesson` exactly
   as `api-spec/lesson-player.tsp` already models it and defers the hardest
   design (§5). Go async only when incremental output becomes a *product*
   feature or when the daemon must be shared across processes with queueing.
7. **Pre-pull images out of the request path** — at worker startup and on
   version promotion, not on first check (§6).

The one thing this recommendation deliberately does *not* do is build the
submit→stream pipeline the current TypeSpec header speculates about. Everything
in §5 says that pipeline is buildable, and that it buys nothing until streaming
output is a feature someone asked for.

---

## 1. Docker access from Go

**The module split is already settled in this repo's favour.** `go mod graph`
shows `github.com/moby/moby/client@v0.5.1` and `github.com/moby/moby/api@v1.55.0`
present as indirect dependencies, and `github.com/docker/docker@...` appearing
**zero** times:

```
$ go mod why github.com/moby/moby/client
hexletbasics/internal/testsupport/testdb
github.com/testcontainers/testcontainers-go
github.com/moby/moby/client
```

Promoting it to a direct dependency therefore costs nothing new.
(`github.com/moby/moby/client` is its own module; the older client lives at
`github.com/docker/docker/client` inside the monolithic module. pkg.go.dev
reports the new package's purpose as "Package client is a Go client for the
Docker Engine API" — https://pkg.go.dev/github.com/moby/moby/client.)

**The API shape is new.** v0.5.1 takes one options struct per call and returns a
result struct — it is *not* the positional-argument API of
`docker/docker/client`, and `NewClientWithOpts` / `WithAPIVersionNegotiation`
are marked deprecated in favour of `client.New`
(https://pkg.go.dev/github.com/moby/moby/client). Verbatim from the module
source (`github.com/moby/moby/client@v0.5.1`):

```go
// container_create.go:18
func (cli *Client) ContainerCreate(ctx context.Context, options ContainerCreateOptions) (ContainerCreateResult, error)

// container_create_opts.go:10-25
type ContainerCreateOptions struct {
	Config           *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
	Platform         *ocispec.Platform
	Name             string
	// Image is a shortcut for Config.Image - only one of Image or Config.Image should be set.
	Image string
}
type ContainerCreateResult struct {
	ID       string
	Warnings []string
}

// container_wait.go:17-41
type ContainerWaitOptions struct { Condition container.WaitCondition }
type ContainerWaitResult struct {
	Result <-chan container.WaitResponse
	Error  <-chan error
}
func (cli *Client) ContainerWait(ctx context.Context, containerID string, options ContainerWaitOptions) ContainerWaitResult

// container_remove.go:9-13
type ContainerRemoveOptions struct { RemoveVolumes, RemoveLinks, Force bool }
```

Note `ContainerWait` returns **no error** — errors arrive on the `Error`
channel, and the call "blocks until the request has been acknowledged by the
server (with a response header)" (`container_wait.go:34-40`).

**Sequence for a one-shot run:**
`ContainerCreate` → `CopyToContainer` (§3) → `ContainerWait` (started *before*
`ContainerStart`, see below) → `ContainerAttach` → `ContainerStart` → read the
attached stream → read the wait channel.

**`AutoRemove` vs explicit remove.** The client documents the ordering hazard
itself: `ContainerWait` "allows the caller to synchronize ContainerWait with
other calls, such as specifying a `next-exit` condition
(`container.WaitConditionNextExit`) before issuing a `Client.ContainerStart`
request" (`container_wait.go:37-40`, sentence at 34-40). The wait conditions are
`not-running` (default), `next-exit`, `removed`
(`github.com/moby/moby/api@v1.55.0`, types/container).

- With `HostConfig.AutoRemove: true`, the container can be gone before you ask
  for anything about it — so **all output must be captured via a pre-start
  attach**, never a post-exit `ContainerLogs`. Issue the wait with `next-exit`
  before `ContainerStart` to close the race.
- With explicit `ContainerRemove` in a `defer`, post-exit `ContainerLogs` is
  available, but a worker crash leaks containers and needs a reaper (legacy
  already has this shape of problem: `ReapStuckVersionBuildsJob`).

**Recommendation:** `AutoRemove: true` + attach-before-start. It has no leak
mode, and streaming output is a prerequisite for §5 anyway. Keep a periodic
prune as a safety net regardless, since a crashed *daemon* leaves containers no
matter which option you pick.

## 2. Resource and time limits

Field names below are from `github.com/moby/moby/api@v1.55.0`
(`types/container.HostConfig` and its embedded `Resources`;
https://pkg.go.dev/github.com/moby/moby/api/types/container).

Legacy flag → API field:

| Legacy flag | Field | Notes |
|---|---|---|
| `--memory=512m` | `Resources.Memory int64` | bytes; docs: minimum allowed value is 6m |
| `--memory-swap=-1` | `Resources.MemorySwap int64` | **-1 = unlimited swap** — see below |
| `--network none` | `HostConfig.NetworkMode NetworkMode` | `"none"` |
| `--rm` | `HostConfig.AutoRemove bool` | §1 |
| `-v host:container` | `HostConfig.Binds []string` / `Mounts []mount.Mount` | replace with `CopyToContainer`, §3 |

On `--memory-swap`, docs.docker.com is explicit: set to `-1`, "the container is
allowed to use unlimited swap, up to the amount available on the host system";
set equal to `--memory`, "the container doesn't have access to swap"
(https://docs.docker.com/engine/containers/resource_constraints/, corroborated
by the CLI reference: "Swap limit equal to memory plus swap: '-1' to enable
unlimited swap", https://docs.docker.com/reference/cli/docker/container/run/).

**So the legacy runner deliberately lets a runaway exercise consume all host
swap.** That is the opposite of a limit. Unless there is a forgotten reason for
it (a language runtime that OOMs at 512m without swap headroom — plausible for
JVM-family courses), `MemorySwap` should equal `Memory`. Flagged as a decision,
not silently changed: it can turn currently-passing exercises red.

Not set by legacy, worth setting:

| Hardening | Field | Documented as |
|---|---|---|
| pids cap | `Resources.PidsLimit *int64` | "Tune container pids limit (set -1 for unlimited)" |
| CPU cap | `Resources.NanoCPUs int64` (or `CPUQuota`/`CPUPeriod`) | `--cpus` ≡ `--cpu-period=100000 --cpu-quota=<n>` |
| read-only rootfs | `HostConfig.ReadonlyRootfs bool` | "Mount the container's root filesystem as read only" |
| writable scratch | `HostConfig.Tmpfs map[string]string` | "Mount a tmpfs directory" |
| drop capabilities | `HostConfig.CapDrop []string` | "Drop Linux capabilities" |
| no-new-privileges | `HostConfig.SecurityOpt []string` | "Security Options" |
| non-root | `Config.User string` | "Username or UID" |
| pid-1 reaping | `HostConfig.Init *bool` | "Run an init inside the container that forwards signals and reaps processes" |
| alt. runtime | `HostConfig.Runtime string` | "Runtime to use for this container" |

(CLI descriptions: https://docs.docker.com/reference/cli/docker/container/run/ —
note `--pids-limit` is documented there but *not* on the resource-constraints
page.)

`PidsLimit` is the one gap that matters most: `--network none` and a memory cap
do nothing against a fork bomb, which is a one-line submission in most of these
languages. `HostConfig.Runtime` is the seam for gVisor/Kata later if a shared
kernel is judged insufficient — but note legacy already accepts a shared kernel,
so tightening that is a raise in the bar, not parity.

## 3. Timeout / infinite-loop detection

The distinction is contractual: `LessonCheckingResponse.result` is
`passed | failed | failed-infinity` (`api-spec/lesson-player.tsp`), and legacy
derives it purely from the exit code — `0 → passed`, `124 → failed-infinity`,
else `failed` (`legacy/app/lib/lesson_tester.rb:43-49`).

**124 is a GNU coreutils guarantee, not a Docker one.** The coreutils manual
documents `timeout`'s exit status as "124 if command times out, and
`--preserve-status` is not specified"; 125 if `timeout` itself fails, 126 if the
command cannot be invoked, 127 if not found, 137 if KILLed (128+9)
(https://www.gnu.org/software/coreutils/manual/html_node/timeout-invocation.html).
Two consequences:

- It requires `timeout` **in the exercise image**. Legacy relies on this today
  (`timeout 6 make --silent -C ...`), so every current image has it — but it is
  an undocumented image contract that a future minimal/distroless image would
  break, silently turning timeouts into `failed`.
- `--kill-after` changes the status to 137, and `--preserve-status` erases the
  signal entirely. Don't add either without remapping.

**Alternative: outer `ContainerWait` + context deadline.** `WaitResponse` gives
`StatusCode int64` and an optional `Error *WaitExitError` (moby api
types/container). A killed container reports the signal-derived status (137 for
SIGKILL), which is *not* distinguishable from an exercise that itself died on
SIGKILL. So the outer clock can tell you "I gave up", but it cannot cleanly tell
you *why* the process stopped.

**Recommendation: keep both, with in-container `timeout` authoritative.**

- In-container `timeout N` → exit 124 → `failed-infinity`. Identical semantics
  to legacy, zero mapping risk, and the timer runs where the work runs.
- Outer deadline at `N + slack` (a `context.WithTimeout` around the wait, then
  kill) → treat as an *infrastructure* error, not `failed-infinity`. It fires
  only when the container never reported at all, which is a daemon problem, not
  a student problem.
- Add a startup assertion that the image has a working `timeout` — a one-off
  probe per image version, so a missing binary is a loud build failure rather
  than silently-wrong grading.

Also carry over the output rule: legacy blanks the output entirely when the
result is `failed-infinity` (`lesson_tester.rb:54`), because a spinning loop's
output is megabytes of noise. That is a good rule and pairs with §4's cap.

## 4. Capturing stdout/stderr

**Attach, not logs.** `ContainerAttach` returns a `HijackedResponse` and must be
called before start to catch the first bytes; `ContainerLogs` reads what the
logging driver retained, which is unavailable once `AutoRemove` fires (§1).
Options are minimal and sufficient (`container_attach.go:10-17`):

```go
type ContainerAttachOptions struct {
	Stream     bool
	Stdin      bool
	Stdout     bool
	Stderr     bool
	DetachKeys string
	Logs       bool
}
```

**Framing depends on `Tty`.** Both `ContainerAttach` and `ContainerLogs` carry
the same doc block (`container_attach.go:29-44`, `container_logs.go:35-50`):

> - If the container is using a TTY, there is only a single stream (stdout) and
>   data is copied directly from the container output stream, no extra
>   multiplexing or headers.
> - If the container is *not* using a TTY, streams for stdout and stderr are
>   multiplexed.
>
> `[8]byte{STREAM_TYPE, 0, 0, 0, SIZE1, SIZE2, SIZE3, SIZE4}[]byte{OUTPUT}`

**Use `Tty: false`** and demux with `stdcopy.StdCopy`, at
`github.com/moby/moby/api/pkg/stdcopy` (present in `api@v1.55.0`, path
`pkg/stdcopy/stdcopy.go`) — *not* the old `github.com/docker/docker/pkg/stdcopy`:

```go
// api/pkg/stdcopy/stdcopy.go:50
func StdCopy(destOut, destErr io.Writer, multiplexedSource io.Reader) (written int64, _ error)
```

`Tty: true` would merge the streams and lose the distinction; the exercise
runner has no interactive terminal, so there is no reason to pay that.
`StdCopy` also surfaces the `Systemerr` (type 3) stream as a returned error
(`stdcopy.go:39-41`), which is a real daemon-side error worth separating from
exercise output.

**Bounding output.** `StdCopy` writes into whatever `io.Writer` you give it, so
the cap belongs in that writer: a small wrapper that stops writing (and records
"truncated") past N bytes. Do not use `io.LimitReader` on the multiplexed
source — cutting mid-frame corrupts the demux. Legacy has no cap at all; it
merely blanks output for `failed-infinity` (§3), so a program that prints
non-stop *just under* the timeout currently returns everything it printed.

Note legacy then base64-encodes the output after `CGI.escapeURIComponent`
(`lesson_tester.rb:54`), commented as scrubbing invalid UTF-8. In Go the
equivalent is scrubbing invalid UTF-8 before it reaches JSON
(`strings.ToValidUTF8`); whether to keep the base64 wire format is a contract
choice — `LessonCheckingResponse.output` is currently just `string`.

## 5. Streaming output back to the browser

### Is the synchronous shape still viable?

**Yes, and it is the recommended v1.** The reason a 6-second blocking request
hurts in Rails is that it pins a Puma thread from a small fixed pool; in Go the
same request costs a goroutine. What it *does* pin either way is a Docker
daemon slot — and moving the run to river does not create daemon capacity, it
only moves the queue. So:

- Synchronous stops being viable when you need **cross-process backpressure**
  (many API replicas, one daemon pool, and you want work to queue durably
  rather than pile up as in-flight requests), or when the run time stops fitting
  a request budget (proxy/load-balancer idle timeouts; at `timeout 6` + pull
  time it comfortably fits).
- Bound it explicitly: a semaphore sized to the daemon's real capacity, and a
  fast 429/503 when full. Unbounded goroutines against a saturated daemon is
  the failure mode Go's cheap concurrency invites.

Concretely: v1 = `checkLesson` handler → runner → Docker, synchronous, as
`api-spec/lesson-player.tsp` already models it. Nothing in this file needs a
contract change to ship that.

### If/when it becomes async

The contract can express it: **ogen 1.23 supports typed SSE responses.** It
recognises `text/event-stream` and auto-enables a data-only event shape, with
`x-ogen-sse-event-shape` (`data-only` | `full` | `full-array`) to override
(`ogen@v1.23.0/openapi/parser/parse_mediatype.go:210-233`,
`openapi/mediatype.go:10-22`), and ships an SSE decoder (`ogen/sse`). So this
is *not* a repeat of the multipart gap that forced
`internal/handlers/attachments.go` outside the generated layer.

Transport from a river worker to an SSE handler in the HTTP process (they are
separate processes, ADR-0010) — three candidates:

| Mechanism | Latency | Durability | Verdict |
|---|---|---|---|
| Postgres `LISTEN`/`NOTIFY` | push, sub-ms | none (fire-and-forget) | **best fit for token/line streaming** |
| watermill-sql topic | ~1s (poll) | durable, at-least-once | wrong tool — see below |
| poll a row | your interval | durable | fine for coarse status, not output |

- **watermill-sql is polling-based.** `SubscriberConfig.PollInterval` defaults
  to `time.Second`
  (`watermill-sql/v4@v4.1.5/pkg/sql/subscriber.go:38-69`), and this repo uses
  the defaults (`internal/events/runtime.go:100-108`). Excellent for durable
  domain events; a 1-second floor per chunk makes it unusable as a token pipe.
  Also: domain events are *facts worth keeping*, and per-line exercise output
  is not one — pushing it through the outbox pollutes an event log that ADR-0004
  scoped deliberately.
- **river's own `Subscribe` is in-process only.** It "panic[s]" on a client
  that will never work jobs (`river@v0.42.0/client.go:1372-1377`), which is
  exactly the insert-only client the HTTP process builds
  (`internal/jobs/jobs.go`, `NewInsertOnlyClient`). It cannot bridge processes.
- **LISTEN/NOTIFY** fits the semantics: output chunks are worthless once the
  viewer is gone. Constraints from the PostgreSQL docs
  (https://www.postgresql.org/docs/current/sql-notify.html): payload "must be
  shorter than 8000 bytes"; "the notify events are not delivered until and
  unless the transaction is committed"; identical channel+payload within one
  transaction is **deduplicated** to a single delivery (so chunks must carry a
  sequence number or they can vanish); a full queue (8GB default) makes
  `NOTIFY` **fail at commit**. Delivery requires a session that has run
  `LISTEN`, so anything sent while no browser is attached is simply lost — hence
  a durable final result row is still needed for the terminal verdict.
  Tooling is in place: pgx exposes `(*pgx.Conn).WaitForNotification(ctx)`
  (`pgx/v5@v5.10.0/conn.go:415`) and the stdlib driver hands back the underlying
  conn via `(*stdlib.Conn).Conn() *pgx.Conn` (`stdlib/sql.go:424`) — the app
  already opens `sql.Open("pgx", dsn)` (`internal/store/store.go:42`). A
  dedicated connection per listener is required (LISTEN is per-session), so
  this needs its own small pool, not the shared one.

**Shape if built:** river job runs the container and NOTIFYs `run:<id>` per
output chunk (sequence-numbered), writes the terminal verdict to a row, and the
SSE handler LISTENs plus reads that row on connect for late joiners and on
close. river's per-job `Timeout(job)` bounds the worker side — "the maximum
amount of time the job is allowed to run before its context is cancelled"
(`river@v0.42.0/worker.go:53-57`).

## 6. Image management

Legacy `ensure_image` shells `docker image inspect`, then `docker pull` + `docker
tag` to a per-version tag, computed in Ruby as `lv<id>` in production for
`hexletbasics`-owned images, `release` otherwise, `latest` outside production
(`legacy/app/models/language/version.rb:99-105`). The Go schema carries
`docker_image` but has **no `image_tag` column** — it was a Rails method, and
there is no column in `legacy/db/schema.rb` either
(`ent/schema/courseversion.go:40-42` has `docker_image`,
`exercise_filename`, `exercise_test_filename`, already commented as "driv[ing]
the runtime exercise runner"). So the tag convention has to be re-decided or
re-derived, not read from the DB.

API equivalents (`github.com/moby/moby/client@v0.5.1`):

```go
// image_pull.go:31
func (cli *Client) ImagePull(ctx context.Context, refStr string, options ImagePullOptions) (ImagePullResponse, error)
// image_pull.go:17-22 — the response is a stream you must consume:
type ImagePullResponse interface {
	io.ReadCloser
	JSONMessages(ctx context.Context) iter.Seq2[jsonstream.Message, error]
	Wait(ctx context.Context) error
}

// image_tag.go:13-22
type ImageTagOptions struct{ Source, Target string }
func (cli *Client) ImageTag(ctx context.Context, options ImageTagOptions) (ImageTagResult, error)

// ImageInspect(ctx, imageID, ...ImageInspectOption) (ImageInspectResult, error)
```

Two gotchas: `ImagePull` **returns before the pull finishes** — you must
`Wait(ctx)` (or drain the reader), or you will start a container against a
half-present image. And `ImageTag` refuses a digest target
(`image_tag.go:35-37`).

**Pre-pull, don't lazy-pull.** A first-check-after-deploy pull of a language
image is seconds to minutes; inside a 6-second grading budget it reads to the
student as a broken lesson. The natural hooks already exist: the course-loader
job promotes a version (`internal/courseloader/loader.go`), and the worker
process has a startup path. Pull on version promotion, and warm current
versions at worker startup; the request path should only ever `ImageInspect` and
fail loudly if the image is missing.

Worth noting the rewrite already dropped Docker from content loading —
`GitFetcher` "[r]eplaces the legacy Docker-image pull — the git tree is
byte-identical to the built image for the content files the parser reads"
(`internal/courseloader/fetch.go:23-24`). So the runner will be the *only*
Docker-daemon consumer in the Go app. That makes the daemon a single, isolable
dependency — and means nothing else in the app currently proves daemon access
works in production.

---

## Decisions that need a human

1. **Sync or async (and therefore: does `checkLesson` change?).** §5 recommends
   keeping it synchronous, which needs no contract change. Overriding that is a
   *product* call — it is only worth the LISTEN/NOTIFY + SSE machinery if
   incremental output during a run is a feature you want. Decide before anyone
   builds, because it determines whether the player is request/response or
   stream-driven.
2. **Hardening posture, and specifically `--memory-swap=-1`.** Legacy grants
   unlimited swap, no pids limit, and (probably) root in-container. Tightening
   all four at once can turn currently-green exercises red. Needs a decision
   plus a canary: run every course's reference solutions under the new limits
   before cutover.
3. **Where the daemon lives.** Worker-local socket, remote daemon, or
   docker-in-docker in k8s (infra is still under `legacy/k8s/`). This decides
   whether bind mounts are even possible (§3 recommends `CopyToContainer` partly
   to make the answer not matter), what the concurrency ceiling is, and how the
   socket is granted — mounting `/var/run/docker.sock` into the worker is
   effectively host root, which deserves its own look.

## Open questions (smaller)

- **Guest checks.** Legacy `check` is `allow_unauthenticated_access` and tracks
  `session[:finished_as_guest]` (`legacy/app/controllers/api/lessons_controller.rb`).
  `checkLesson` in TypeSpec declares no auth, so guests can call it — but
  progress recording depends on the missing `language_members` work
  (`docs/PARITY.md` §1). Sequence them together.
- **Output wire format.** Keep legacy's base64+escape, or send scrubbed UTF-8?
  `LessonCheckingResponse.output` is a plain `string` today.
- **The `timeout` binary as an image contract.** §3 recommends asserting it per
  image version; nothing documents it today.
- **`timeout 6`** is hard-coded in legacy. Per-course budgets (compiled
  languages need more) may be worth a column while the schema is in flux.

## ADR needed

One ADR, "how untrusted exercise code is executed", covering: the client
(`moby/moby/client`), sync-vs-async and why, the isolation baseline as a
concrete `HostConfig`, and the timeout/exit-code contract. It has hard-to-reverse
consequences in two directions — the public contract shape and the security
posture — which is the bar `docs/adr/` has held so far.
