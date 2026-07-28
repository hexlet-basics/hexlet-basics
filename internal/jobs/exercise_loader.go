package jobs

import (
	"context"

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

// InsertOpts caps the job at a single attempt: a build failure is recorded on the
// version row (state `failed`) by the loader itself, and re-running a broken build
// 25 times (river's default) would just re-fail. Transient issues are re-triggered
// manually (or by the next webhook), matching legacy's no-retry behavior.
func (ExerciseLoaderArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{MaxAttempts: 1}
}

type exerciseLoaderWorker struct {
	river.WorkerDefaults[ExerciseLoaderArgs]
	loader *courseloader.Loader
}

// Work runs the loader for the job's version. The loader records failure state on
// the version; returning the error here surfaces it in river's job history (the
// job is discarded after its single attempt).
func (w *exerciseLoaderWorker) Work(ctx context.Context, job *river.Job[ExerciseLoaderArgs]) error {
	return w.loader.Run(ctx, job.Args.VersionID)
}
