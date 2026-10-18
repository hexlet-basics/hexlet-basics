package testsupport

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/oauthclient"
	"hexletbasics/internal/api"
	"hexletbasics/internal/auth"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/ident"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport/testdb"
)

const (
	// OAuthIssuer is the canonical origin the test authorization server speaks as.
	OAuthIssuer = "https://code-basics.test"
	// OAuthClientID is the pre-registered public client the flow helper seeds.
	OAuthClientID = "test-mcp-client"
	// OAuthRedirectURI is its registered loopback redirect URI.
	OAuthRedirectURI = "http://127.0.0.1:8976/callback"
	// oauthEncryptionKey is the secret the test server encrypts tokens under.
	oauthEncryptionKey = "test-oauth-encryption-key-test-only-test-only"
	oauthVerifier      = "test-pkce-verifier-with-enough-entropy-0123456789-abcdefghijklmnopqrstuvwxyz"
)

// OAuth is the OAuth authorization server wired the way production wires it,
// over the test's rolled-back transaction: the real router (the server mounted
// beside the generated API), the real handlers behind the consent operations,
// and a controllable clock. Requests go through ServeHTTP in process, no socket.
//
// Connect is what every later OAuth and MCP test starts from: it runs the whole
// authorization-code flow with PKCE for a user and a list of Scopes and hands
// back the tokens, so a test about tools or Connected Apps never re-implements
// the dance. Use Start for a test about the flow itself (denying, expiring,
// replaying), where each step is driven by hand.
type OAuth struct {
	t *testing.T
	// DB queries the transaction everything wrote through.
	DB *ent.Client
	// Server is the authorization server, for VerifyAccessToken.
	Server *oauthserver.Server
	// Services holds the outbox, for PublishedEvents and Deliver (a read of
	// personal data over MCP leaves a fact there, as over HTTP).
	Services *Services
	// Router is the production router; mounted OAuth endpoints and the generated
	// API (the consent operations) are both reachable through it.
	Router http.Handler
	// Documents is the web host of Client ID Metadata Documents: the server
	// fetches the documents it publishes the way it would fetch them from the
	// internet (https, a hostname), but they are served from this process.
	Documents *ClientDocuments

	now time.Time
}

type oauthOptions struct {
	mcpLimits mcpserver.Limits
	onPool    bool
}

// OAuthOption changes how NewOAuth builds the stack.
type OAuthOption func(*oauthOptions)

// WithMCPLimits sets the MCP endpoint's limits (page size, response size,
// timeout, rates); a field left zero keeps its default. Tests of a limit set
// only that one, small, so they do not have to produce a hundred kilobytes or
// wait half a minute.
func WithMCPLimits(limits mcpserver.Limits) OAuthOption {
	return func(o *oauthOptions) { o.mcpLimits = limits }
}

// OnPool builds the stack over a connection pool of its own instead of the
// test's rolled-back transaction, so that every write commits and separate
// requests run in separate transactions. Only a test about concurrent
// transactions wants it; the pool is a clone dropped when the test ends, so
// there is nothing to delete.
func OnPool() OAuthOption {
	return func(o *oauthOptions) { o.onPool = true }
}

// NewOAuth builds the stack. The clock starts at the real time and moves only
// through Advance, so a token can be expired without sleeping.
//
// Client ID Metadata Documents are fetched from o.Documents (see ClientDocuments).
func NewOAuth(t *testing.T, opts ...OAuthOption) *OAuth {
	t.Helper()
	return newOAuth(t, true, oauthserver.Config{}, opts)
}

// NewOAuthWith is NewOAuth with settings of the server open to change: the
// token lifetimes. The issuer and the key are always the test's own.
func NewOAuthWith(t *testing.T, cfg oauthserver.Config, opts ...OAuthOption) *OAuth {
	t.Helper()
	return newOAuth(t, true, cfg, opts)
}

// NewOAuthWithProductionFetcher is NewOAuth with the metadata document client
// the application wires in production, which refuses every non-public address.
// A test about what the server will not fetch uses it; there is no document
// host behind it.
func NewOAuthWithProductionFetcher(t *testing.T, opts ...OAuthOption) *OAuth {
	t.Helper()
	return newOAuth(t, false, oauthserver.Config{}, opts)
}

