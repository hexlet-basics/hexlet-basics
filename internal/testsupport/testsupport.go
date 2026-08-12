// Package testsupport provides Rails-style integration-test plumbing: each test
// runs against a real Postgres transaction that is rolled back on cleanup, over
// the migrated and fixture-loaded container owned by testdb. Handlers hit a real
// database; every test's writes are discarded on rollback, so the package's
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
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pkgz/auth/v2/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/riverqueue/river"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"

	"hexletbasics/ent"
	"hexletbasics/ent/user"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/config"
	"hexletbasics/internal/events"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/ids"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

// testConfig gives handlers fixed public hosts so URL-building assertions are
// deterministic (independent of the ambient env a `config.Load` would read).
var testConfig = &config.Config{
	AppHost:           "code-basics.com",
	PublicURL:         "http://localhost:3001",
	JWTSecret:         "test-secret",
	CourseRepoBaseURL: "https://github.com/hexlet-basics",
}

// NewClient opens an ent client bound to a fresh sql.Tx that is rolled back when
// the test finishes, matching Rails' transactional test lifecycle. Every test
// sees the container's fixture baseline and its own writes vanish on rollback.
func NewClient(t *testing.T) *ent.Client {
	t.Helper()
	client, _ := NewClientWithTransactor(t)
	return client
}

// NewClientWithTransactor returns the test's transaction-bound ent client and a
// savepoint-backed transaction adapter. Business modules can therefore exercise
// their production transaction seam without escaping the outer rollback that
// keeps the shared fixture database clean.
func NewClientWithTransactor(t *testing.T) (*ent.Client, store.Transactor) {
	t.Helper()
	client, tx := NewClientWithTx(t)
	return client, newSavepointTransactor(tx, client)
}

