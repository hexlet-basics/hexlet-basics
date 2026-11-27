package jobs

import (
	"context"

	"github.com/riverqueue/river"

	"hexletbasics/ent"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/events"
)

// RecordPersonalDataReadArgs carries the fact that a staff member read personal
// data (ADR-0020). The request only publishes it; this job writes the Audit
// Record in the worker, so recording never slows the read.
type RecordPersonalDataReadArgs struct {
	Event events.StaffReadPersonalData `json:"event"`
}

// Kind is River's stable job discriminator; do not rename once jobs are enqueued.
func (RecordPersonalDataReadArgs) Kind() string { return "record_personal_data_read" }

// recordPersonalDataReadMaxAttempts is about an hour of River's attempt^4
// backoff. The job is one insert, so it fails only while the database does;
// that covers a restart or a failover, and a record still missing after it is
// an incident the final attempt reports.
const recordPersonalDataReadMaxAttempts = 8

// InsertOpts bounds retries. There is no uniqueness: each publish is one read.
func (RecordPersonalDataReadArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{MaxAttempts: recordPersonalDataReadMaxAttempts}
}

// recordPersonalDataRead is the events.Consumer that queues this job.
func recordPersonalDataRead(event events.Event) (river.JobArgs, bool) {
	read, ok := event.(events.StaffReadPersonalData)
	if !ok {
		return nil, false
	}

	return RecordPersonalDataReadArgs{Event: read}, true
}

// RecordPersonalDataRead writes the Audit Records of staff reads of personal
// data.
func RecordPersonalDataRead(db *ent.Client) Module {
	return module(RecordPersonalDataReadWorker(db))
}

// RecordPersonalDataReadWorker is the worker alone, so that a test can work a
// job through it in a transaction (rivertest.Worker).
func RecordPersonalDataReadWorker(db *ent.Client) river.Worker[RecordPersonalDataReadArgs] {
	return &recordPersonalDataReadWorker{db: db, WorkerDefaults: river.WorkerDefaults[RecordPersonalDataReadArgs]{}}
}

// recordPersonalDataReadWorker writes the record through the audit write path
// with the time the read happened, not the time of this run, and holds
// identifiers and a count only.
//
// Not idempotent, by name (ADR-0027): a job that wrote its record and was then
// lost before River marked it done runs again and leaves a second record of the
// same read. An extra line in a trail is the safe direction to err, and
// ADR-0020 already accepts that a record may be lost on a crash.
type recordPersonalDataReadWorker struct {
	river.WorkerDefaults[RecordPersonalDataReadArgs]

	db *ent.Client
}

func (w *recordPersonalDataReadWorker) Work(ctx context.Context, job *river.Job[RecordPersonalDataReadArgs]) error {
	event := job.Args.Event
	entry := audit.Entry{
		Actor: event.ActorUserID,
		Source: audit.Source{
			Channel:        audit.Channel(event.Channel),
			ConnectedAppID: event.ConnectedAppID,
			TokenID:        event.TokenID,
		},
		RequestID:  event.RequestID,
		Address:    event.Address,
		UserAgent:  event.UserAgent,
		Operation:  event.Operation,
		Locale:     event.Locale,
		Outcome:    audit.OutcomeSucceeded,
		OccurredAt: event.OccurredAt,
		EntityType: event.EntityType,
		Reads:      &audit.Reads{Count: event.Count, IDs: event.EntityIDs},
		EntityID:   "",
	}
	// A read of exactly one row is filed under it, like a change.
	if event.Count == 1 && len(event.EntityIDs) == 1 {
		entry.EntityID = event.EntityIDs[0]
	}

	err := audit.Write(ctx, w.db, entry, nil)
	if err != nil {
		return cancelUnfixable(errs.Wrapf(err, "record read of personal data"))
	}

	return nil
}