func newOAuth(t *testing.T, documentHost bool, cfg oauthserver.Config, opts []OAuthOption) *OAuth {
	t.Helper()
	cfg.Issuer, cfg.EncryptionKey = OAuthIssuer, oauthEncryptionKey
	var options oauthOptions
	for _, opt := range opts {
		opt(&options)
	}
	var (
		db         *ent.Client
		tx         *sql.Tx
		transactor store.Transactor
	)
	if options.onPool {
		pool := testdb.Clone(t)
		db, transactor = store.NewClient(pool), store.New(pool)
	} else {
		db, tx, transactor = NewClientWithTxAndTransactor(t)
	}
	o := &OAuth{t: t, DB: db, now: time.Now()}
	if documentHost {
		o.Documents = newClientDocuments(t)
	}

	translator := NewTranslator(t)
	errorHandler := NewAPIErrorHandler(t, translator)
	services := NewServices(t, tx, transactor)
	o.Services = services
	mocks := NewMocks()
	tracker := progress.New(db, transactor, services.Events, mocks.Runner)

	var serverOpts []oauthserver.Option
	if documentHost {
		serverOpts = append(serverOpts, oauthserver.WithMetadataHTTPClient(o.Documents.client))
	}
	server, err := oauthserver.New(oauthserver.Deps{
		Config:        cfg,
		DB:            db,
		Transactor:    transactor,
		App:           handlers.OAuthSession{},
		ClientAddress: handlers.ClientIP,
		Now:           func() time.Time { return o.now },
	}, serverOpts...)
	require.NoError(t, err)
	o.Server = server

	manager := auth.New(db, auth.WithErrorFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		errorHandler.Write(r.Context(), w, r, err)
	}))
	// The harness's collaborators, so every tool reaches the real module a
	// request does (a nil collaborator would panic a tool in a test); only the
	// OAuth server, which the harness does not run, is added.
	deps := handlerDeps(t, db, transactor, services, mocks, tracker, translator, errorHandler, manager)
	deps.OAuth = server
	handler := handlers.NewServer(deps)
	apiServer, err := api.NewServer(
		handler,
		handler.Security(),
		api.WithErrorHandler(errorHandler.Write),
		// The production unit of work, so an admin write over MCP leaves its Audit
		// Record here as it does in the running server.
		// One call: a second WithMiddleware would replace the first.
		api.WithMiddleware(
			handlers.NewAuditMiddleware(transactor, translator, errorHandler),
			handlers.NewPersonalDataReadMiddleware(services.Events, translator, errorHandler),
		),
		api.WithNotFound(handlers.NewNotFoundHandler(errorHandler)),
		api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(errorHandler)),
	)
	require.NoError(t, err)
	// The endpoint's rate limiter runs on the same controllable clock as the
	// OAuth server, unless a test brings its own.
	limits := options.mcpLimits
	if limits.Now == nil {
		limits.Now = func() time.Time { return o.now }
	}
	mcpEndpoint, err := mcpserver.New(mcpserver.Deps{Verifier: server, Catalogue: handler, MyDashboard: handler.MyDashboard, MyCourseProgress: handler.MyCourseProgress, HeldCapabilities: handler.Security().HeldCapabilities, API: handlers.NewAdminDispatcher(apiServer, translator), Issuer: OAuthIssuer, Limits: limits})
	require.NoError(t, err)
	// The attachment and webhook handlers are not under test; their routes are
	// registered but never requested.
	o.Router = translator.Middleware(handlers.NewRouter(handler.Security().Pipeline(apiServer), nil, nil, server.Handler(), mcpEndpoint.Handler()))
	return o
}

// Now is the server's clock.
func (o *OAuth) Now() time.Time { return o.now }

// Advance moves the server's clock forward.
func (o *OAuth) Advance(d time.Duration) { o.now = o.now.Add(d) }

// SeedClient registers the public test client allowed to request the given
// custom scopes (replacing an earlier registration), and returns its row.
func (o *OAuth) SeedClient(scopes ...string) *ent.OAuthClient {
	o.t.Helper()
	return o.SeedClientURIs([]string{OAuthRedirectURI}, scopes...)
}

