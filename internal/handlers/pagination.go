package handlers

import "hexletbasics/internal/api"

// pagination is a resolved 1-based page window. It centralizes the offset math
// that every list endpoint would otherwise re-inline.
type pagination struct {
	Page    int32
	PerPage int32
}

// newPagination resolves absent query params. The page size has no Go default:
// the contract declares it once (`perPage?: int32 = 20` in TypeSpec) and the
// generated ogen decoder fills it in, so a request reaching a handler always
// carries one. The same decoder rejects values outside the declared bounds.
func newPagination(page, perPage api.OptInt32) pagination {
	return pagination{
		Page:    page.Or(1),
		PerPage: perPage.Value,
	}
}

// Offset is the SQL OFFSET for this window.
func (p pagination) Offset() int { return int((p.Page - 1) * p.PerPage) }

// Limit is the SQL LIMIT for this window.
func (p pagination) Limit() int { return int(p.PerPage) }
