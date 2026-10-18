package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"

	"hexletbasics/internal/courseloader"
)

// ExerciseLoaderArgs enqueues a build of one course version: fetch its exercises
// repo, parse it, and (on success) promote it live. This is the river successor
// to the legacy ExerciseLoaderJob, triggered by the admin createVersion action
// and the GitHub webhook.
type ExerciseLoaderArgs struct {
	// VersionID is the pre-created CourseVersion row (state `created`) to build.
	VersionID int `json:"version_id"`
}

// Kind is river's stable job discriminator; do not rename once jobs are enqueued.
func (ExerciseLoaderArgs) Kind() string { return "exercise_loader" }

// exerciseLoaderMaxAttempts is the one named exception to idempotent retries:
// a build is not safe to re-run by the queue. It fetches, parses, uploads and
// then promotes a version, and the version row records the failure and is
// re-triggered by an admin or the next webhook, so a blind retry would only
// repeat a build the row already marks failed.
const exerciseLoaderMaxAttempts = 1

// InsertOpts caps the job at a single attempt and deduplicates jobs for the same
// version. The loader's SQL claim remains the correctness boundary, while River
// uniqueness avoids spending worker capacity on normal duplicate triggers.
// Build failures are recorded on the version row and re-triggered manually (or
// by the next webhook), matching legacy's no-retry behavior.
func (ExerciseLoaderArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		MaxAttempts: exerciseLoaderMaxAttempts,
		Queue:       QueueBuilds,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	}
}

// ExerciseLoader runs course-version builds with loader.
func ExerciseLoader(loader *courseloader.Loader) Module {
	return module(ExerciseLoaderWorker(loader))
}

// ExerciseLoaderWorker is the worker alone, without its registration, so that a test
// can work a job through it in a transaction (rivertest.Worker).
func ExerciseLoaderWorker(loader *courseloader.Loader) river.Worker[ExerciseLoaderArgs] {
	return &exerciseLoaderWorker{loader: loader}
}

type exerciseLoaderWorker struct {
	river.WorkerDefaults[ExerciseLoaderArgs]
	loader *courseloader.Loader
}

// exerciseLoaderTimeout bounds one build. River's default job timeout is one
// minute, far below a git fetch plus parse plus image uploads, and the build's
// single attempt would be cancelled mid-way.
const exerciseLoaderTimeout = 15 * time.Minute

// Timeout overrides River's one-minute default for course builds.
func (*exerciseLoaderWorker) Timeout(*river.Job[ExerciseLoaderArgs]) time.Duration {
	return exerciseLoaderTimeout
}

// Work runs the loader for the job's version. The loader records failure state on
// the version; returning the error here surfaces it in river's job history (the
// job is discarded after its single attempt).
func (w *exerciseLoaderWorker) Work(ctx context.Context, job *river.Job[ExerciseLoaderArgs]) error {
	return cancelUnfixable(w.loader.Run(ctx, job.Args.VersionID))
}