// NewClientWithTx returns the same rolled-back client alongside the raw
// transaction behind it, for the rare test that must issue SQL ent cannot
// express — DDL, or a migration file executed verbatim so that what is verified
// is what ships.
func NewClientWithTx(t *testing.T) (*ent.Client, *sql.Tx) {
	t.Helper()
	db, err := store.NewDB(testdb.DatabaseURL())
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	tx, err := db.BeginTx(t.Context(), nil)
	if err != nil {
		t.Fatalf("begin test transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })

	return store.NewTxClient(tx), tx
}

// NewTranslator loads the same embedded backend catalogs as production.
func NewTranslator(t testing.TB) *localization.Translator {
	t.Helper()
	translator, err := localization.New()
	if err != nil {
		t.Fatalf("new translator: %v", err)
	}
	return translator
}

// NewAPIErrorHandler wires the production handler seam with a silent logger.
// Harness requests do not install a Sentry hub, so they cannot send events.
func NewAPIErrorHandler(t testing.TB, translator *localization.Translator) *handlers.APIErrorHandler {
	t.Helper()
	return handlers.NewAPIErrorHandler(
		translator,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
}

// Harness bundles the generated API client, wired to an in-process server, with
// the ent client behind it — both over the same rolled-back SQL transaction, so
// the client's writes are visible to DB assertions and discarded together.
//
// One-transaction caveat: the whole test shares a single transaction with no
// per-request savepoint, so a write that violates a DB constraint aborts it
// (Postgres 25P02). After an intentional conflict you therefore cannot query the
// DB, nor issue a further "then a valid write" call — assert the conflict's status
// code and stop. In production each request runs in its own autocommit, so this is
// a test-harness artifact, not a behavioural one.
type Harness struct {
	// Client issues typed calls (URLs and bodies are generated, never hand-written).
	Client *Client
	// DB queries the same transaction the handlers wrote through, for assertions.
	DB *ent.Client
	// Enqueuer records background jobs the handlers scheduled, so a test can
	// assert an operation both wrote to the DB and enqueued its follow-up work.
	Enqueuer  *RecordingEnqueuer
	Registrar *RecordingRegistrar
	Events    *RecordingEventPublisher
	// Runner is the exercise runner behind the check, returning a canned outcome
	// instead of starting a container.
	Runner *StubExerciseRunner
	// UserID is the fixture user the harness authenticates as, so a test can
	// assert on rows belonging to the caller without hard-coding a fixture id.
	UserID int
	doer   *inProcessDoer
}

// HarnessUser reports the user every harness request is authenticated as.
func HarnessUser(h *Harness) (int, bool) {
	return h.UserID, h.UserID != 0
}

// LastStatus is the HTTP status of the most recent client call. It remains a
// concise assertion helper even though default errors now decode to the
// generated ProblemDetailsStatusCode type.
func (h *Harness) LastStatus() int { return h.doer.status }

// ResponseCookies are the raw Set-Cookie headers of the most recent call. The
// generated client decodes that header as a single comma-separated array and so
// reports only the first cookie; a test that cares which cookies were actually
// set has to look at the response itself.
func (h *Harness) ResponseCookies() []string { return h.doer.setCookies }

// NewHarness builds the in-process test stack: an ent client over a fresh SQL
// transaction, the handlers.Server, the ogen api.Server with the production
// ErrorHandler, and a generated client whose transport dispatches straight into
// that server's ServeHTTP — so tests run the full request pipeline without a
// socket, and exercise the exact ent-error -> status mapping the server uses.
func NewHarness(t *testing.T) *Harness {
	t.Helper()

	db, transactor := NewClientWithTransactor(t)

	enqueuer := &RecordingEnqueuer{DB: db}
	registrar := NewRecordingRegistrar(db)
	eventPublisher := &RecordingEventPublisher{}
	translator := NewTranslator(t)
	errorHandler := NewAPIErrorHandler(t, translator)
	// The real progress module, over the test's transaction: the gate and the
	// transitions under test are the production ones, only the event transport
	// is recorded instead of written to the outbox and the submission is not
	// really run — running it means Docker, which is the one thing this seam
	// exists to keep out of the tests.
	runner := NewStubExerciseRunner()
	tracker := progress.New(db, transactor, eventPublisher, runner)
	handler := handlers.NewServer(db, testConfig, enqueuer, enqueuer, tracker, registrar, eventPublisher, translator, errorHandler)
	srv, err := api.NewServer(
		handler,
		handler.AuthHandler(),
		api.WithErrorHandler(errorHandler.Write),
		api.WithNotFound(handlers.NewNotFoundHandler(translator)),
		api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(translator)),
	)
	if err != nil {
		t.Fatalf("new api server: %v", err)
	}

	security := newHarnessSecurity(t, db)
	doer := &inProcessDoer{
		server: translator.Middleware(handler.AuthHandler().Trace(
			handler.AuthHandler().Identify(handler.AuthHandler().CarryGuestProgress(srv)),
		)),
		jwt: security.jwt,
	}
	client, err := api.NewClient("http://test", security, api.WithClient(doer))
	if err != nil {
		t.Fatalf("new api client: %v", err)
	}

	return &Harness{
		Client: &Client{Client: client}, DB: db, doer: doer,
		Enqueuer: enqueuer, Registrar: registrar, Events: eventPublisher,
		Runner: runner,
		UserID: security.userID,
	}
}

// StubExerciseRunner answers every submission with a canned outcome and records
// what it was asked to run. A passing outcome is the default because it is the
// one that moves progress; a test after the failure path sets Outcome itself.
type StubExerciseRunner struct {
	Outcome     progress.Outcome
	Err         error
	Submissions []progress.Submission
}

// NewStubExerciseRunner builds the stub with a passing outcome.
func NewStubExerciseRunner() *StubExerciseRunner {
	return &StubExerciseRunner{Outcome: PassingOutcome()}
}

// PassingOutcome is what the runner returns for a solution that passes.
func PassingOutcome() progress.Outcome {
	return progress.Outcome{Passed: true, Result: progress.ResultPassed, Output: "", Status: 0}
}

// FailingOutcome is an ordinary test failure: the submission ran and lost.
func FailingOutcome() progress.Outcome {
	return progress.Outcome{
		Passed: false,
		Result: progress.ResultFailed,
		Output: "expected 3, got 2",
		Status: 1,
	}
}

// Run records the submission and returns the canned outcome.
func (r *StubExerciseRunner) Run(
	_ context.Context,
	submission progress.Submission,
) (progress.Outcome, error) {
	r.Submissions = append(r.Submissions, submission)
	if r.Err != nil {
		return progress.Outcome{}, r.Err
	}
	return r.Outcome, nil
}

// HarnessUserPassword is the password GivePassword sets, so a test can sign in
// as a fixture user without inventing credentials of its own.
const HarnessUserPassword = "harness-password"

// GivePassword gives a fixture user a usable password. The fixtures carry no
// digest — the harness normally mints a JWT directly — so a test that exercises
// the password flow has to supply one.
func GivePassword(t *testing.T, h *Harness, userID int) string {
	t.Helper()
	digest, err := bcrypt.GenerateFromPassword([]byte(HarnessUserPassword), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash harness password: %v", err)
	}
	u := h.DB.User.UpdateOneID(userID).
		SetPasswordDigest(string(digest)).
		SaveX(t.Context())
	return lo.FromPtr(u.Email)
}

// NewGuestHarness is the harness with a signed guest-progress cookie on every
// request, for the flows that turn a visitor's cookie into account rows.
func NewGuestHarness(t *testing.T, guest progress.GuestProgress) *Harness {
	t.Helper()
	h := NewHarness(t)
	value, err := progress.NewGuestCodec(testConfig.JWTSecret).Encode(guest)
	if err != nil {
		t.Fatalf("encode guest progress: %v", err)
	}
	h.doer.guest = value
	return h
}

// NewAnonymousHarness is the same in-process stack with no session cookie, for
// the public reads that must answer a visitor as well as a learner.
func NewAnonymousHarness(t *testing.T) *Harness {
	t.Helper()
	h := NewHarness(t)
	h.doer.anonymous = true
	return h
}

// NewVisitorHarness is a visitor with no account and the progress cookie a
// previous visit left them: the only client of the guest storage that is not
// also signing in.
func NewVisitorHarness(t *testing.T, guest progress.GuestProgress) *Harness {
	t.Helper()
	h := NewGuestHarness(t, guest)
	h.doer.anonymous = true
	return h
}

// SpeakTo makes every following request ask for a locale, as a browser does
// with Accept-Language. The generated client has no seam for request headers,
// and the locale decides which translation of a lesson is served.
func SpeakTo(h *Harness, locale string) {
	h.doer.locale = locale
}

// ForgeGuestCookie makes the harness carry progress signed with the wrong
// secret — what a visitor editing their own cookie can produce. The server must
// treat it as no cookie at all.
func ForgeGuestCookie(t *testing.T, h *Harness, guest progress.GuestProgress) {
	t.Helper()
	value, err := progress.NewGuestCodec("not-the-server-secret").Encode(guest)
	if err != nil {
		t.Fatalf("encode forged guest progress: %v", err)
	}
	h.doer.guest = value
}

// DecodeGuestCookie reads the guest progress a response set, verifying its
// signature exactly as the server does on the next request. It fails the test
// when no guest cookie was set, so a caller asserting on progress cannot pass
// by accident on an unchanged one.
func DecodeGuestCookie(t *testing.T, setCookies []string) progress.GuestProgress {
	t.Helper()

	for _, raw := range setCookies {
		parsed, err := http.ParseSetCookie(raw)
		if err != nil || parsed.Name != progress.GuestCookieName {
			continue
		}
		guest, err := progress.NewGuestCodec(testConfig.JWTSecret).Decode(parsed.Value)
		if err != nil {
			t.Fatalf("guest cookie does not verify: %v", err)
		}
		return guest
	}
	t.Fatalf("no %s cookie in %v", progress.GuestCookieName, setCookies)
	return progress.GuestProgress{}
}

type harnessSecurity struct {
	jwt    string
	xsrf   string
	userID int
}

func newHarnessSecurity(t *testing.T, db *ent.Client) *harnessSecurity {
	t.Helper()

	u, err := db.User.Query().Where(user.AdminEQ(true)).First(t.Context())
	if err != nil {
		t.Fatalf("load harness admin: %v", err)
	}

	jti := ids.New()
	rec := httptest.NewRecorder()
	_, err = token.NewService(token.Opts{
		SecretReader: token.SecretFunc(func(string) (string, error) {
			return testConfig.JWTSecret, nil
		}),
	}).Set(rec, token.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:       jti,
			Audience: jwt.ClaimStrings{"hexlet-basics"},
		},
		User: &token.User{
			ID:   fmt.Sprintf("%d", u.ID),
			Name: *u.Email,
		},
		AuthProvider: &token.AuthProvider{Name: "password"},
	})
	if err != nil {
		t.Fatalf("create harness JWT: %v", err)
	}

	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "JWT" {
			return &harnessSecurity{jwt: cookie.Value, xsrf: jti, userID: u.ID}
		}
	}
	t.Fatal("JWT cookie was not issued")
	return nil
}

