package handlers

import (
	"context"

	"hexletbasics/internal/api"
	"hexletbasics/internal/paging"
)

// pageQuery is the slice of an ent query builder the paginated list needs. It
// self-references Q so the fluent Offset/Limit/Clone calls keep returning the
// concrete builder type — every generated ent `*XQuery` satisfies it, which is
// what lets listPage stay generic over otherwise-unrelated ent query types.
type pageQuery[E any, Q any] interface {
	Clone() Q
	Count(context.Context) (int, error)
	Offset(int) Q
	Limit(int) Q
	All(context.Context) ([]E, error)
}

// listPage runs the standard admin list once, generically: a total count plus a
// single ordered, paginated page, converted to the API model and wrapped in the
// resource's `*Page`. The caller supplies only the resource-specific bits — the
// scoped, ordered, eager-loading query, the converter, and the `*Page`
// constructor. The count runs on a Clone, so the page's offset/limit never
// leak into it, and the one query carries the same filter for both (a bare
// Count would over-report rows the scoping hides).
func listPage[E any, A any, P any, Q pageQuery[E, Q]](
	ctx context.Context,
	pageParam, perPageParam api.OptInt32,
	query Q,
	conv func([]E) []A,
	mkPage func(items []A, total, page, perPage int32) P,
) (P, error) {
	return listPageWith(ctx, pageParam, perPageParam, query,
		func(_ context.Context, rows []E) ([]A, error) { return conv(rows), nil },
		mkPage,
	)
}

// listPageWith is listPage for a read model that needs more than the rows it
// is handed — a converter that queries (batched covers, like counts) and so can
// fail. listPage is the pure-converter case layered on top of it.
func listPageWith[E any, A any, P any, Q pageQuery[E, Q]](
	ctx context.Context,
	pageParam, perPageParam api.OptInt32,
	query Q,
	conv func(context.Context, []E) ([]A, error),
	mkPage func(items []A, total, page, perPage int32) P,
) (P, error) {
	var zero P
	page := paging.New(pageParam, perPageParam)

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return zero, err
	}

	rows, err := query.Offset(page.Offset()).Limit(page.Limit()).All(ctx)
	if err != nil {
		return zero, err
	}

	items, err := conv(ctx, rows)
	if err != nil {
		return zero, err
	}

	return mkPage(items, int32(total), page.Page, page.PerPage), nil
}

// listAll runs a non-paginated list (a bare array response, e.g. the nested QnA
// items) : fetch all rows from the caller's scoped+ordered query, then convert.
func listAll[E any, A any](
	ctx context.Context,
	all func(context.Context) ([]E, error),
	conv func([]E) []A,
) ([]A, error) {
	rows, err := all(ctx)
	if err != nil {
		return nil, err
	}
	return conv(rows), nil
}

// getOne fetches a single row by id and converts it, returning a pointer to the
// API value. The fetch error (including ent's not-found) is returned as-is for
// the central ErrorHandler to map to 404. This is the shared body every
// AdminGet* handler used to repeat inline.
func getOne[E any, A any](
	ctx context.Context,
	id int,
	fetch func(context.Context, int) (E, error),
	conv func(E) A,
) (*A, error) {
	row, err := fetch(ctx, id)
	if err != nil {
		return nil, err
	}
	item := conv(row)
	return &item, nil
}
