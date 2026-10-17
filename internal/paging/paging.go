// Package paging resolves the 1-based page window of the admin lists. It is a
// leaf package (it imports only the generated contract) so both the handlers
// and the ent-generated AdminList helpers can share one implementation.
package paging

import "hexletbasics/internal/api"

// Window is a resolved 1-based page window. It centralizes the offset math
// that every list endpoint would otherwise re-inline.
type Window struct {
	Page    int32
	PerPage int32
}

// New resolves absent query params. The page size has no Go default: the
// contract declares it once (`perPage?: int32 = 20` in TypeSpec) and the
// generated ogen decoder fills it in, so a request reaching a handler always
// carries one. The same decoder rejects values outside the declared bounds.
func New(page, perPage api.OptInt32) Window {
	return Window{
		Page:    page.Or(1),
		PerPage: perPage.Value,
	}
}

// Offset is the SQL OFFSET for this window.
func (w Window) Offset() int { return int((w.Page - 1) * w.PerPage) }

// Limit is the SQL LIMIT for this window.
func (w Window) Limit() int { return int(w.PerPage) }
