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
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gocloud.dev/blob/memblob"

	"hexletbasics/ent"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/banners"
	"hexletbasics/internal/blogposts"
	"hexletbasics/internal/books"
	"hexletbasics/internal/config"
	"hexletbasics/internal/coursecategories"
	"hexletbasics/internal/courselessons"
	"hexletbasics/internal/courses"
	"hexletbasics/internal/dashboard"
	"hexletbasics/internal/emailtokens"
	"hexletbasics/internal/feeds"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/landingpages"
	"hexletbasics/internal/leads"
	"hexletbasics/internal/lessonreviews"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/qnaitems"
	"hexletbasics/internal/relatedcourses"
	"hexletbasics/internal/reviews"
	"hexletbasics/internal/sessions"
	"hexletbasics/internal/sitemap"
	"hexletbasics/internal/staffroles"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

// testConfig gives handlers fixed public hosts so URL-building assertions are
// deterministic (independent of the ambient env a `config.Load` would read).
var testConfig = &config.Config{
	AppHost:           "code-basics.com",
	PublicURL:         "http://localhost:3001",
	GuestCookieSecret: "test-secret",
	EmailTokenSecret:  "test-email-secret",
	CourseRepoBaseURL: "https://github.com/hexlet-basics",
	SiteURL:           "https://code-basics.com",
	BookBlobKey:       "book.pdf",
}

