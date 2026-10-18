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

// Faults a test injects; each is told apart from the others by identity.
var (
	errCause = errors.New("connection reset")
)

func TestAmoCRMLeadWorkerCancelsPermanentFailures(t *testing.T) {
	cause := classifiedLeadError{retryable: false}
	worker := &amoCRMLeadWorker{
		syncer:         leadSyncerStub{err: cause},
		WorkerDefaults: river.WorkerDefaults[AmoCRMLeadArgs]{},
	}

	err := worker.Work(t.Context(), &river.Job[AmoCRMLeadArgs]{})

	var cancelErr *river.JobCancelError
	require.ErrorAs(t, err, &cancelErr)
	require.ErrorIs(t, err, cause)
}

func TestAmoCRMLeadWorkerReturnsRetryableFailures(t *testing.T) {
	cause := classifiedLeadError{retryable: true}
	worker := &amoCRMLeadWorker{
		syncer:         leadSyncerStub{err: cause},
		WorkerDefaults: river.WorkerDefaults[AmoCRMLeadArgs]{},
	}

	err := worker.Work(t.Context(), &river.Job[AmoCRMLeadArgs]{})

	require.ErrorIs(t, err, cause)

	var cancelErr *river.JobCancelError
	assert.NotErrorAs(t, err, &cancelErr)
}

func TestAmoCRMLeadWorkerRetriesUnclassifiedFailures(t *testing.T) {
	worker := &amoCRMLeadWorker{
		syncer:         leadSyncerStub{err: errCause},
		WorkerDefaults: river.WorkerDefaults[AmoCRMLeadArgs]{},
	}

	err := worker.Work(t.Context(), &river.Job[AmoCRMLeadArgs]{})

	require.ErrorIs(t, err, errCause)
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
