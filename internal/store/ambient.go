package store

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"

	"hexletbasics/ent"
)

// The ambient transaction is how one request-wide unit of work reaches code
// that never sees a transaction. The generated ogen server hands handlers a
// context and nothing else (ADR-0011), so the unit of work (the audit middleware,
// ADR-0020 addendum) opens the transaction, parks it in the context with
// WithAmbientTx, and from then on:
//
//   - the application ent client built by NewClient routes every query whose
//     context carries it through that transaction instead of the pool, so a
//     handler written against a plain `s.db` takes part without knowing, and
//     reads after a write see the write;
//   - Store.WithinTx called with such a context joins it instead of opening a
//     second, independent transaction, so modules that own an atomic step
//     (blog.Remover, relatedcourses.Replacer) are inside the same commit.
//
// Without the key in the context nothing changes: autocommit and own
// transactions as before.
type ambientKey struct{}

type ambient struct {
	tx     *sql.Tx
	client *ent.Client
	driver dialect.Driver
}

// WithAmbientTx returns a context that carries tx (and the ent client bound to
// it) as the transaction every store-built client and Store.WithinTx join.
// The caller keeps owning commit and rollback.
func WithAmbientTx(ctx context.Context, tx *sql.Tx, client *ent.Client) context.Context {
	return context.WithValue(ctx, ambientKey{}, &ambient{tx: tx, client: client, driver: newTxDriver(tx)})
}

func ambientFrom(ctx context.Context) (*ambient, bool) {
	a, ok := ctx.Value(ambientKey{}).(*ambient)

	return a, ok
}

// routingDriver sends a statement to the ambient transaction when the context
// carries one and to the pool otherwise. ent passes the caller's context to
// every Exec, Query and Tx, which is what makes the choice per request instead
// of per client.
type routingDriver struct{ dialect.Driver }

func (d routingDriver) Exec(ctx context.Context, query string, args, dest any) error {
	if a, ok := ambientFrom(ctx); ok {
		err := a.driver.Exec(ctx, query, args, dest)
		if err != nil {
			return errs.Wrapf(err, "exec")
		}

		return nil
	}

	err := d.Driver.Exec(ctx, query, args, dest)
	if err != nil {
		return errs.Wrapf(err, "exec")
	}

	return nil
}

func (d routingDriver) Query(ctx context.Context, query string, args, dest any) error {
	if a, ok := ambientFrom(ctx); ok {
		err := a.driver.Query(ctx, query, args, dest)
		if err != nil {
			return errs.Wrapf(err, "query")
		}

		return nil
	}

	err := d.Driver.Query(ctx, query, args, dest)
	if err != nil {
		return errs.Wrapf(err, "query")
	}

	return nil
}

// Tx joins the ambient transaction (a no-op commit and rollback, like
// NewTxClient) so ent's own client.Tx does not open a nested one.
func (d routingDriver) Tx(ctx context.Context) (dialect.Tx, error) {
	if a, ok := ambientFrom(ctx); ok {
		result, err := a.driver.Tx(ctx)
		if err != nil {
			return result, errs.Wrapf(err, "begin transaction")
		}

		return result, nil
	}

	result, err := d.Driver.Tx(ctx)
	if err != nil {
		return result, errs.Wrapf(err, "begin transaction")
	}

	return result, nil
}