// EmailTokens signs Magic Link and Password Reset tokens the handlers under
// test accept, so a test can follow a link without a mailbox. Options such as
// emailtokens.WithClock issue one that has already expired.
func EmailTokens(opts ...emailtokens.Option) *emailtokens.Tokens {
	return emailtokens.New(testConfig.EmailTokenSecret, opts...)
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

// NewClientWithTxAndTransactor returns the rolled-back client, the raw
// transaction behind it (to read the outbox and the job queue) and the
// savepoint transactor, for tests that wire a few real services by hand.
func NewClientWithTxAndTransactor(t *testing.T) (*ent.Client, *sql.Tx, store.Transactor) {
	t.Helper()
	client, tx := NewClientWithTx(t)
	return client, tx, newSavepointTransactor(tx, client)
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
		slog.New(slog.DiscardHandler),
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
	// Services are the real event publisher and job enqueuers over that
	// transaction (PublishedEvents and EnqueuedJobs read what they left behind);
	// Mocks are the LLM and Docker boundaries. Both are embedded, so h.Starter,
	// h.Runner and h.Assistant read directly.
	*Services
	*Mocks
	// UserID is the fixture user the harness authenticates as, so a test can
	// assert on rows belonging to the caller without hard-coding a fixture id.
	UserID int
	// Tx is the transaction behind DB and the handlers, for the test that makes
	// the storage fail on purpose (renaming a table inside it) to see how an
	// operation answers a database error. DDL is transactional, so the rollback
	// at cleanup undoes it.
	Tx *sql.Tx
	// Transactor is the savepoint transactor the handlers' modules use, for a
	// test that drives a piece of the operation layer (the audit middleware)
	// directly instead of through HTTP.
	Transactor store.Transactor
	doer       *inProcessDoer
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

// NewHarness builds the in-process test stack, authenticated as the fixture user
// UserAdminLearner: an administrator who is enrolled in JavaScript and has
// finished its first lesson, so one caller reaches both the admin surface and
// the learner's own. The stack is: an ent client over a fresh SQL
// transaction, the handlers.Server, the ogen api.Server with the production
// ErrorHandler, and a generated client whose transport dispatches straight into
// that server's ServeHTTP — so tests run the full request pipeline without a
// socket, and exercise the exact ent-error -> status mapping the server uses.
func NewHarness(t *testing.T) *Harness {
	t.Helper()
	return NewHarnessAs(t, fixtures.UserAdminLearner)
}

// Today is the day the harness's clock reads: the day the fixtures were written,
// so a question the fixtures say was asked today is asked today whenever a test
// runs. The assistant's daily quota counts from the start of it.
var Today = time.Date(2026, time.July, 27, 12, 0, 0, 0, time.UTC)

// NewHarnessAs is NewHarness for a request authenticated as the fixture user
// with the label, for a scenario that belongs to a user of its own.
func NewHarnessAs(t *testing.T, userLabel fixtures.Label) *Harness {
	t.Helper()

	db, tx := NewClientWithTx(t)
	transactor := newSavepointTransactor(tx, db)
	services := NewServices(t, tx, transactor)
	mocks := NewMocks()
	// The real registrar: it hashes, writes the user and raises the sign-up fact
	// through the transactor, so the fact lands in the outbox like in production.
	registrar := accounts.NewRegistrar(transactor, services.Events)
	translator := NewTranslator(t)
	errorHandler := NewAPIErrorHandler(t, translator)
	// The real session manager over the test's transaction: a test's session rows
	// are rolled back with everything else.
	manager := sessions.New(db, sessions.WithErrorFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		errorHandler.Write(r.Context(), w, r, err)
	}))
	// The real progress module, over the test's transaction: the
	// transitions under test are the production ones, only the event transport
	// is recorded instead of written to the outbox and the submission is not
	// really run — running it means Docker, which is the one thing this seam
	// exists to keep out of the tests.
	tracker := progress.New(db, transactor, services.Events, mocks.Runner)
	// An in-memory bucket keeps uploads off disk and out of the network while
	// the attachment metadata still lands in the test's transaction.
	bucket := memblob.OpenBucket(nil)
	t.Cleanup(func() { _ = bucket.Close() })
	assets := assetstore.New(db, bucket, testConfig.PublicURL)
	// The real lead recorder too: a lead test asserts the stored row and the
	// published fact, and both go through the test's transaction.
	leadRecorder := leads.NewRecorder(transactor, services.Events)
	// The real chat over the real progress module: only the model is mocked.
	handler := handlers.NewServer(handlers.Deps{
		DB:      db,
		Config:  testConfig,
		Starter: services.Starter,
		Jobs:    services.Jobs,
		// The real replacer over the savepoint transactor, so a failed replace
		// rolls back exactly as it does in production.
		RelatedCoursesSet: relatedcourses.NewReplacer(transactor),
		Progress:          tracker,
		Assets:            assets,
		Sessions:          manager,
		Registrar:         registrar,
		// The real remover, over the savepoint transactor: what a test asserts
		// about a removed account is what production does to one.
		Remover:       accounts.NewRemover(transactor),
		UserEditor:    accounts.NewEditor(db),
		UserDirectory: accounts.NewDirectory(db),
		Events:        services.Events,
		Leads:         leadRecorder,
		Books:         books.NewRecorder(transactor, services.Events),
		I18n:          translator,
		Errors:        errorHandler,
		YandexFeed:    feeds.NewYandex(db, testConfig.AppHost),
		Assistant: assistant.NewChat(db, transactor, tracker, mocks.Assistant, translator,
			assistant.WithClock(func() time.Time { return Today })),
		// The real remover over the savepoint transactor, like the others.
		PostRemover:       blogposts.NewRemover(transactor),
		Permissions:       staffroles.NewPermissions(db, transactor),
		Dashboard:         dashboard.New(db, landingpages.NewCatalog(db), tracker, slog.New(slog.DiscardHandler)),
		Sitemap:           sitemap.New(db),
		LessonReviews:     lessonreviews.NewScheduler(db, services.Jobs),
		RelatedSuggester:  relatedcourses.NewScheduler(db, services.Jobs),
		LandingPageQna:    qnaitems.NewLandingPage(db),
		CategoryQna:       qnaitems.NewCategory(db),
		PostEditor:        blogposts.NewEditor(db),
		ReviewDirectory:   reviews.NewDirectory(db),
		CourseDirectory:   courses.NewDirectory(db),
		LeadDirectory:     leads.NewDirectory(db),
		MessageDirectory:  assistant.NewDirectory(db),
		CategoryDirectory: coursecategories.NewDirectory(db),
		BannerDirectory:   banners.NewDirectory(db),
		LessonDirectory:   courselessons.NewDirectory(db),
	})
	srv, err := api.NewServer(
		handler,
		handler.Security(),
		api.WithErrorHandler(errorHandler.Write),
		// The production unit of work, over the savepoint transactor: the audit
		// record and the change share a savepoint, as they share a commit.
		api.WithMiddleware(handlers.NewAuditMiddleware(transactor, translator, errorHandler)),
		api.WithNotFound(handlers.NewNotFoundHandler(translator)),
		api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(translator)),
	)
	if err != nil {
		t.Fatalf("new api server: %v", err)
	}

	security := newHarnessSecurity(t, db, manager, userLabel)
	doer := &inProcessDoer{
		// The production pipeline around the generated server, so a test passes
		// the same gates a request does.
		server:  translator.Middleware(handler.Security().Pipeline(srv)),
		session: security.session,
	}
	client, err := api.NewClient("http://test", security, api.WithClient(doer))
	if err != nil {
		t.Fatalf("new api client: %v", err)
	}

	return &Harness{
		Client: &Client{Client: client}, DB: db, Tx: tx, doer: doer,
		Services: services, Mocks: mocks,
		Transactor: transactor,
		UserID:     security.userID,
	}
}

// Mocks are the two external boundaries the harness cannot run for real, preset
// to what most tests want: a passing exercise runner and a short canned chat
// reply. A test that needs another answer sets the mock's Func; it never asserts
// on what the boundary was asked.
type Mocks struct {
	// Assistant is the model behind the in-lesson chat, streaming canned deltas
	// instead of calling the provider.
	Assistant *StreamerMock
	// Runner is the exercise runner behind the check, returning a passing
	// outcome instead of starting a container.
	Runner *ExerciseRunnerMock
}

