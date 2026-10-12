package jobs

import (
	"context"
	"errors"

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

// amoCRMLeadMaxAttempts is about an hour of River's attempt^4 backoff, enough
// to ride out an amoCRM outage; River's default of 25 retries for weeks.
const amoCRMLeadMaxAttempts = 8

// InsertOpts deduplicates a redelivered event by its args and bounds retries.
// Uniqueness lives here, not at the Watermill call site, so every insert of the
// kind gets it.
func (AmoCRMLeadArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		MaxAttempts: amoCRMLeadMaxAttempts,
		UniqueOpts:  river.UniqueOpts{ByArgs: true},
	}
}

// AmoCRMLead syncs created leads to amoCRM through syncer.
func AmoCRMLead(syncer LeadSyncer) Module {
	return module(&amoCRMLeadWorker{syncer: syncer})
}

type amoCRMLeadWorker struct {
	river.WorkerDefaults[AmoCRMLeadArgs]
	syncer LeadSyncer
}

func (w *amoCRMLeadWorker) Work(ctx context.Context, job *river.Job[AmoCRMLeadArgs]) error {
	err := w.syncer.CreateLead(ctx, job.Args.Event)
	if err == nil {
		return nil
	}
	var classified retryableError
	if errors.As(err, &classified) && !classified.Retryable() {
		return river.JobCancel(err)
	}
	return err
}

type retryableError interface {
	error
	Retryable() bool
}
