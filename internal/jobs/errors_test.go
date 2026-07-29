package jobs_test

import (
	"errors"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/jobs"
)

func TestErrorHandlerCapturesJobErrors(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)
	handler := jobs.NewErrorHandler(client)
	job := &rivertype.JobRow{
		ID:          42,
		Kind:        "exercise_loader",
		Queue:       "default",
		Attempt:     2,
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
	assert.Equal(t, "2", events[0].Tags["job.attempt"])
	assert.Equal(t, "25", events[0].Tags["job.max_attempts"])
}

func TestErrorHandlerCapturesJobPanics(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)
	handler := jobs.NewErrorHandler(client)
	job := &rivertype.JobRow{ID: 7, Kind: "ping", Queue: "critical", Attempt: 1, MaxAttempts: 3}

	result := handler.HandlePanic(t.Context(), job, "boom", "goroutine 1 [running]")

	assert.Nil(t, result, "observability must not cancel the job")
	events := transport.Events()
	require.Len(t, events, 1)
	require.Len(t, events[0].Exception, 1)
	assert.Equal(t, "job panic: boom", events[0].Exception[0].Value)
	assert.Equal(t, "goroutine 1 [running]", events[0].Contexts["river"]["panic_trace"])
}