// NewMocks builds the boundary mocks with their defaults.
func NewMocks() *Mocks {
	return &Mocks{
		Assistant: &StreamerMock{
			StreamFunc: func(_ context.Context, _ []assistant.Turn, onDelta func(string) error) (assistant.Usage, error) {
				for _, delta := range []string{"Try ", "console.log."} {
					if err := onDelta(delta); err != nil {
						return assistant.Usage{}, err
					}
				}
				return assistant.Usage{InputTokens: 100, OutputTokens: 5}, nil
			},
		},
		Runner: &ExerciseRunnerMock{
			RunFunc: func(context.Context, progress.Submission) (progress.Outcome, error) {
				return PassingOutcome(), nil
			},
		},
	}
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

// NewGuestHarness is the harness with a signed guest-progress cookie on every
// request, for the flows that turn a visitor's cookie into account rows.
func NewGuestHarness(t *testing.T, guest progress.GuestProgress) *Harness {
	t.Helper()
	h := NewHarness(t)
	value, err := progress.NewGuestCodec(testConfig.GuestCookieSecret).Encode(guest)
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

// Identify makes every following request carry the given request id and user
// agent, the way a proxy and a browser would, so a test can read them back off
// an Audit Record.
func Identify(h *Harness, requestID, agent string) {
	h.doer.requestID = requestID
	h.doer.agent = agent
}

// ArriveFrom makes every following request come through the ingress on
// behalf of a client at ip, the way production sees a browser.
func ArriveFrom(h *Harness, ip string) {
	h.doer.realIP = ip
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
		guest, err := progress.NewGuestCodec(testConfig.GuestCookieSecret).Decode(parsed.Value)
		if err != nil {
			t.Fatalf("guest cookie does not verify: %v", err)
		}
		return guest
	}
	t.Fatalf("no %s cookie in %v", progress.GuestCookieName, setCookies)
	return progress.GuestProgress{}
}

type harnessSecurity struct {
	session string
	userID  int
}

// newHarnessSecurity signs the fixture user in through the real session manager
// and keeps the cookie value that names the session, which is what a browser
// would send on every request.
func newHarnessSecurity(t *testing.T, db *ent.Client, manager *sessions.Manager, userLabel fixtures.Label) *harnessSecurity {
	t.Helper()

	u, err := db.User.Get(t.Context(), fixtures.ID(userLabel))
	if err != nil {
		t.Fatalf("load harness user: %v", err)
	}
	token, err := manager.Issue(t.Context(), u)
	if err != nil {
		t.Fatalf("issue harness session: %v", err)
	}
	return &harnessSecurity{session: token, userID: u.ID}
}

func (s *harnessSecurity) AdminSession(context.Context, api.OperationName) (api.AdminSession, error) {
	return api.AdminSession{APIKey: s.session}, nil
}

func (s *harnessSecurity) UserSession(context.Context, api.OperationName) (api.UserSession, error) {
	return api.UserSession{APIKey: s.session}, nil
}

// inProcessDoer satisfies ogen's http client interface by serving each request
// straight through the server's ServeHTTP against an in-memory recorder — no TCP
// listener, no port. It records the last status for LastStatus().
type inProcessDoer struct {
	server http.Handler
	status int
	// session is the session cookie a browser would send on EVERY request,
	// including the public reads whose contract declares no security — those
	// still answer a signed-in learner differently.
	session string
	// anonymous drops it, so a public read can be exercised exactly as a visitor
	// reaches it.
	anonymous bool
	// guest is the signed guest-progress cookie a visitor would carry.
	guest string
	// locale is the Accept-Language a browser would send.
	locale string
	// realIP is the client address the ingress would forward.
	realIP string
	// requestID and agent are what a client or proxy would send; empty leaves
	// the header to the default.
	requestID, agent string
	// setCookies are the raw Set-Cookie headers of the last response.
	setCookies []string
}

func (d *inProcessDoer) Do(r *http.Request) (*http.Response, error) {
	switch {
	case d.anonymous:
		r.Header.Del("Cookie")
	case d.session != "":
		r.AddCookie(&http.Cookie{Name: sessions.CookieName, Value: d.session})
	}
	// Added after the anonymous branch cleared the header: a visitor carries
	// progress without carrying a session, which is the whole point of it.
	if d.guest != "" {
		r.AddCookie(&http.Cookie{Name: progress.GuestCookieName, Value: d.guest})
	}
	if d.locale != "" {
		r.Header.Set("Accept-Language", d.locale)
	}
	if d.realIP != "" {
		r.Header.Set("X-Real-IP", d.realIP)
	}
	if d.requestID != "" {
		r.Header.Set("X-Request-ID", d.requestID)
	}
	if d.agent != "" {
		r.Header.Set("User-Agent", d.agent)
	}
	// A streamed body (the generated multipart encoder writes through a pipe)
	// has no length on the client side; over a socket net/http sends it chunked
	// and the server sees -1 (unknown), never 0, which it would read as empty.
	if r.Body != nil && r.Body != http.NoBody && r.ContentLength == 0 {
		r.ContentLength = -1
	}
	rec := httptest.NewRecorder()
	d.server.ServeHTTP(rec, r)
	d.status = rec.Code
	result := rec.Result()
	d.setCookies = result.Header.Values("Set-Cookie")
	return result, nil
}