// SeedClientURIs is SeedClient with the redirect URIs the client registered, for
// a test about which redirects match.
func (o *OAuth) SeedClientURIs(redirectURIs []string, scopes ...string) *ent.OAuthClient {
	o.t.Helper()
	ctx := o.t.Context()
	// Updated in place when it exists: deleting the row would cascade away the
	// tokens of an earlier Connect, and a test about families needs several.
	existing, err := o.DB.OAuthClient.Query().Where(oauthclient.ClientID(OAuthClientID)).Only(ctx)
	if ent.IsNotFound(err) {
		row, err := o.DB.OAuthClient.Create().
			SetClientID(OAuthClientID).
			SetName("Test MCP client").
			SetRedirectUris(redirectURIs).
			SetScopes(scopes).
			Save(ctx)
		require.NoError(o.t, err)
		return row
	}
	require.NoError(o.t, err)
	row, err := existing.Update().SetRedirectUris(redirectURIs).SetScopes(scopes).Save(ctx)
	require.NoError(o.t, err)
	return row
}

// Session is a signed-in browser: the session cookie it sends.
type Session struct {
	Cookie *http.Cookie
}

// SessionFor starts a session for a fixture user through the real session
// manager, as the sign-in operations would have, and returns its cookie.
func SessionFor(t *testing.T, db *ent.Client, userID int) Session {
	t.Helper()
	u, err := db.User.Get(t.Context(), userID)
	require.NoError(t, err)
	value, err := auth.New(db).Issue(t.Context(), u)
	require.NoError(t, err)
	return Session{Cookie: &http.Cookie{Name: auth.CookieName, Value: value}}
}

// Do sends one request through the router and returns the recorded response.
func (o *OAuth) Do(method, target string, form url.Values, mutate ...func(*http.Request)) *http.Response {
	o.t.Helper()
	var body *strings.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	} else {
		body = strings.NewReader("")
	}
	req := httptest.NewRequestWithContext(o.t.Context(), method, target, body)
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for _, m := range mutate {
		m(req)
	}
	rec := httptest.NewRecorder()
	o.Router.ServeHTTP(rec, req)
	return rec.Result()
}

// Flow is one authorization-code flow with PKCE in progress, driven step by step.
type Flow struct {
	o *OAuth
	// UserID is the user who signs in and decides.
	UserID int
	// ClientID is the client the flow runs as.
	ClientID string
	// State is the opaque value the client expects echoed back.
	State string
	// RequestID is the authorization request id the consent screen is opened with
	// (from the authorize redirect).
	RequestID string
	// RedirectURI is the redirect_uri the authorization request named; the token
	// request repeats it.
	RedirectURI string
	// Callback is where Approve was sent back to at the client, with the code
	// and iss it carries.
	Callback *url.URL
	// Authorize is the response of GET /oauth/authorize.
	Authorize *http.Response
	session   Session
}

// Start begins a flow for a user and Scopes: it registers the client for them
// (plus offline_access, which the library requires to issue a refresh token)
// and sends the authorization request. The response, normally the redirect to
// the consent screen, is in Authorize.
func (o *OAuth) Start(userID int, scopes ...string) *Flow {
	o.t.Helper()
	return o.StartWith(userID, StartOptions{Scopes: scopes})
}

// StartOptions shapes the flow of StartWith beyond the defaults of Start.
type StartOptions struct {
	// Scopes are registered for the client and requested.
	Scopes []string
	// RedirectURIs are what the client registers; the default is OAuthRedirectURI.
	RedirectURIs []string
	// ClientID, when set, runs the flow as an already registered client (one that
	// registered itself through /oauth/register): nothing is seeded and Scopes
	// and RedirectURIs are only what the request names.
	ClientID string
	// Params edits the authorization request's query before it is sent (to add a
	// resource, change the redirect_uri, or drop a parameter).
	Params func(url.Values)
}

