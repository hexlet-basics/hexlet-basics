package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"hexletbasics/internal/api"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/telemetry"
)

type failingCoursesHandler struct {
	api.UnimplementedHandler
	errorHandler *handlers.APIErrorHandler
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func (h *failingCoursesHandler) ListCourses(context.Context) ([]api.CourseCatalogItem, error) {
	return nil, errors.New("database unavailable")
}

func (h *failingCoursesHandler) NewError(ctx context.Context, err error) *api.ProblemDetailsStatusCode {
	return h.errorHandler.NewError(ctx, err)
}

func TestAPIErrorHandlerReportsUnexpectedErrors(t *testing.T) {
	translator, err := localization.New()
	require.NoError(t, err)

	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)

	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	ctx, span := tracerProvider.Tracer("test").Start(t.Context(), "request")
	defer span.End()
	ctx = sentry.SetHubOnContext(ctx, sentry.NewHub(client, sentry.NewScope()))

	handler := handlers.NewAPIErrorHandler(translator, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/admin/languages", nil).WithContext(ctx)
	cause := errors.New("database unavailable")
	handler.Write(ctx, recorder, request, cause)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "application/problem+json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
		"type":"about:blank",
		"title":"Internal Server Error",
		"status":500
	}`, recorder.Body.String())
	assert.Contains(t, logs.String(), `"msg":"request failed"`)
	assert.Contains(t, logs.String(), `"method":"POST"`)
	assert.Contains(t, logs.String(), `"path":"/admin/languages"`)
	assert.Contains(t, logs.String(), `"trace_id":"`+span.SpanContext().TraceID().String()+`"`)

	events := transport.Events()
	require.Len(t, events, 1)
	require.Len(t, events[0].Exception, 1)
	assert.Equal(t, cause.Error(), events[0].Exception[0].Value)
	assert.Equal(t, span.SpanContext().TraceID().String(), events[0].Tags["trace_id"])
	assert.Equal(t, span.SpanContext().SpanID().String(), events[0].Tags["span_id"])
}

func TestAPIErrorHandlerDoesNotReportExpectedErrors(t *testing.T) {
	translator, err := localization.New()
	require.NoError(t, err)

	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)

	handler := handlers.NewAPIErrorHandler(translator, logger)
	recorder := httptest.NewRecorder()
	ctx := sentry.SetHubOnContext(t.Context(), sentry.NewHub(client, sentry.NewScope()))
	request := httptest.NewRequest(http.MethodGet, "/api/languages/missing", nil).WithContext(ctx)
	handler.Write(request.Context(), recorder, request, &ogenerrors.DecodeParamsError{
		OperationContext: ogenerrors.OperationContext{Name: "ListCourses", ID: "listCourses"},
		Err:              errors.New("invalid query"),
	})

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, "application/problem+json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
		"type":"about:blank",
		"title":"Bad Request",
		"status":400
	}`, recorder.Body.String())
	assert.Empty(t, logs.String())
	assert.Empty(t, transport.Events())
}

func TestGeneratedClientDecodesCentralErrorsAsProblemDetails(t *testing.T) {
	translator, err := localization.New()
	require.NoError(t, err)
	transport := &sentry.MockTransport{}
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)
	errorHandler := handlers.NewAPIErrorHandler(
		translator,
		slog.New(slog.NewJSONHandler(&bytes.Buffer{}, nil)),
	)
	server, err := api.NewServer(
		&failingCoursesHandler{errorHandler: errorHandler},
		api.WithErrorHandler(errorHandler.Write),
	)
	require.NoError(t, err)
	client, err := api.NewClient("http://test", api.WithClient(&http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			recorder := httptest.NewRecorder()
			telemetry.NewSentryHTTPHandler(sentryClient, server).ServeHTTP(recorder, request)
			return recorder.Result(), nil
		}),
	}))
	require.NoError(t, err)

	_, err = client.ListCourses(t.Context())

	var problem *api.ProblemDetailsStatusCode
	require.ErrorAs(t, err, &problem)
	assert.Equal(t, http.StatusInternalServerError, problem.StatusCode)
	assert.Equal(t, api.ProblemDetails{
		Type:   "about:blank",
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
	}, problem.Response)
	events := transport.Events()
	require.Len(t, events, 1)
	assert.Equal(t, http.MethodGet, events[0].Request.Method)
	assert.Equal(t, "http://test/languages", events[0].Request.URL)
}

var _ trace.TracerProvider = (*sdktrace.TracerProvider)(nil)
