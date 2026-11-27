package events_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport"
)

// Faults a test injects; each is told apart from the others by identity.
var errQueueUnavailable = errors.New("queue unavailable")

func staffRead() events.StaffReadPersonalData {
	actor := 5

	return events.StaffReadPersonalData{
		ActorUserID:    &actor,
		Channel:        "admin",
		ConnectedAppID: nil,
		TokenID:        nil,
		Operation:      "AdminListUsers",
		EntityType:     "User",
		EntityIDs:      []string{"7"},
		Count:          1,
		RequestID:      "req-publish-1",
		Address:        "",
		UserAgent:      "",
		Locale:         "en",
		OccurredAt:     time.Now(),
	}
}

// A fact becomes the job of each consumer that reacts to it, in the caller's
// transaction, and of no other: the amoCRM consumer does not react to a read.
func TestAFactQueuesTheJobOfEveryConsumerThatReactsToIt(t *testing.T) {
	_, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	services := testsupport.NewServices(t, tx, transactor)
	read := staffRead()

	require.NoError(t, services.Events.PublishStandalone(t.Context(), read))

	queued := testsupport.EnqueuedJobs[jobs.RecordPersonalDataReadArgs](t, services)
	require.Len(t, queued, 1)
	assert.Equal(t, read.RequestID, queued[0].Event.RequestID)
	assert.Empty(t, testsupport.EnqueuedJobs[jobs.AmoCRMLeadArgs](t, services))
}

type failingQueue struct{}

func (failingQueue) InsertTx(context.Context, *sql.Tx, river.JobArgs, *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	return nil, errQueueUnavailable
}

// A job that cannot be queued fails the publish, so the caller's transaction
// rolls back with the state change it was part of.
func TestAFactWhoseJobCannotBeQueuedFailsThePublish(t *testing.T) {
	_, transactor := testsupport.NewClientWithTransactor(t)
	publisher := events.NewPublisher(transactor, failingQueue{}, jobs.EventConsumers()...)

	err := publisher.PublishStandalone(t.Context(), staffRead())

	require.ErrorIs(t, err, errQueueUnavailable)
}
