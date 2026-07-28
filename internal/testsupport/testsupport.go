// Package testsupport provides Rails-style integration-test plumbing: each test
// runs against a real Postgres transaction that is rolled back on cleanup
// (go-txdb), over a baseline prepared once by `make test-prepare` — the atlas
// migrations plus the fixtures/ YAML loaded by the testfixtures CLI. Handlers
// hit a real database; every test's writes are discarded on rollback, so the
// shared baseline is never mutated and nothing is left behind.
//
// Handlers are exercised end-to-end through the generated ogen client against an
// in-process server (NewHarness) — no TCP listener is opened. This is the Go
// equivalent of a Rails controller test: dispatch through the real HTTP request
// pipeline (routing, codec, the central ErrorHandler) and assert the response
// status code, then assert database state through the same transaction.
package testsupport

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	txdb "github.com/DATA-DOG/go-txdb"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/riverqueue/river"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/config"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/jobs"
)

// testConfig gives handlers fixed public hosts so URL-building assertions are
// deterministic (independent of the ambient env a `config.Load` would read).
var testConfig = &config.Config{
	AppHost:   "code-basics.com",
	PublicURL: "http://localhost:3001",
}

const defaultTestDSN = "postgres://postgres:postgres@127.0.0.1:54330/code_basics_test"

var registerOnce sync.Once

func testDSN() string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	return defaultTestDSN
}

// NewClient opens an ent client bound to a fresh go-txdb transaction that is
// rolled back when the test finishes. Fixtures are not loaded here: they are the
// pre-loaded baseline from `make test-prepare`, which every test sees and whose
// own writes vanish on rollback. Run `make test-prepare` before `go test`.
func NewClient(t *testing.T) *ent.Client {
	t.Helper()

	registerOnce.Do(func() {
		txdb.Register("txdb", "pgx", testDSN())
	})

	db, err := sql.Open("txdb", t.Name())
	if err != nil {
		t.Fatalf("open txdb: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() }) // rolls the transaction back

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}

// Harness bundles the generated API client, wired to an in-process server, with
// the ent client behind it — both over the same rolled-back txdb transaction, so
// the client's writes are visible to DB assertions and discarded together.
//
// One-transaction caveat: the whole test shares a single txdb transaction with no
// per-request savepoint, so a write that violates a DB constraint aborts it
// (Postgres 25P02). After an intentional conflict you therefore cannot query the
// DB, nor issue a further "then a valid write" call — assert the conflict's status
// code and stop. In production each request runs in its own autocommit, so this is
// a test-harness artifact, not a behavioural one.
type Harness struct {
	// Client issues typed calls (URLs and bodies are generated, never hand-written).
	Client *api.Client
	// DB queries the same transaction the handlers wrote through, for assertions.
	DB *ent.Client
	// Enqueuer records background jobs the handlers scheduled, so a test can
	// assert an operation both wrote to the DB and enqueued its follow-up work.
	Enqueuer *RecordingEnqueuer
	doer     *inProcessDoer
}

// LastStatus is the HTTP status of the most recent client call. Error responses
// (404/409) are undeclared statuses the generated client surfaces as an error;
// assert the code here instead of destructuring a typed error member.
func (h *Harness) LastStatus() int { return h.doer.status }

// NewHarness builds the in-process test stack: an ent client over a fresh txdb
// transaction, the handlers.Server, the ogen api.Server with the production
// ErrorHandler, and a generated client whose transport dispatches straight into
// that server's ServeHTTP — so tests run the full request pipeline without a
// socket, and exercise the exact ent-error -> status mapping the server uses.
func NewHarness(t *testing.T) *Harness {
	t.Helper()

	db := NewClient(t)

	enqueuer := &RecordingEnqueuer{DB: db}
	srv, err := api.NewServer(handlers.NewServer(db, testConfig, enqueuer), api.WithErrorHandler(handlers.APIErrorHandler))
	if err != nil {
		t.Fatalf("new api server: %v", err)
	}

	doer := &inProcessDoer{server: srv}
	client, err := api.NewClient("http://test", api.WithClient(doer))
	if err != nil {
		t.Fatalf("new api client: %v", err)
	}

	return &Harness{Client: client, DB: db, doer: doer, Enqueuer: enqueuer}
}

// RecordingEnqueuer is a test adapter for handlers.VersionBuildStarter. It
// creates the version through the harness's rollback-only ent client and records
// the job args without touching River.
type RecordingEnqueuer struct {
	DB       *ent.Client
	Inserted []river.JobArgs
}

// Start mirrors the production operation's visible result.
func (e *RecordingEnqueuer) Start(ctx context.Context, courseID int) (*ent.CourseVersion, error) {
	version, err := e.DB.CourseVersion.Create().
		SetLanguageID(courseID).
		SetState("created").
		Save(ctx)
	if err != nil {
		return nil, err
	}
	args := jobs.ExerciseLoaderArgs{VersionID: version.ID}
	e.Inserted = append(e.Inserted, args)
	return version, nil
}

// inProcessDoer satisfies ogen's http client interface by serving each request
// straight through the server's ServeHTTP against an in-memory recorder — no TCP
// listener, no port. It records the last status for LastStatus().
type inProcessDoer struct {
	server *api.Server
	status int
}

func (d *inProcessDoer) Do(r *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	d.server.ServeHTTP(rec, r)
	d.status = rec.Code
	return rec.Result(), nil
}
