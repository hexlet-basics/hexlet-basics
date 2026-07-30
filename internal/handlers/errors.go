package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/ogenerrors"
	"go.opentelemetry.io/otel/trace"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/localization"
)

// APIErrorHandler maps application and transport errors to the contract's
// generated RFC 9457 response and reports unexpected failures.
//
// It sits at both ogen error seams:
//
//   - NewError implements ogen's convenient-error method for raw errors returned
//     by application handlers.
//   - Write is installed with api.WithErrorHandler for request decoding,
//     security, and response-encoding failures that occur outside a handler.
//
// Classification and response construction stay here, while the wire schema and
// JSON encoder remain generated from TypeSpec.
type APIErrorHandler struct {
	translator *localization.Translator
	logger     *slog.Logger
}

// fieldValidationError carries a declared field-error response through the same
// writer as transport and application failures.
type fieldValidationError struct {
	errors api.ValidationErrorErrors
}

func (e *fieldValidationError) Error() string {
	return "request validation failed"
}

// statusError assigns an expected HTTP status without teaching the central
// writer about an adapter's dependency-specific sentinel errors.
type statusError struct {
	status int
	err    error
}

func (e *statusError) Error() string {
	return e.err.Error()
}

func (e *statusError) Unwrap() error {
	return e.err
}

func (e *statusError) HTTPStatus() int {
	return e.status
}

func newValidationError(field, message string) error {
	return &fieldValidationError{
		errors: api.ValidationErrorErrors{field: {message}},
	}
}

func withHTTPStatus(status int, err error) error {
	return &statusError{status: status, err: err}
}

// NewAPIErrorHandler builds the central error adapter.
func NewAPIErrorHandler(
	translator *localization.Translator,
	logger *slog.Logger,
) *APIErrorHandler {
	return &APIErrorHandler{
		translator: translator,
		logger:     logger,
	}
}

// NewError creates ogen's generated default response for an error returned by
// an application handler.
func (h *APIErrorHandler) NewError(ctx context.Context, err error) *api.ProblemDetailsStatusCode {
	status := errorStatus(err)
	h.report(ctx, nil, status, err)
	return &api.ProblemDetailsStatusCode{
		StatusCode: status,
		Response:   h.problem(ctx, status),
	}
}

// NewError implements ogen's generated convenient-error seam.
func (s *Server) NewError(ctx context.Context, err error) *api.ProblemDetailsStatusCode {
	return s.errors.NewError(ctx, err)
}

// Write writes contract-shaped errors raised before or after an application
// handler runs. ogen's generated convenient-error encoder owns handler-returned
// errors; this hook covers decode, security, and response-encoding errors.
func (h *APIErrorHandler) Write(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	var validation *fieldValidationError
	if errors.As(err, &validation) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		encoder := new(jx.Encoder)
		response := api.ValidationError{Errors: validation.errors}
		response.Encode(encoder)
		if _, writeErr := encoder.WriteTo(w); writeErr != nil {
			h.logger.ErrorContext(ctx, "writing validation error response", "err", writeErr)
		}
		return
	}

	status := errorStatus(err)
	h.report(ctx, r, status, err)

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	encoder := new(jx.Encoder)
	problem := h.problem(ctx, status)
	problem.Encode(encoder)
	if _, writeErr := encoder.WriteTo(w); writeErr != nil {
		h.logger.ErrorContext(ctx, "writing error response", "err", writeErr)
	}
}

func errorStatus(err error) int {
	var ogenErr ogenerrors.Error
	var statusErr interface{ HTTPStatus() int }
	switch {
	case errors.As(err, &statusErr):
		return statusErr.HTTPStatus()
	case errors.As(err, &ogenErr):
		return ogenErr.Code()
	case ent.IsNotFound(err):
		return http.StatusNotFound
	case ent.IsConstraintError(err):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func (h *APIErrorHandler) problem(ctx context.Context, status int) api.ProblemDetails {
	return api.ProblemDetails{
		Type:   "about:blank",
		Title:  h.translator.StatusText(ctx, status),
		Status: int32(status),
	}
}

func (h *APIErrorHandler) report(ctx context.Context, r *http.Request, status int, err error) {
	if status < http.StatusInternalServerError {
		return
	}

	spanContext := trace.SpanContextFromContext(ctx)
	traceID := spanContext.TraceID().String()
	spanID := spanContext.SpanID().String()
	attrs := []any{
		"err", err,
		"status", status,
		"trace_id", traceID,
		"span_id", spanID,
	}
	tags := map[string]string{
		"trace_id": traceID,
		"span_id":  spanID,
	}
	httpContext := sentry.Context{"status": status}
	if r != nil {
		attrs = append(attrs, "method", r.Method, "path", r.URL.Path)
		tags["http.method"] = r.Method
		tags["http.path"] = r.URL.Path
		httpContext["method"] = r.Method
		httpContext["path"] = r.URL.Path
	}

	h.logger.ErrorContext(ctx, "request failed", attrs...)
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		return
	}
	hub.WithScope(func(scope *sentry.Scope) {
		scope.SetTags(tags)
		scope.SetContext("http", httpContext)
		hub.CaptureException(err)
	})
}

// NewNotFoundHandler localizes requests that do not match an ogen route.
func NewNotFoundHandler(translator *localization.Translator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, translator.StatusText(r.Context(), http.StatusNotFound), http.StatusNotFound)
	}
}

// NewMethodNotAllowedHandler preserves ogen's Allow/OPTIONS behavior while
// localizing the human-readable 405 body.
func NewMethodNotAllowedHandler(translator *localization.Translator) func(http.ResponseWriter, *http.Request, string) {
	return func(w http.ResponseWriter, r *http.Request, allowed string) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Allow", allowed)
		http.Error(w, translator.StatusText(r.Context(), http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
