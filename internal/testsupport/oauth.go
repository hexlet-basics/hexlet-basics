package testsupport

import (
	"crypto/sha256"
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
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/config"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/ids"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/session"
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
	// Router is the production router; mounted OAuth endpoints and the generated
	// API (the consent operations) are both reachable through it.
	Router http.Handler

	now time.Time
}

// NewOAuth builds the stack. The clock starts at the real time and moves only
// through Advance, so a token can be expired without sleeping.
func NewOAuth(t *testing.T) *OAuth {
	t.Helper()
	db, tx, transactor := NewClientWithTxAndTransactor(t)
	o := &OAuth{t: t, DB: db, now: time.Now()}

	translator := NewTranslator(t)
	errorHandler := NewAPIErrorHandler(t, translator)
	services := NewServices(t, tx, transactor)
	mocks := NewMocks()

	server, err := oauthserver.New(oauthserver.Deps{
		Config:     oauthserver.Config{Issuer: OAuthIssuer, EncryptionKey: oauthEncryptionKey},
		DB:         db,
		Transactor: transactor,
		App:        handlers.OAuthSession{},
		Now:        func() time.Time { return o.now },
	})
	require.NoError(t, err)
	o.Server = server

	sessions := session.New(db, session.WithErrorFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		errorHandler.Write(r.Context(), w, r, err)
	}))
	handler := handlers.NewServer(handlers.Deps{
		DB:        db,
		Config:    &config.Config{GuestCookieSecret: testConfig.GuestCookieSecret, EmailTokenSecret: testConfig.EmailTokenSecret},
		Sessions:  sessions,
		Starter:   services.Starter,
		Jobs:      services.Jobs,
		Progress:  progress.New(db, transactor, services.Events, mocks.Runner),
		Registrar: accounts.NewRegistrar(transactor, services.Events),
		Remover:   accounts.NewRemover(transactor),
		Events:    services.Events,
		I18n:      translator,
		Errors:    errorHandler,
		OAuth:     server,
	})
	apiServer, err := api.NewServer(
		handler,
		handler.Security(),
		api.WithErrorHandler(errorHandler.Write),
		api.WithNotFound(handlers.NewNotFoundHandler(translator)),
		api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(translator)),
	)
	require.NoError(t, err)
	// The attachment and webhook handlers are not under test; their routes are
	// registered but never requested.
	o.Router = translator.Middleware(handlers.NewRouter(handler.Security().Pipeline(apiServer), nil, nil, server.Handler()))
	return o
}

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
	_, err := o.DB.OAuthClient.Delete().Exec(ctx)
	require.NoError(o.t, err)
	row, err := o.DB.OAuthClient.Create().
		SetClientID(OAuthClientID).
		SetName("Test MCP client").
		SetRedirectUris(redirectURIs).
		SetScopes(scopes).
		Save(ctx)
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
	value, err := session.New(db).Issue(t.Context(), u)
	require.NoError(t, err)
	return Session{Cookie: &http.Cookie{Name: session.CookieName, Value: value}}
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
	// Params edits the authorization request's query before it is sent (to add a
	// resource, change the redirect_uri, or drop a parameter).
	Params func(url.Values)
}

// StartWith is Start with the registration and the request open to change.
func (o *OAuth) StartWith(userID int, opts StartOptions) *Flow {
	o.t.Helper()
	scopes := opts.Scopes
	if opts.RedirectURIs == nil {
		opts.RedirectURIs = []string{OAuthRedirectURI}
	}
	o.SeedClientURIs(opts.RedirectURIs, scopes...)
	f := &Flow{o: o, UserID: userID, State: "state-" + ids.New(), session: SessionFor(o.t, o.DB, userID)}
	query := url.Values{
		"client_id":             {OAuthClientID},
		"redirect_uri":          {OAuthRedirectURI},
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
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	IDToken      string   `json:"id_token"`
	TokenType    string   `json:"token_type"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        string   `json:"scope"`
	Scopes       []string `json:"-"`
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
		"client_id":     {OAuthClientID},
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
