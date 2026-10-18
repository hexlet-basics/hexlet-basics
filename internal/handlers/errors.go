package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/ogenerrors"
	"go.opentelemetry.io/otel/trace"

	"hexletbasics/internal/account"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/selfprotection"
	"hexletbasics/internal/transition"
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
		Response:   h.problem(ctx, status, err),
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
	status := errorStatus(err)
	h.report(ctx, r, status, err)

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	encoder := new(jx.Encoder)
	problem := h.problem(ctx, status, err)
	problem.Encode(encoder)
	if _, writeErr := encoder.WriteTo(w); writeErr != nil {
		h.logger.ErrorContext(ctx, "writing error response", "err", writeErr)
	}
}

// statusClientClosedRequest is nginx's 499: the caller went away (closed the
// tab, navigated off a lesson) and cancelled the request context. Nobody reads
// the response, and it is not a server failure, so it stays out of the error
// log and Sentry (report only looks at 5xx).
const statusClientClosedRequest = 499

// domainStatus is the one table from a domain error to an HTTP status. An
// adapter-specific failure joins it by wrapping one of these (a self-protection
// refusal wraps ErrForbidden, the asset store's missing key wraps ErrNotFound) instead of
// carrying a status of its own, and ent's errors reach it already translated
// by domain.FromStore.
var domainStatus = []struct {
	err    error
	status int
}{
	{domain.ErrNotFound, http.StatusNotFound},
	{domain.ErrConflict, http.StatusConflict},
	{domain.ErrInvalid, http.StatusUnprocessableEntity},
	{domain.ErrForbidden, http.StatusForbidden},
}

func errorStatus(err error) int {
	// The domain's classes come first: ogen wraps what a security handler
	// returns in its own 401, and a refusal of a signed-in caller (ErrForbidden)
	// must stay a 403.
	translated := domain.FromStore(err)
	for _, row := range domainStatus {
		if errors.Is(translated, row.err) {
			return row.status
		}
	}
	if ogenErr, ok := errors.AsType[ogenerrors.Error](err); ok {
		return ogenErr.Code()
	}
	if errors.Is(err, context.Canceled) {
		return statusClientClosedRequest
	}
	return http.StatusInternalServerError
}

// problem builds the response body. `detail` is a localized explanation that
// detailFor picks for the errors the adapter knows (a refusal, an illegal move,
// a named not-found); every other error stays a bare status, so nothing
// internal is exposed. A 422 adds `errors`, the field messages the caller can
// fix.
func (h *APIErrorHandler) problem(ctx context.Context, status int, err error) api.ProblemDetails {
	problem := api.ProblemDetails{
		Type:   "about:blank",
		Title:  h.translator.StatusText(ctx, status),
		Status: apiconv.Int32FromInt(status),
	}
	if message, data, ok := detailFor(err); ok {
		translated := make(map[string]any, len(data))
		for key, value := range data {
			if named, isMessage := value.(localization.Message); isMessage {
				value = h.translator.Text(ctx, named)
			}
			translated[key] = value
		}
		problem.Detail = api.NewOptString(h.translator.TextWith(ctx, message, translated))
	}
	if status == http.StatusUnprocessableEntity {
		if fields := h.fieldErrors(ctx, status, err); len(fields) > 0 {
			problem.Errors = api.NewOptProblemDetailsErrors(fields)
		}
	}
	return problem
}

// fieldErrors is the `errors` extension of a 422: the field an invalid input
// names and what is wrong with it, in the caller's language. An input the
// domain marked invalid without a message of its own gets the status text.
func (h *APIErrorHandler) fieldErrors(ctx context.Context, status int, err error) api.ProblemDetailsErrors {
	if input, ok := errors.AsType[*invalidInputError](err); ok {
		if input.text != "" {
			return api.ProblemDetailsErrors{input.field: {input.text}}
		}
		return api.ProblemDetailsErrors{input.field: {h.translator.Text(ctx, input.message)}}
	}
	if field, ok := domain.FieldOf(err); ok {
		return api.ProblemDetailsErrors{field: {h.translator.StatusText(ctx, status)}}
	}
	return nil
}

// invalidInputError is a 422 whose field and localized message the handler knows. It
// is a domain.ErrInvalid, so the one errorStatus table answers it.
type invalidInputError struct {
	field   string
	message localization.Message
	// text is the message when it is not a localized one: names of the caller's
	// own input echoed back (an unknown permission), which no catalog holds.
	text string
}

