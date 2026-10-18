package jobs_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river/rivertype"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/config"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/telemetry"
)

func TestErrorHandlerSendsOopsAttributesOfJobErrors(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := telemetry.NewSentryClient(
		&config.Config{SentryDSN: "https://public@example.com/1"},
		telemetry.WithSentryTransport(transport),
	)
	require.NoError(t, err)
	handler := jobs.NewErrorHandler(client, slog.New(slog.DiscardHandler))
	job := &rivertype.JobRow{ID: 1, Kind: "exercise_loader", Queue: "default", Attempt: 3, MaxAttempts: 3}

	handler.HandleError(t.Context(), job, oops.In("courseloader").With("course_id", 7).Wrapf(errors.New("clone failed"), "load course"))

	events := transport.Events()
	require.Len(t, events, 1)
	assert.Equal(t, "courseloader", events[0].Tags["oops.domain"])
	assert.Equal(t, 7, events[0].Contexts["oops"]["course_id"])
}

func TestErrorHandlerCapturesJobErrors(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)
	handler := jobs.NewErrorHandler(client, slog.New(slog.DiscardHandler))
	job := &rivertype.JobRow{
		ID:          42,
		Kind:        "exercise_loader",
		Queue:       "default",
		Attempt:     25,
		MaxAttempts: 25,
	}

	result := handler.HandleError(t.Context(), job, errors.New("clone failed"))

	assert.Nil(t, result, "observability must not change River retry behavior")
	events := transport.Events()
	require.Len(t, events, 1)
	require.Len(t, events[0].Exception, 1)
	assert.Equal(t, "clone failed", events[0].Exception[0].Value)
	assert.Equal(t, "42", events[0].Tags["job.id"])
	assert.Equal(t, "exercise_loader", events[0].Tags["job.kind"])
	assert.Equal(t, "default", events[0].Tags["job.queue"])
	assert.Equal(t, "25", events[0].Tags["job.attempt"])
	assert.Equal(t, "25", events[0].Tags["job.max_attempts"])
}

func TestErrorHandlerCapturesJobPanics(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)
	handler := jobs.NewErrorHandler(client, slog.New(slog.DiscardHandler))
	job := &rivertype.JobRow{ID: 7, Kind: "ping", Queue: "critical", Attempt: 3, MaxAttempts: 3}

	result := handler.HandlePanic(t.Context(), job, "boom", "goroutine 1 [running]")

	assert.Nil(t, result, "observability must not cancel the job")
	events := transport.Events()
	require.Len(t, events, 1)
	// The panic value travels as an attribute (read by telemetry's BeforeSend), not in the message.
	require.NotEmpty(t, events[0].Exception)
	assert.Equal(t, "job panic", events[0].Exception[len(events[0].Exception)-1].Value)
	assert.Equal(t, "goroutine 1 [running]", events[0].Contexts["river"]["panic_trace"])
}

// A job that will be retried is routine; paging someone for every attempt of a
// flaky call would bury the one event that says the job is out of tries.
func TestErrorHandlerSendsOnlyTheFinalAttemptToSentry(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)
	var logs bytes.Buffer
	handler := jobs.NewErrorHandler(client, slog.New(slog.NewJSONHandler(&logs, nil)))
	job := &rivertype.JobRow{ID: 9, Kind: "account_email", Queue: "mail", Attempt: 1, MaxAttempts: 3}

	handler.HandleError(t.Context(), job, errors.New("smtp down"))

	assert.Empty(t, transport.Events(), "a retry is not an incident")
	assert.Contains(t, logs.String(), "smtp down", "the attempt still reaches the logs")
	assert.Contains(t, logs.String(), `"job_kind":"account_email"`)

	job.Attempt = 3
	handler.HandleError(t.Context(), job, errors.New("smtp down"))

	require.Len(t, transport.Events(), 1)
}

func TestErrorHandlerDoesNotSendAPanicOfANonFinalAttempt(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)
	handler := jobs.NewErrorHandler(client, slog.New(slog.DiscardHandler))

	handler.HandlePanic(t.Context(), &rivertype.JobRow{ID: 1, Kind: "k", Attempt: 1, MaxAttempts: 2}, "boom", "trace")

	assert.Empty(t, transport.Events())
}
