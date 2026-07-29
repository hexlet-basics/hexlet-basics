package jobs

import (
	"context"
	"fmt"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// ErrorHandler reports worker errors and panics without changing River's retry
// decisions. Returning nil delegates cancellation and scheduling to River.
type ErrorHandler struct {
	sentryClient *sentry.Client
}

// NewErrorHandler builds River's native error-handler adapter.
func NewErrorHandler(sentryClient *sentry.Client) *ErrorHandler {
	return &ErrorHandler{sentryClient: sentryClient}
}

// HandleError captures a failed job attempt.
func (h *ErrorHandler) HandleError(
	ctx context.Context,
	job *rivertype.JobRow,
	err error,
) *river.ErrorHandlerResult {
	h.capture(ctx, job, err, "")
	return nil
}

// HandlePanic captures the recovered panic value and River's stack trace.
func (h *ErrorHandler) HandlePanic(
	ctx context.Context,
	job *rivertype.JobRow,
	panicVal any,
	trace string,
) *river.ErrorHandlerResult {
	h.capture(ctx, job, fmt.Errorf("job panic: %v", panicVal), trace)
	return nil
}

func (h *ErrorHandler) capture(ctx context.Context, job *rivertype.JobRow, err error, panicTrace string) {
	scope := sentry.NewScope()
	scope.SetTags(map[string]string{
		"job.id":           strconv.FormatInt(job.ID, 10),
		"job.kind":         job.Kind,
		"job.queue":        job.Queue,
		"job.attempt":      strconv.Itoa(job.Attempt),
		"job.max_attempts": strconv.Itoa(job.MaxAttempts),
	})
	riverContext := sentry.Context{
		"id":           job.ID,
		"kind":         job.Kind,
		"queue":        job.Queue,
		"attempt":      job.Attempt,
		"max_attempts": job.MaxAttempts,
	}
	if panicTrace != "" {
		riverContext["panic_trace"] = panicTrace
	}
	scope.SetContext("river", riverContext)
	h.sentryClient.CaptureException(err, &sentry.EventHint{Context: ctx}, scope)
}
