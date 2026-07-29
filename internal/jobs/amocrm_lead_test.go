package jobs

import (
	"context"
	"errors"
	"testing"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/events"
)

func TestAmoCRMLeadWorkerCancelsPermanentFailures(t *testing.T) {
	cause := classifiedLeadError{retryable: false}
	worker := &amoCRMLeadWorker{syncer: leadSyncerStub{err: cause}}

	err := worker.Work(t.Context(), &river.Job[AmoCRMLeadArgs]{})

	var cancelErr *river.JobCancelError
	require.ErrorAs(t, err, &cancelErr)
	assert.ErrorIs(t, err, cause)
}

func TestAmoCRMLeadWorkerReturnsRetryableFailures(t *testing.T) {
	cause := classifiedLeadError{retryable: true}
	worker := &amoCRMLeadWorker{syncer: leadSyncerStub{err: cause}}

	err := worker.Work(t.Context(), &river.Job[AmoCRMLeadArgs]{})

	assert.ErrorIs(t, err, cause)
	var cancelErr *river.JobCancelError
	assert.False(t, errors.As(err, &cancelErr))
}

func TestAmoCRMLeadWorkerRetriesUnclassifiedFailures(t *testing.T) {
	cause := errors.New("connection reset")
	worker := &amoCRMLeadWorker{syncer: leadSyncerStub{err: cause}}

	err := worker.Work(t.Context(), &river.Job[AmoCRMLeadArgs]{})

	assert.ErrorIs(t, err, cause)
}

type leadSyncerStub struct {
	err error
}

func (s leadSyncerStub) CreateLead(context.Context, events.LeadCreated) error {
	return s.err
}

type classifiedLeadError struct {
	retryable bool
}

func (e classifiedLeadError) Error() string {
	return "amoCRM rejected lead"
}

func (e classifiedLeadError) Retryable() bool {
	return e.retryable
}
