package handlers_test

import (
	"context"
	"errors"

	entgo "entgo.io/ent"

	"hexletbasics/ent"
	"hexletbasics/internal/testsupport"
)

// What an Audit Record does when the database fails under the operation it
// belongs to. The harness has one transaction and no way to take the database
// down mid-request, so the failure is injected where ent runs the statement: a
// hook refuses the writes of an entity, an interceptor refuses its reads. The
// handlers and the audit hook share h.DB's hook stack, so both see it. These are
// the only storage-failure tests left in the package: the guarantee they check
// (no change without its record) is the audit trail's own.

var errStorageDown = errors.New("storage is down")

// breakReads makes every read of the entity fail.
func breakReads(h *testsupport.Harness, entity string) {
	h.DB.Intercept(ent.InterceptFunc(func(next ent.Querier) ent.Querier {
		return ent.QuerierFunc(func(ctx context.Context, q ent.Query) (ent.Value, error) {
			if entgo.QueryFromContext(ctx).Type == entity {
				return nil, errStorageDown
			}
			return next.Query(ctx, q)
		})
	}))
}

// breakWrites makes every mutation of the entity fail, while its reads keep
// working.
func breakWrites(h *testsupport.Harness, entity string) {
	h.DB.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if m.Type() == entity {
				return nil, errStorageDown
			}
			return next.Mutate(ctx, m)
		})
	})
}
