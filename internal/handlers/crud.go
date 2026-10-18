package handlers

import (
	"context"

	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/paging"
)

// listFrom is the page a module's list answers, wrapped in the resource's
// `*Page`: the window comes from the request, the rows and the total from the
// module, and the converter and the `*Page` constructor are the resource's.
func listFrom[E any, A any, P any](
	ctx context.Context,
	pageParam, perPageParam api.OptInt32,
	fetch func(context.Context, domain.Page) ([]E, int, error),
	conv func([]E) []A,
	mkPage func(items []A, total, page, perPage int32) P,
) (P, error) {
	return listFromWith(ctx, pageParam, perPageParam, fetch,
		func(_ context.Context, rows []E) ([]A, error) { return conv(rows), nil },
		mkPage,
	)
}

// listFromWith is listFrom for a read model that needs more than the rows it is
// handed: a converter that queries (batched covers, like counts) and so can
// fail. listFrom is the pure-converter case layered on top of it.
func listFromWith[E any, A any, P any](
	ctx context.Context,
	pageParam, perPageParam api.OptInt32,
	fetch func(context.Context, domain.Page) ([]E, int, error),
	conv func(context.Context, []E) ([]A, error),
	mkPage func(items []A, total, page, perPage int32) P,
) (P, error) {
	var zero P

	window := paging.New(pageParam, perPageParam)

	rows, total, err := fetch(ctx, domain.Page{Offset: window.Offset(), Limit: window.Limit()})
	if err != nil {
		return zero, err
	}

	items, err := conv(ctx, rows)
	if err != nil {
		return zero, err
	}

	return mkPage(items, apiconv.Int32FromInt(total), window.Page, window.PerPage), nil
}
