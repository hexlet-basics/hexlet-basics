package sessions

import (
	"context"
	"time"

	"hexletbasics/ent"
	session "hexletbasics/ent/session"
)

// store is scs's session store over the application's ent client.
//
// It is not scs's own postgresstore for two reasons. That store takes a
// *sql.DB and so cannot join the transaction a test rolls back, which would leave
// session rows behind every test. And it starts a cleanup goroutine of its own,
// which ADR-0010 forbids: expired rows are removed by a River job instead
// (deleteExpired).
//
// scs hands the store the token already hashed (HashTokenInStore), so the key
// here is a digest, not something a cookie can be built from.
type store struct {
	db    *ent.Client
	clock func() time.Time
}

// FindCtx answers only for a live session. An expired row is not found rather
// than an error: it is a visitor whose session ended, not a fault, and the row
// itself is left for deleteExpired.
func (s *store) FindCtx(ctx context.Context, token string) ([]byte, bool, error) {
	row, err := s.db.Session.Query().
		Where(session.ID(token), session.ExpiryGT(s.clock().UTC())).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return row.Data, true, nil
}

// CommitCtx writes the session, replacing the data and expiry of one that is
// already there: scs commits again on every request when an idle timeout moves
// the expiry forward.
func (s *store) CommitCtx(ctx context.Context, token string, data []byte, expiry time.Time) error {
	return s.db.Session.Create().
		SetID(token).
		SetData(data).
		SetExpiry(expiry.UTC()).
		OnConflictColumns(session.FieldID).
		UpdateNewValues().
		Exec(ctx)
}

// DeleteCtx removes a session. A token that is not there is not an error, as
// scs requires: signing out twice is not a failure.
func (s *store) DeleteCtx(ctx context.Context, token string) error {
	_, err := s.db.Session.Delete().Where(session.ID(token)).Exec(ctx)
	return err
}

// deleteExpired removes every session whose expiry has passed and reports how
// many it removed.
func (s *store) deleteExpired(ctx context.Context) (int, error) {
	return s.db.Session.Delete().Where(session.ExpiryLTE(s.clock().UTC())).Exec(ctx)
}

// The context-free methods scs's Store interface demands. scs prefers the Ctx
// variants when a store has them, so these are only the interface's price.

func (s *store) Find(token string) ([]byte, bool, error) {
	return s.FindCtx(context.Background(), token)
}

func (s *store) Commit(token string, data []byte, expiry time.Time) error {
	return s.CommitCtx(context.Background(), token, data, expiry)
}

func (s *store) Delete(token string) error {
	return s.DeleteCtx(context.Background(), token)
}