func (s *harnessSecurity) AdminSession(context.Context, api.OperationName) (api.AdminSession, error) {
	return api.AdminSession{APIKey: s.jwt}, nil
}

func (s *harnessSecurity) UserSession(context.Context, api.OperationName) (api.UserSession, error) {
	return api.UserSession{APIKey: s.jwt}, nil
}

func (s *harnessSecurity) XsrfToken(context.Context, api.OperationName) (api.XsrfToken, error) {
	return api.XsrfToken{APIKey: s.xsrf}, nil
}

// RecordingEventPublisher captures standalone domain facts in handler tests.
type RecordingEventPublisher struct {
	Published []events.Event
	Err       error
}

// PublishStandalone records the fact without touching the SQL outbox.
func (p *RecordingEventPublisher) PublishStandalone(
	_ context.Context,
	event events.Event,
) error {
	if p.Err != nil {
		return p.Err
	}
	p.Published = append(p.Published, event)
	return nil
}

// Publish records a fact raised inside a caller-owned transaction. The harness
// asserts on what was published; Watermill's outbox write is production's job.
func (p *RecordingEventPublisher) Publish(
	_ context.Context,
	_ *sql.Tx,
	event events.Event,
) error {
	if p.Err != nil {
		return p.Err
	}
	p.Published = append(p.Published, event)
	return nil
}

