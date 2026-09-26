package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/riverqueue/river"

	"hexletbasics/internal/courseloader"
)

// reapStuckVersionBuildsInterval matches the legacy `every hour` schedule in
// config/recurring.yml. The reaper only acts on builds idle for
// courseloader.StuckBuildAfter, so the interval bounds how long past that
// threshold a dead build can stay visible as `building`.
const reapStuckVersionBuildsInterval = time.Hour

// ReapStuckVersionBuildsArgs is the periodic sweep of course version builds
// whose worker died mid-build. It is the River successor to the legacy
// ReapStuckVersionBuildsJob and is enqueued only by the periodic scheduler.
type ReapStuckVersionBuildsArgs struct{}

// Kind is River's stable job discriminator; do not rename once jobs are enqueued.
func (ReapStuckVersionBuildsArgs) Kind() string { return "reap_stuck_version_builds" }

type reapStuckVersionBuildsWorker struct {
	river.WorkerDefaults[ReapStuckVersionBuildsArgs]
	loader *courseloader.Loader
	logger *slog.Logger
}

// Work reaps the stuck builds and logs each one at warn, as legacy did, so a
// killed build leaves a trace in the worker logs as well as on the version row.
func (w *reapStuckVersionBuildsWorker) Work(ctx context.Context, _ *river.Job[ReapStuckVersionBuildsArgs]) error {
	reaped, err := w.loader.ReapStuck(ctx)
	for _, version := range reaped {
		w.logger.WarnContext(ctx, "reaped stuck version build",
			slog.Int("version_id", version.ID),
			slog.Int("course_id", version.CourseID),
		)
	}
	return err
}

// periodicJobs lists the worker's recurring jobs. River runs the scheduler only
// on the elected leader, so each tick inserts one job across all worker
// replicas. RunOnStart hedges against the scheduler's in-memory state: every
// deploy or leader change restarts the hourly countdown, and without it a
// frequently restarted worker could never reach its first tick. The job takes
// no unique options on purpose — with empty args, River's default unique states
// (which include `completed`) would block every run after the first.
func periodicJobs() []*river.PeriodicJob {
	return []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(reapStuckVersionBuildsInterval),
			func() (river.JobArgs, *river.InsertOpts) {
				return ReapStuckVersionBuildsArgs{}, nil
			},
			&river.PeriodicJobOpts{
				ID:         ReapStuckVersionBuildsArgs{}.Kind(),
				RunOnStart: true,
			},
		),
	}
}
