package testsupport

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"testing"

	wmsql "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivertest"
	"github.com/stretchr/testify/require"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/staffroles"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
	"hexletbasics/internal/versionbuilds"
)

// Services are the production event publisher and job enqueuers, wired over the
// test's savepoint transactor. They write the outbox and river_job rows inside
// the test's rollback-only transaction, so a test asserts on what an operation
// left behind (PublishedEvents, EnqueuedJobs) rather than on what was called.
type Services struct {
	// River is the insert-only client behind the enqueuers; EnqueuedJobs lists the
	// queue through its own JobListTx.
	River   *river.Client[*sql.Tx]
	Events  *events.Publisher
	Starter *versionbuilds.Starter
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
		Events:  events.NewPublisher(transactor, slog.New(slog.DiscardHandler)),
		Starter: versionbuilds.NewStarter(transactor, enqueuer),
		Jobs:    enqueuer,
	}
}

// Grant gives a Staff Role Permissions it did not hold, through the production
// writer behind the roles screen (staffroles.Permissions). That writer sets a
// role's grants exactly, so the grants the role holds are read back and kept.
// Granting one the role already holds changes nothing.
func Grant(t testing.TB, db *ent.Client, transactor store.Transactor, roleID int, permissions ...string) {
	t.Helper()
	held, err := db.StaffRoleGrant.Query().
		Where(staffrolegrant.RoleID(roleID)).
		Select(staffrolegrant.FieldPermission).
		Strings(t.Context())
	require.NoError(t, err)
	require.NoError(t, staffroles.NewPermissions(db, transactor).Apply(t.Context(), roleID, append(held, permissions...)))
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

// PublishedEvents lists the domain facts of one type the test has written to
// the outbox, oldest first. Watermill's own subscriber query hides rows of the
// still-open transaction, so the topic table is read directly, but through the
// schema's row mapping and the events' marshaler.
func PublishedEvents[E events.Event](t *testing.T, s *Services) []E {
	t.Helper()
	var zero E
	name := events.Marshaler.Name(zero)

	var found []E
	for _, msg := range outbox(t, s) {
		if events.Marshaler.NameFromMessage(msg) != name {
			continue
		}
		var event E
		require.NoError(t, events.Marshaler.Unmarshal(msg, &event))
		found = append(found, event)
	}
	return found
}

// PublishedCount is how many domain facts of any type the test has written.
func PublishedCount(t *testing.T, s *Services) int {
	t.Helper()
	return len(outbox(t, s))
}

// outbox reads the messages of the domain-events topic table in order.
func outbox(t *testing.T, s *Services) []*message.Message {
	t.Helper()
	schema := wmsql.DefaultPostgreSQLSchema{}
	rows, err := s.tx.QueryContext(t.Context(),
		`SELECT "offset", transaction_id, uuid, payload, metadata FROM `+
			schema.MessagesTable(events.Topic)+` ORDER BY transaction_id, "offset"`)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	var messages []*message.Message
	for rows.Next() {
		row, err := schema.UnmarshalMessage(wmsql.UnmarshalMessageParams{Row: rows})
		require.NoError(t, err)
		messages = append(messages, row.Msg)
	}
	require.NoError(t, rows.Err())
	return messages
}

// Deliver hands every fact in the outbox that handler consumes to it, oldest
// first, the way the worker's Watermill router does after the request is over.
// The harness runs no router, so a test that asserts on what a consumer leaves
// behind (an Audit Record) calls this after the operation and then reads the
// rows. Facts of other types are skipped.
func Deliver(t *testing.T, s *Services, handler cqrs.EventHandler) {
	t.Helper()
	name := events.Marshaler.Name(handler.NewEvent())
	for _, msg := range outbox(t, s) {
		if events.Marshaler.NameFromMessage(msg) != name {
			continue
		}
		event := handler.NewEvent()
		require.NoError(t, events.Marshaler.Unmarshal(msg, event))
		require.NoError(t, handler.Handle(t.Context(), event))
	}
}
