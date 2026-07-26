package handlers

import "hexletbasics/internal/api"

// defaultPerPage matches the legacy pagy default page size.
const defaultPerPage = 20

// pagination is a resolved 1-based page window. It centralizes the default and
// offset math that every list endpoint would otherwise re-inline.
type pagination struct {
	Page    int32
	PerPage int32
}

// newPagination resolves the optional page/perPage query params, applying
// defaults and ignoring non-positive values (page defaults to 1, perPage to
// defaultPerPage).
func newPagination(page, perPage api.OptInt32) pagination {
	return pagination{
		Page:    optPositive(page, 1),
		PerPage: optPositive(perPage, defaultPerPage),
	}
}

// Offset is the SQL OFFSET for this window.
func (p pagination) Offset() int { return int((p.Page - 1) * p.PerPage) }

// Limit is the SQL LIMIT for this window.
func (p pagination) Limit() int { return int(p.PerPage) }

// optPositive returns v when present and > 0, otherwise def.
func optPositive(v api.OptInt32, def int32) int32 {
	if n, ok := v.Get(); ok && n > 0 {
		return n
	}
	return def
}
