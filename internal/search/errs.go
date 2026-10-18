package search

import "github.com/samber/oops"

// errs is the package template for errors (ADR-0018): static messages, the data
// in With, wrapped once at the exported function that knows the operation.
var errs = oops.In("search")
