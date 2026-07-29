package jobs

import (
	"context"

	"github.com/riverqueue/river"

	"hexletbasics/internal/events"
)

// LeadSyncer is the amoCRM integration seam used by the River worker.
type LeadSyncer interface {
	CreateLead(ctx context.Context, event events.LeadCreated) error
}

// AmoCRMLeadArgs carries the immutable lead snapshot from the domain event.
type AmoCRMLeadArgs struct {
	Event events.LeadCreated `json:"event"`
}

// Kind is River's stable job discriminator.
func (AmoCRMLeadArgs) Kind() string { return "amocrm_lead" }

type amoCRMLeadWorker struct {
	river.WorkerDefaults[AmoCRMLeadArgs]
	syncer LeadSyncer
}

func (w *amoCRMLeadWorker) Work(ctx context.Context, job *river.Job[AmoCRMLeadArgs]) error {
	return w.syncer.CreateLead(ctx, job.Args.Event)
}
