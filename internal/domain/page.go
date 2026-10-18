package domain

import "context"

// Page is a window of rows, as an offset and a size.
type Page struct {
	Offset int
	Limit  int
}

// pagedQuery is the slice of an ent query builder a paged list needs. It
// self-references Q so the fluent Offset and Limit calls keep returning the
// concrete builder type: every generated ent `*XQuery` satisfies it, which is
// what lets Paged stay generic over otherwise unrelated query types.
type pagedQuery[E any, Q any] interface {
	Clone() Q
	Count(ctx context.Context) (int, error)
	Offset(offset int) Q
	Limit(limit int) Q
	All(ctx context.Context) ([]E, error)
}

// Paged runs the standard list once: how many rows the query matches and the
// one window of them. The count runs on a clone, so the window never leaks into
// it, and one query carries the scope and the order for both: a bare Count
// would report rows the scope hides.
func Paged[E any, Q pagedQuery[E, Q]](ctx context.Context, query Q, page Page) ([]E, int, error) {
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, FromStore(err)
	}

	rows, err := query.Offset(page.Offset).Limit(page.Limit).All(ctx)
	if err != nil {
		return nil, 0, FromStore(err)
	}

	return rows, total, nil
}
