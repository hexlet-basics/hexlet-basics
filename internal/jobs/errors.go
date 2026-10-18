package jobs

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// ErrorHandler reports worker errors and panics without changing River's retry
// decisions. Returning nil delegates cancellation and scheduling to River.
//
// Only the final attempt of a job goes to Sentry: an earlier one will be
// retried, so it is routine and would page for every try of a flaky call. Every
// attempt, final or not, is logged, so the history stays readable in the logs.
type ErrorHandler struct {
	sentryClient *sentry.Client
	logger       *slog.Logger
}

// NewErrorHandler builds River's native error-handler adapter.
func NewErrorHandler(sentryClient *sentry.Client, logger *slog.Logger) *ErrorHandler {
	return &ErrorHandler{sentryClient: sentryClient, logger: logger}
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
	h.capture(ctx, job, errs.With("panic", panicVal).Errorf("job panic"), trace)
	return nil
}

func (h *ErrorHandler) capture(ctx context.Context, job *rivertype.JobRow, err error, panicTrace string) {
	final := job.Attempt >= job.MaxAttempts
	h.logger.ErrorContext(ctx, "job attempt failed",
		slog.Any("err", err),
		slog.Int64("job_id", job.ID),
		slog.String("job_kind", job.Kind),
		slog.String("job_queue", job.Queue),
		slog.Int("job_attempt", job.Attempt),
		slog.Int("job_max_attempts", job.MaxAttempts),
		slog.Bool("final_attempt", final),
	)
	if !final {
		return
	}
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
	h.sentryClient.CaptureException(err, &sentry.EventHint{Context: ctx, OriginalException: err}, scope)
}