// RecordingRegistrar mirrors account creation inside the harness transaction
// and records the requested registration without touching Watermill.
type RecordingRegistrar struct {
	DB            *ent.Client
	Registrations []accounts.Registration
}

// NewRecordingRegistrar builds the auth test adapter.
func NewRecordingRegistrar(db *ent.Client) *RecordingRegistrar {
	return &RecordingRegistrar{DB: db}
}

// Register persists a bcrypt-compatible user so login tests exercise the real
// credential path after registration.
func (r *RecordingRegistrar) Register(
	ctx context.Context,
	input accounts.Registration,
) (*ent.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, accounts.ErrPasswordProcessing
	}
	u, err := r.DB.User.Create().
		SetEmail(input.Email).
		SetPasswordDigest(string(hash)).
		SetNillableFirstName(input.FirstName).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	r.Registrations = append(r.Registrations, input)
	return u, nil
}

// RecordingEnqueuer is a test adapter for the handlers' job seams
// (VersionBuildStarter, LessonReviewEnqueuer). It performs the visible DB
// writes through the harness's rollback-only ent client and records the job
// args without touching River.
type RecordingEnqueuer struct {
	DB       *ent.Client
	Inserted []river.JobArgs
}

// EnqueueLessonReviews records one review job per lesson info.
func (e *RecordingEnqueuer) EnqueueLessonReviews(_ context.Context, lessonInfoIDs []int) error {
	for _, id := range lessonInfoIDs {
		e.Inserted = append(e.Inserted, jobs.ReviewLessonArgs{LessonInfoID: id})
	}
	return nil
}

// Start mirrors the production operation's visible result.
func (e *RecordingEnqueuer) Start(ctx context.Context, courseID int) (*ent.CourseVersion, error) {
	version, err := e.DB.CourseVersion.Create().
		SetCourseID(courseID).
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
	server http.Handler
	status int
	// jwt is the session a browser would send on EVERY request, including the
	// public reads whose contract declares no security — those still answer a
	// signed-in learner differently.
	jwt string
	// anonymous drops it, so a public read can be exercised exactly as a visitor
	// reaches it.
	anonymous bool
	// guest is the signed guest-progress cookie a visitor would carry.
	guest string
	// locale is the Accept-Language a browser would send.
	locale string
	// setCookies are the raw Set-Cookie headers of the last response.
	setCookies []string
}

func (d *inProcessDoer) Do(r *http.Request) (*http.Response, error) {
	switch {
	case d.anonymous:
		r.Header.Del("Cookie")
		r.Header.Del("X-XSRF-TOKEN")
	case d.jwt != "":
		r.AddCookie(&http.Cookie{Name: "JWT", Value: d.jwt})
	}
	// Added after the anonymous branch cleared the header: a visitor carries
	// progress without carrying a session, which is the whole point of it.
	if d.guest != "" {
		r.AddCookie(&http.Cookie{Name: progress.GuestCookieName, Value: d.guest})
	}
	if d.locale != "" {
		r.Header.Set("Accept-Language", d.locale)
	}
	rec := httptest.NewRecorder()
	d.server.ServeHTTP(rec, r)
	d.status = rec.Code
	result := rec.Result()
	d.setCookies = result.Header.Values("Set-Cookie")
	return result, nil
}
