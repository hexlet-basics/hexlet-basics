package handlers

import (
	"context"
	"fmt"
	"net/http"

	"hexletbasics/ent"
)

// APIErrorHandler is the central ogen ErrorHandler. It maps ent's error
// taxonomy to HTTP status codes so every handler can return the raw ent error
// (`return nil, err`) instead of assembling a typed error DTO per operation.
//
// This is the one place that decides what a database failure means to a client:
//
//   - ent.IsNotFound        -> 404: a Get/Update/Delete against a missing row.
//   - ent.IsConstraintError -> 409: a DB unique-index violation (the uniqueness
//     rules once checked by hand in the handler now live as constraints).
//   - anything else         -> 500.
//
// Wired via api.WithErrorHandler in the DI container and reused by the test
// harness, so tests exercise the exact same mapping the server runs.
func APIErrorHandler(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	status := http.StatusInternalServerError
	switch {
	case ent.IsNotFound(err):
		status = http.StatusNotFound
	case ent.IsConstraintError(err):
		status = http.StatusConflict
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// The status code carries the meaning; the body is a minimal, uniform
	// envelope. ogen's generated client treats these undeclared statuses as
	// errors, which the handler tests assert on via the recorded status code.
	_, _ = fmt.Fprintf(w, `{"error":%q}`, http.StatusText(status))
}
