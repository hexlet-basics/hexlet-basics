// Package books owns the free book funnel (legacy `BooksController`): a
// signed-in user requests the book once, and the request records whether they
// have downloaded it since.
package books

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"hexletbasics/ent"
	"hexletbasics/ent/bookrequest"
	"hexletbasics/internal/events"
	"hexletbasics/internal/store"
)

// Request states, legacy BookRequest's typed_enum values.
const (
	StateRequested  = "requested"
	StateDownloaded = "downloaded"
)

// Requester is the seam the handler depends on.
type Requester interface {
	Request(ctx context.Context, userID int, locale string) error
}

// Recorder atomically stores a user's first book request and its
// BookRequested outbox record, so the fact is raised exactly once per user.
type Recorder struct {
	store     store.Transactor
	publisher events.TxPublisher
	now       func() time.Time
}

// NewRecorder builds the production book request recorder.
func NewRecorder(txStore store.Transactor, publisher events.TxPublisher) *Recorder {
	return &Recorder{store: txStore, publisher: publisher, now: time.Now}
}

// Request records the user's request, idempotently. Legacy did
// `find_or_initialize_by` and published only for a new record; here the unique
// index on user_id decides instead, so two concurrent first requests still
// store one row and raise one event. `ON CONFLICT DO NOTHING` returns no row
// for a repeat, which is how a repeat is told apart from a first request, and
// it does not abort the transaction.
func (r *Recorder) Request(ctx context.Context, userID int, locale string) error {
	return r.store.WithinTx(ctx, func(tx *sql.Tx, db *ent.Client) error {
		_, err := db.BookRequest.Create().
			SetUserID(userID).
			SetState(StateRequested).
			OnConflictColumns(bookrequest.FieldUserID).
			DoNothing().
			ID(ctx)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("create book request for user %d: %w", userID, err)
		}
		if err := r.publisher.Publish(ctx, tx, events.BookRequested{
			Locale:     locale,
			OccurredAt: r.now().UTC(),
		}); err != nil {
			return fmt.Errorf("publish book requested: %w", err)
		}
		return nil
	})
}
