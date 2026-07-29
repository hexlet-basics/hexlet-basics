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

// newPagination resolves absent page/perPage query params. The generated ogen
// decoder rejects values outside the bounds declared by ListQuery in TypeSpec
// before a request reaches a handler.
func newPagination(page, perPage api.OptInt32) pagination {
	return pagination{
		Page:    optOrDefault(page, 1),
		PerPage: optOrDefault(perPage, defaultPerPage),
	}
}

// Offset is the SQL OFFSET for this window.
func (p pagination) Offset() int { return int((p.Page - 1) * p.PerPage) }

// Limit is the SQL LIMIT for this window.
func (p pagination) Limit() int { return int(p.PerPage) }

// optOrDefault returns v when present, otherwise def.
func optOrDefault(v api.OptInt32, def int32) int32 {
	if n, ok := v.Get(); ok {
		return n
	}
	return def
}
