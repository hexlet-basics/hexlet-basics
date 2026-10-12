package eventhandlers_test

import (
	"testing"

	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"hexletbasics/internal/eventhandlers"
	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

// A lead fact becomes one amoCRM job, and redelivery of the same fact does not
// queue a second one: River's uniqueness over the args deduplicates it, so
// RequireInserted (which fails on more than one job) passes only for one.
func TestLeadCreatedEnqueuesTheAmoCRMJobOnce(t *testing.T) {
	pool, err := store.NewDB(testdb.DatabaseURL())
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.ExecContext(t.Context(), "DELETE FROM river_job WHERE kind = $1", jobs.AmoCRMLeadArgs{}.Kind())
		_ = pool.Close()
	})
	client, err := jobs.NewInsertOnlyClient(pool, nil, tracenoop.NewTracerProvider(), metricnoop.NewMeterProvider())
	require.NoError(t, err)

	handler := eventhandlers.LeadCreated(client)
	event := &events.LeadCreated{LeadID: 42, UserID: 7}

	require.NoError(t, handler.Handle(t.Context(), event))
	require.NoError(t, handler.Handle(t.Context(), event))

	job := rivertest.RequireInserted[*riverdatabasesql.Driver](t.Context(), t, riverdatabasesql.New(pool), jobs.AmoCRMLeadArgs{}, nil)
	assert.Equal(t, *event, job.Args.Event)
	assert.Equal(t, "lead_created_to_amocrm_v1", handler.HandlerName())
}
