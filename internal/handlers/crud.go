package handlers

import (
	"context"

	"hexletbasics/internal/api"
)

// pageQuery is the slice of an ent query builder the paginated list needs. It
// self-references Q so the fluent Offset/Limit calls keep returning the concrete
// builder type — every generated ent `*XQuery` satisfies it, which is what lets
// listPage stay generic over otherwise-unrelated ent query types.
type pageQuery[E any, Q any] interface {
	Count(context.Context) (int, error)
	Offset(int) Q
	Limit(int) Q
	All(context.Context) ([]E, error)
}

// listPage runs the standard admin list once, generically: a total count plus a
// single ordered, paginated page, converted to the API model and wrapped in the
// resource's `*Page`. The caller supplies only the resource-specific bits — a
// fresh ordered query per call (count and page read from independent builders,
// so eager-load/limit state never leaks between them), the converter, and the
// `*Page` constructor. This is the shared body every AdminList* handler used to
// repeat inline.
func listPage[E any, A any, P any, Q pageQuery[E, Q]](
	ctx context.Context,
	pageParam, perPageParam api.OptInt32,
	newQuery func() Q,
	conv func([]E) []A,
	mkPage func(items []A, total, page, perPage int32) P,
) (P, error) {
	page := newPagination(pageParam, perPageParam)

	total, err := newQuery().Count(ctx)
	if err != nil {
		var zero P
		return zero, err
	}

	rows, err := newQuery().Offset(page.Offset()).Limit(page.Limit()).All(ctx)
	if err != nil {
		var zero P
		return zero, err
	}

	return mkPage(conv(rows), int32(total), page.Page, page.PerPage), nil
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
