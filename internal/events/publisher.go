package events

import (
	"context"
	"database/sql"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	"hexletbasics/ent"
	"hexletbasics/internal/store"
)

// TxPublisher is the transactional seam used by business modules.
type TxPublisher interface {
	Publish(ctx context.Context, tx *sql.Tx, event Event) error
}

// StandalonePublisher records facts with no accompanying state mutation.
type StandalonePublisher interface {
	PublishStandalone(ctx context.Context, event Event) error
}

// Consumer is one reaction to domain facts: the job it wants queued for event,
// or false when it does not react to that event. The job's own InsertOpts carry
// its queue, attempts and uniqueness, as for any other insert of its kind.
type Consumer func(event Event) (river.JobArgs, bool)

// jobInserter is the part of the River client the publisher uses: an insert
// that joins the caller's transaction.
type jobInserter interface {
	InsertTx(ctx context.Context, tx *sql.Tx, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

// Publisher fans a fact out to its consumers as River jobs. The consumers are
// named once, at composition (di), so a producer raises a fact without knowing
// who reacts to it and a new consumer touches no producer.
type Publisher struct {
	store     store.Transactor
	jobs      jobInserter
	consumers []Consumer
}

// NewPublisher builds the publisher over the River client and the consumers.
func NewPublisher(txStore store.Transactor, jobs jobInserter, consumers ...Consumer) *Publisher {
	return &Publisher{store: txStore, jobs: jobs, consumers: consumers}
}

// PublishStandalone owns a short transaction for a fact with no state change of
// its own, such as a read of personal data.
func (p *Publisher) PublishStandalone(ctx context.Context, event Event) error {
	return p.store.WithinTx(ctx, func(tx *sql.Tx, _ *ent.Client) error {
		return p.Publish(ctx, tx, event)
	})
}

// Publish inserts, in tx, the job of every consumer that reacts to event, so the
// jobs exist exactly when the caller's state change commits. An error rolls the
// caller's transaction back.
func (p *Publisher) Publish(ctx context.Context, tx *sql.Tx, event Event) error {
	for _, consume := range p.consumers {
		args, ok := consume(event)
		if !ok {
			continue
		}

		_, err := p.jobs.InsertTx(ctx, tx, args, nil)
		if err != nil {
			return errs.With("event", event.eventName(), "job_kind", args.Kind()).Wrapf(err, "queue event consumer")
		}
	}

	return nil
}
