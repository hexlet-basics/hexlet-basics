package eventhandler_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivertest"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/eventhandler"
	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport"
)

// txInserter queues jobs in the test's transaction, so they land where the test
// reads them and are discarded with it.
type txInserter struct {
	client *river.Client[*sql.Tx]
	tx     *sql.Tx
}

func (i txInserter) Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	return i.client.InsertTx(ctx, i.tx, args, opts)
}

// A lead fact becomes one amoCRM job, and redelivery of the same fact does not
// queue a second one: River's uniqueness over the args deduplicates it, so
// RequireInsertedTx (which fails on more than one job) passes only for one.
func TestLeadCreatedEnqueuesTheAmoCRMJobOnce(t *testing.T) {
	_, tx := testsupport.NewClientWithTx(t)
	client, err := jobs.NewInsertOnlyClient(nil, nil, nil, nil)
	require.NoError(t, err)

	handler := eventhandler.LeadCreated(txInserter{client: client, tx: tx})
	event := &events.LeadCreated{LeadID: 42, UserID: 7}

	require.NoError(t, handler.Handle(t.Context(), event))
	require.NoError(t, handler.Handle(t.Context(), event))

	job := rivertest.RequireInsertedTx[*riverdatabasesql.Driver](t.Context(), t, tx, jobs.AmoCRMLeadArgs{}, nil)
	assert.Equal(t, *event, job.Args.Event)
	assert.Equal(t, "lead_created_to_amocrm_v1", handler.HandlerName())
}
