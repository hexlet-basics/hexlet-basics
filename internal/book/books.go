// Package book owns the free book funnel (legacy `BooksController`): a
// signed-in user requests the book once, and the request records whether they
// have downloaded it since.
package book

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/bookrequest"
	"hexletbasics/internal/events"
	"hexletbasics/internal/store"
)

// errs is the package template for errors (ADR-0018): static messages, the data
// in With, wrapped once at the exported function that knows the operation.
var errs = oops.In("books")

// Request states, legacy BookRequest's typed_enum values.
const (
	StateRequested  = "requested"
	StateDownloaded = "downloaded"
)

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
	err := r.store.WithinTx(ctx, func(tx *sql.Tx, db *ent.Client) error {
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
			return errs.Wrapf(err, "save book request")
		}

		if err := r.publisher.Publish(ctx, tx, events.BookRequested{
			Locale:     locale,
			OccurredAt: r.now().UTC(),
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errs.With("user_id", userID).Wrapf(err, "record book request")
	}

	return nil
}

// Requested reports whether the user has requested the book, in any state: the
// page picks the request or the download button by it.
func (r *Recorder) Requested(ctx context.Context, userID int) (bool, error) {
	var requested bool

	err := r.store.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		exists, err := db.BookRequest.Query().Where(bookrequest.UserID(userID)).Exist(ctx)
		requested = exists

		return errs.Wrapf(err, "check book request")
	})
	if err != nil {
		return false, errs.With("user_id", userID).Wrapf(err, "check book request")
	}

	return requested, nil
}

// MarkDownloaded moves the user's request to downloaded, on every download as
// legacy did, and reports whether there was a request to move.
func (r *Recorder) MarkDownloaded(ctx context.Context, userID int) (bool, error) {
	var marked bool

	err := r.store.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		count, err := db.BookRequest.Update().
			Where(bookrequest.UserID(userID)).
			SetState(StateDownloaded).
			Save(ctx)
		marked = count > 0

		return errs.Wrapf(err, "save book request")
	})
	if err != nil {
		return false, errs.With("user_id", userID).Wrapf(err, "mark book downloaded")
	}

	return marked, nil
}
