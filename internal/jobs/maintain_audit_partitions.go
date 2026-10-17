package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/riverqueue/river"

	"hexletbasics/internal/audit"
)

// maintainAuditPartitionsInterval is daily: partitions are created
// audit.PartitionsAhead months ahead, so a day's delay is irrelevant, and
// retention only changes at month granularity.
const maintainAuditPartitionsInterval = 24 * time.Hour

// MaintainAuditPartitionsArgs is the periodic upkeep of the audit table
// (ADR-0018, ticket #846): create the partitions of the coming months, then
// drop the whole partitions older than the retention. It is enqueued only by
// the periodic scheduler.
type MaintainAuditPartitionsArgs struct{}

// Kind is River's stable job discriminator; do not rename once jobs are enqueued.
func (MaintainAuditPartitionsArgs) Kind() string { return "maintain_audit_partitions" }

type maintainAuditPartitionsWorker struct {
	river.WorkerDefaults[MaintainAuditPartitionsArgs]
	partitions      *audit.Partitions
	retentionMonths int
	logger          *slog.Logger
	now             func() time.Time
}

// Work creates before it drops, so a failure of the drop never leaves the
// table without next month's partition. Dropping is the only way a record
// leaves the trail; each dropped partition is logged at info, because losing
// audit data should never be silent.
func (w *maintainAuditPartitionsWorker) Work(ctx context.Context, _ *river.Job[MaintainAuditPartitionsArgs]) error {
	now := w.now()
	if err := w.partitions.EnsureAhead(ctx, now); err != nil {
		return fmt.Errorf("maintain audit partitions: %w", err)
	}
	dropped, err := w.partitions.DropExpired(ctx, now, w.retentionMonths)
	if err != nil {
		return fmt.Errorf("maintain audit partitions: %w", err)
	}
	for _, name := range dropped {
		w.logger.InfoContext(ctx, "dropped expired audit partition",
			slog.String("partition", name),
			slog.Int("retention_months", w.retentionMonths),
		)
	}
	return nil
}

// MaintainAuditPartitions keeps the audit table's monthly partitions: next
// months ahead of time, expired ones dropped.
func MaintainAuditPartitions(partitions *audit.Partitions, retentionMonths int, logger *slog.Logger) Module {
	return module(MaintainAuditPartitionsWorker(partitions, retentionMonths, logger), maintainAuditPartitionsSchedule())
}

// MaintainAuditPartitionsWorker is the worker alone, without its schedule, so
// that a test can work the upkeep through it in a transaction (rivertest.Worker).
func MaintainAuditPartitionsWorker(partitions *audit.Partitions, retentionMonths int, logger *slog.Logger) river.Worker[MaintainAuditPartitionsArgs] {
	return &maintainAuditPartitionsWorker{
		partitions:      partitions,
		retentionMonths: retentionMonths,
		logger:          logger,
		now:             time.Now,
	}
}

// maintainAuditPartitionsSchedule is the recurring insert. RunOnStart makes a
// fresh deploy (or a worker that was down over a month boundary) create the
// partitions at once instead of after the first interval; no unique options,
// for the same reason as the stuck-build reaper (empty args and River's
// default unique states would block every run after the first).
func maintainAuditPartitionsSchedule() *river.PeriodicJob {
	return river.NewPeriodicJob(
		river.PeriodicInterval(maintainAuditPartitionsInterval),
		func() (river.JobArgs, *river.InsertOpts) {
			return MaintainAuditPartitionsArgs{}, nil
		},
		&river.PeriodicJobOpts{
			ID:         MaintainAuditPartitionsArgs{}.Kind(),
			RunOnStart: true,
		},
	)
}