func (e *invalidInputError) Error() string { return domain.ErrInvalid.Error() + ": field " + e.field }

func (e *invalidInputError) Unwrap() error { return domain.ErrInvalid }

// invalid reports input the caller can fix: field names the request property
// and message says what is wrong with it.
func invalid(field string, message localization.Message) error {
	return &invalidInputError{field: field, message: message}
}

// invalidText is invalid for a message that echoes the caller's own input and
// has no translation.
func invalidText(field, text string) error {
	return &invalidInputError{field: field, text: text}
}

// notFoundError is a domain.ErrNotFound that says what was not found, so details can
// give it an explanation of its own.
type notFoundError string

func (n notFoundError) Error() string { return string(n) + ": " + domain.ErrNotFound.Error() }

// Is makes every notFoundError match domain.ErrNotFound.
func (notFoundError) Is(target error) bool { return target == domain.ErrNotFound }

// Not-founds that answer with their own explanation instead of the bare status
// text. Each is a domain.ErrNotFound, listed in details.
var (
	errLessonNotFound       error = notFoundError("lesson")
	errCourseNotFound       error = notFoundError("course")
	errEmailLinkInvalid     error = notFoundError("email link")
	errOAuthRequestNotFound error = notFoundError("oauth request")
	errConnectedAppNotFound error = notFoundError("connected app")
)

// details is the table from an error to the localized explanation the client
// may show as `detail`. It lives in the adapter, not in the domain, so no
// module imports localization to word its own failure.
var details = []struct {
	err     error
	message localization.Message
}{
	{errLessonNotFound, localization.LessonNotFound},
	{errCourseNotFound, localization.CourseNotFound},
	{errEmailLinkInvalid, localization.EmailLinkInvalid},
	{errOAuthRequestNotFound, localization.OAuthRequestNotFound},
	{errConnectedAppNotFound, localization.ConnectedAppNotFound},
	{selfprotection.ErrOwnStaffRecord, localization.OwnStaffRecord},
	{selfprotection.ErrOwnRoleGrants, localization.OwnRoleGrants},
	{selfprotection.ErrAdministratorFlagRestricted, localization.AdministratorFlagRestricted},
	{selfprotection.ErrOwnAccountDeletion, localization.OwnAccountDeletion},
	{account.ErrLastAdmin, localization.LastAdministrator},
}

// stateNames are the localized names of the states in each vocabulary
// (transition.Vocabulary): Publication States (GLOSSARY) are shared by every
// resource that has them, a Course's Readiness has a vocabulary of its own.
var stateNames = map[string]map[string]localization.Message{
	"publication_state": {
		"draft":     localization.PublicationStateDraft,
		"published": localization.PublicationStatePublished,
		"archived":  localization.PublicationStateArchived,
	},
	"course_readiness": {
		"draft":          localization.ReadinessDraft,
		"in_development": localization.ReadinessInDevelopment,
		"completed":      localization.ReadinessCompleted,
	},
}

// detailFor finds the localized explanation of err, and the values for its
// placeholders (a value that is itself a localization.Message is translated
// first). A refused state move names both states in their vocabulary.
func detailFor(err error) (localization.Message, map[string]any, bool) {
	if illegal, ok := transition.IsIllegal(err); ok {
		return localization.IllegalStateTransition, map[string]any{
			"From": stateName(illegal.Vocabulary, illegal.From),
			"To":   stateName(illegal.Vocabulary, illegal.To),
		}, true
	}
	for _, row := range details {
		if errors.Is(err, row.err) {
			return row.message, nil, true
		}
	}
	return localization.Message{}, nil, false
}

// stateName is the localized name of a state, or the raw value of one the
// vocabulary does not know (a legacy row with a state of its own).
func stateName(vocabulary, state string) any {
	if message, ok := stateNames[vocabulary][state]; ok {
		return message
	}
	return state
}

func (h *APIErrorHandler) report(ctx context.Context, r *http.Request, status int, err error) {
	if status < http.StatusInternalServerError {
		return
	}

	spanContext := trace.SpanContextFromContext(ctx)
	traceID := spanContext.TraceID().String()
	spanID := spanContext.SpanID().String()
	// The logger adds trace_id and span_id itself (internal/logging); the Sentry
	// tags below stay, since Sentry has no access to the log pipeline.
	attrs := []any{
		"err", err,
		"status", status,
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