// StartWith is Start with the registration and the request open to change.
func (o *OAuth) StartWith(userID int, opts StartOptions) *Flow {
	o.t.Helper()
	scopes := opts.Scopes
	if opts.RedirectURIs == nil && opts.ClientID == "" {
		opts.RedirectURIs = []string{OAuthRedirectURI}
	}
	clientID := opts.ClientID
	if clientID == "" {
		clientID = OAuthClientID
		o.SeedClientURIs(opts.RedirectURIs, scopes...)
	}
	redirectURI := OAuthRedirectURI
	if opts.ClientID != "" {
		redirectURI = opts.RedirectURIs[0]
	}
	f := &Flow{o: o, UserID: userID, ClientID: clientID, State: "state-" + ident.New(), session: SessionFor(o.t, o.DB, userID)}
	query := url.Values{
		"client_id":             {clientID},
		"redirect_uri":          {redirectURI},
		"response_type":         {"code"},
		"scope":                 {strings.Join(append(append([]string{}, scopes...), "offline_access"), " ")},
		"state":                 {f.State},
		"code_challenge":        {PKCEChallenge()},
		"code_challenge_method": {"S256"},
	}
	if opts.Params != nil {
		opts.Params(query)
	}
	f.RedirectURI = query.Get("redirect_uri")
	f.Authorize = o.Do(http.MethodGet, "/oauth/authorize?"+query.Encode(), nil)
	if loc, err := f.Authorize.Location(); err == nil {
		f.RequestID = loc.Query().Get("id")
	}
	return f
}

