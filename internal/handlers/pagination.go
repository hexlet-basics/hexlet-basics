package handlers

import "hexletbasics/internal/paging"

// newPagination resolves the page window from the list query params; the
// implementation lives in internal/paging so the ent-generated AdminList
// helpers share it.
var newPagination = paging.New
