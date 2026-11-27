package testsupport

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivertest"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/require"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/permission"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
	"hexletbasics/internal/versionbuild"
)

// Services are the production event publisher and job enqueuers, wired over the
// test's savepoint transactor. They write river_job rows inside the test's
// rollback-only transaction, so a test asserts on what an operation left behind
// (RequireEnqueued, EnqueuedJobs) rather than on what was called. A published
// fact is the jobs of its consumers, so it is asserted the same way.
type Services struct {
	// River is the insert-only client behind the enqueuers; EnqueuedJobs lists the
	// queue through its own JobListTx.
	River   *river.Client[*sql.Tx]
	Events  *events.Publisher
	Starter *versionbuild.Starter
	Jobs    *jobs.Enqueuer

	tx *sql.Tx
}

// NewServices builds the services. The River client is insert-only and only
// needs a pool to construct; every insert goes through the transactor's
// transaction (InsertTx), never the pool.
func NewServices(t *testing.T, tx *sql.Tx, transactor store.Transactor) *Services {
	t.Helper()
	pool, err := store.NewDB(testdb.DatabaseURL())
	require.NoError(t, err)
	t.Cleanup(func() { _ = pool.Close() })
	riverClient, err := jobs.NewInsertOnlyClient(
		pool,
		nil,
		tracenoop.NewTracerProvider(),
		metricnoop.NewMeterProvider(),
	)
	require.NoError(t, err)
	enqueuer := jobs.NewEnqueuer(transactor, riverClient)
	return &Services{
		tx:      tx,
		River:   riverClient,
		Events:  events.NewPublisher(transactor, riverClient, jobs.EventConsumers()...),
		Starter: versionbuild.NewStarter(transactor, enqueuer),
		Jobs:    enqueuer,
	}
}

// Grant gives a Staff Role Permissions it did not hold, through the production
// writer behind the roles screen (permission.Permissions). That writer sets a
// role's grants exactly, so the grants the role holds are read back and kept.
// Granting one the role already holds changes nothing.
func Grant(t testing.TB, db *ent.Client, transactor store.Transactor, roleID int, permissions ...string) {
	t.Helper()
	held, err := db.StaffRoleGrant.Query().
		Where(staffrolegrant.RoleID(roleID)).
		Select(staffrolegrant.FieldPermission).
		Strings(t.Context())
	require.NoError(t, err)
	require.NoError(t, permission.NewPermissions(db, transactor).Apply(t.Context(), roleID, append(held, permissions...)))
}

// Grant is Grant for a harness: the role's grants change in the test's
// transaction, where the handlers see them.
func (h *Harness) Grant(t testing.TB, roleID int, permissions ...string) {
	t.Helper()
	Grant(t, h.DB, h.Transactor, roleID, permissions...)
}

// RequireEnqueued asserts the queue holds exactly one job of kind's type and
// returns its args. It is River's own rivertest assertion: it fails when there
// is none or more than one, so a repeated request skipped by the job's
// uniqueness passes and a duplicate fails.
func RequireEnqueued[A river.JobArgs](t *testing.T, s *Services, kind A) A {
	t.Helper()
	job := rivertest.RequireInsertedTx[*riverdatabasesql.Driver](t.Context(), t, s.tx, kind, nil)
	return job.Args
}

// RequireNotEnqueued asserts the queue holds no job of kind's type.
func RequireNotEnqueued[A river.JobArgs](t *testing.T, s *Services, kind A) {
	t.Helper()
	rivertest.RequireNotInsertedTx[*riverdatabasesql.Driver](t.Context(), t, s.tx, kind, nil)
}

// EnqueuedJobs lists every job of one kind on the queue, oldest first, for the
// operations that enqueue several at once (River's JobListTx).
func EnqueuedJobs[A river.JobArgs](t *testing.T, s *Services) []A {
	t.Helper()
	var kind A
	list, err := s.River.JobListTx(t.Context(), s.tx,
		river.NewJobListParams().Kinds(kind.Kind()).OrderBy(river.JobListOrderByID, river.SortOrderAsc))
	require.NoError(t, err)

	found := make([]A, 0, len(list.Jobs))
	for _, job := range list.Jobs {
		var args A
		require.NoError(t, json.Unmarshal(job.EncodedArgs, &args))
		found = append(found, args)
	}
	return found
}

// WorkEnqueued works every available job of A's kind the test has queued,
// oldest first, through worker and in the test's transaction: what the worker
// process does once the request is over. The harness runs no worker, so a test
// that asserts on what a job leaves behind (an Audit Record) calls this after
// the operation and then reads the rows. A job already worked is not available
// any more, so a second call works only what was queued since.
func WorkEnqueued[A river.JobArgs](t *testing.T, s *Services, worker river.Worker[A]) {
	t.Helper()
	var kind A
	list, err := s.River.JobListTx(t.Context(), s.tx,
		river.NewJobListParams().
			Kinds(kind.Kind()).
			States(rivertype.JobStateAvailable).
			OrderBy(river.JobListOrderByID, river.SortOrderAsc))
	require.NoError(t, err)

	runner := rivertest.NewWorker(t, riverdatabasesql.New(nil), &river.Config{}, worker)
	for _, job := range list.Jobs {
		result, err := runner.WorkJob(t.Context(), t, s.tx, job)
		require.NoError(t, err)
		require.Equal(t, river.EventKindJobCompleted, result.EventKind, "the job is done")
	}
}