// PKCEChallenge is the S256 challenge of the verifier every flow uses.
func PKCEChallenge() string {
	sum := sha256.Sum256([]byte(oauthVerifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// Decide approves or denies as the signed-in user, through the contract
// operations the consent screen calls, and returns where the browser is sent.
func (f *Flow) Decide(approve bool) *http.Response {
	f.o.t.Helper()
	action := "deny"
	if approve {
		action = "approve"
	}
	return f.o.Do(http.MethodPost, "/api/oauth/consents/"+url.PathEscape(f.RequestID)+"/"+action, nil,
		func(r *http.Request) {
			r.AddCookie(f.session.Cookie)
		})
}

// Approve decides yes, follows the callback, and returns the authorization code
// the client receives at its redirect URI.
func (f *Flow) Approve() string {
	f.o.t.Helper()
	res := f.Decide(true)
	require.Equal(f.o.t, http.StatusOK, res.StatusCode)
	var decision struct {
		RedirectURL string `json:"redirectUrl"`
	}
	require.NoError(f.o.t, json.NewDecoder(res.Body).Decode(&decision))
	callback, err := url.Parse(decision.RedirectURL)
	require.NoError(f.o.t, err)

	res = f.o.Do(http.MethodGet, callback.RequestURI(), nil)
	require.Equal(f.o.t, http.StatusFound, res.StatusCode)
	back, err := res.Location()
	require.NoError(f.o.t, err)
	f.Callback = back
	require.Equal(f.o.t, f.State, back.Query().Get("state"))
	code := back.Query().Get("code")
	require.NotEmpty(f.o.t, code, "redirect carried no code: %s", back)
	return code
}

// Tokens is a token response.
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	// ExpiresIn is unsigned, as the library writes it: it counts against the
	// real clock, so a token issued under a shifted test clock is already past
	// and the negative lifetime wraps to a huge uint64 on amd64 (CI) and
	// saturates to 0 on arm64. An int would fail to decode the wrapped value.
	ExpiresIn uint64   `json:"expires_in"`
	Scope     string   `json:"scope"`
	Scopes    []string `json:"-"`
}

// Exchange trades an authorization code for tokens at /oauth/token and returns
// the raw response too, so a test can assert on a refusal.
func (f *Flow) Exchange(code, verifier string) (Tokens, *http.Response) {
	f.o.t.Helper()
	return f.ExchangeWith(code, verifier, nil)
}

// ExchangeWith is Exchange with extra token request parameters (a resource).
func (f *Flow) ExchangeWith(code, verifier string, extra url.Values) (Tokens, *http.Response) {
	f.o.t.Helper()
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {f.ClientID},
		"redirect_uri":  {f.RedirectURI},
		"code":          {code},
		"code_verifier": {verifier},
	}
	for k, v := range extra {
		form[k] = v
	}
	res := f.o.Do(http.MethodPost, "/oauth/token", form)
	var tokens Tokens
	if res.StatusCode == http.StatusOK {
		require.NoError(f.o.t, json.NewDecoder(res.Body).Decode(&tokens))
		tokens.Scopes = strings.Fields(tokens.Scope)
	}
	return tokens, res
}

// PKCEVerifier is the verifier matching PKCEChallenge.
func PKCEVerifier() string { return oauthVerifier }

// Connect completes the whole flow for a user and Scopes on the test
// transaction, the way a person connecting an assistant would: authorize, sign
// in and approve on the consent screen, exchange the code. It fails the test if
// any step does not behave, and returns tokens that VerifyAccessToken accepts.
// The client is registered for the Scopes; offline_access is added so a refresh
// token comes back.
func (o *OAuth) Connect(userID int, scopes ...string) Tokens {
	o.t.Helper()
	f := o.Start(userID, scopes...)
	require.Equal(o.t, http.StatusFound, f.Authorize.StatusCode)
	require.NotEmpty(o.t, f.RequestID)
	code := f.Approve()
	tokens, res := f.Exchange(code, PKCEVerifier())
	require.Equal(o.t, http.StatusOK, res.StatusCode)
	return tokens
}

// ConnectAs is Connect through a named Connected App, for a test with several
// users or apps connected at once. Connect's registration replaces every client,
// and deleting a client revokes its tokens, so a second Connect would end the
// first one's session; ConnectAs registers clientID only when it is not
// registered yet and leaves the others alone. Use OAuthClientID to connect
// another user through the app Connect registered (call Connect first), and any
// other id for a second app. A test of what is counted per client needs two apps
// whose users overlap.
func (o *OAuth) ConnectAs(clientID string, userID int, scopes ...string) Tokens {
	o.t.Helper()
	known, err := o.DB.OAuthClient.Query().Where(oauthclient.ClientID(clientID)).Exist(o.t.Context())
	require.NoError(o.t, err)
	if !known {
		_, err = o.DB.OAuthClient.Create().
			SetClientID(clientID).
			SetName("Test MCP client " + clientID).
			SetRedirectUris([]string{OAuthRedirectURI}).
			SetScopes(scopes).
			Save(o.t.Context())
		require.NoError(o.t, err)
	}
	f := o.StartWith(userID, StartOptions{Scopes: scopes, ClientID: clientID, RedirectURIs: []string{OAuthRedirectURI}})
	require.Equal(o.t, http.StatusFound, f.Authorize.StatusCode)
	tokens, res := f.Exchange(f.Approve(), PKCEVerifier())
	require.Equal(o.t, http.StatusOK, res.StatusCode)
	return tokens
}

// Register posts a Dynamic Client Registration request (RFC 7591) as an
// anonymous caller. The body is sent as it is when it is a string or []byte,
// and as JSON otherwise; mutate edits the request (the client address header,
// the content type).
func (o *OAuth) Register(body any, mutate ...func(*http.Request)) *http.Response {
	o.t.Helper()
	var raw []byte
	switch b := body.(type) {
	case string:
		raw = []byte(b)
	case []byte:
		raw = b
	default:
		var err error
		raw, err = json.Marshal(body)
		require.NoError(o.t, err)
	}
	req := httptest.NewRequestWithContext(o.t.Context(), http.MethodPost, "/oauth/register", strings.NewReader(string(raw)))
	req.Header.Set("Content-Type", "application/json")
	for _, m := range mutate {
		m(req)
	}
	rec := httptest.NewRecorder()
	o.Router.ServeHTTP(rec, req)
	return rec.Result()
}

// FromAddress makes a request come from a client address, as the ingress would
// forward it.
func FromAddress(ip string) func(*http.Request) {
	return func(r *http.Request) { r.Header.Set("X-Real-IP", ip) }
}

// Refresh presents a refresh token at /oauth/token as a client and returns the
// new tokens (zero when refused) with the raw response.
func (o *OAuth) Refresh(clientID, refreshToken string) (Tokens, *http.Response) {
	o.t.Helper()
	res := o.Do(http.MethodPost, "/oauth/token", url.Values{
		"grant_type": {"refresh_token"}, "client_id": {clientID}, "refresh_token": {refreshToken},
	})
	var tokens Tokens
	if res.StatusCode == http.StatusOK {
		require.NoError(o.t, json.NewDecoder(res.Body).Decode(&tokens))
		tokens.Scopes = strings.Fields(tokens.Scope)
	}
	return tokens, res
}

// Revoke posts a token to /oauth/revoke (RFC 7009) as a client.
func (o *OAuth) Revoke(clientID, token string) *http.Response {
	o.t.Helper()
	return o.Do(http.MethodPost, "/oauth/revoke", url.Values{"client_id": {clientID}, "token": {token}})
}

// Accepts reports whether the server still accepts an access token.
func (o *OAuth) Accepts(accessToken string) bool {
	o.t.Helper()
	_, err := o.Server.VerifyAccessToken(o.t.Context(), accessToken)
	return err == nil
}
