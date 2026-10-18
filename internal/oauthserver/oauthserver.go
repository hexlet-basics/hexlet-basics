// Package oauthserver is the in-process OAuth 2.1 authorization server that lets
// remote MCP clients act for a signed-in user (ADR-0022).
//
// It is self-contained on purpose, so it can be lifted out of the application
// later: it owns its tables and its configuration, and it must not import the
// HTTP handlers package (a test enforces that). The application gives it one
// thing, a SignedIn function: who is signed in, and how to name them on the consent
// screen. In return the application's handlers call Pending, Approve and Deny
// to serve the consent screen, which is a contract operation like every other
// thing our own frontend calls (ADR-0011); the protocol endpoints an RFC defines
// are served by Handler, outside the contract.
//
// The protocol itself is the zitadel/oidc provider library; this package is the
// storage behind it, the policy around it (PKCE S256 only, public clients only)
// and the consent seam.
package oauthserver

import (
	"context"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"net/http"
	"net/url"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/rs/cors"
	"github.com/samber/oops"
	"github.com/zitadel/oidc/v3/pkg/op"

	"hexletbasics/ent"
	"hexletbasics/internal/store"
)

// errs is the package template for errors (ADR-0018): static messages, the data
// in With, wrapped once at the function that knows the operation.
var errs = oops.In("oauthserver")

const (
	// Endpoint paths. Every authorization-server endpoint sits under /oauth/ with
	// the library defaults overridden, because the root belongs to the SSR server
	// (ADR-0022). The consent screen is the SSR route /connect, outside the prefix.
	authorizePath     = "/oauth/authorize"
	callbackPath      = "/oauth/authorize/callback"
	tokenPath         = "/oauth/token" //nolint:gosec // G101: a URL path named token, not a credential
	revokePath        = "/oauth/revoke"
	userinfoPath      = "/oauth/userinfo"
	endSessionPath    = "/oauth/end_session"
	keysPath          = "/oauth/keys"
	introspectPath    = "/oauth/introspect"
	deviceAuthzPath   = "/oauth/device_authorization"
	consentPagePath   = "/connect"
	metadataPath      = "/.well-known/oauth-authorization-server"
	openIDMetadataURI = "/.well-known/openid-configuration"

	authRequestLifetime = 10 * time.Minute
	// Defaults of Config.AccessTokenLifetime and Config.RefreshTokenLifetime.
	defaultAccessTokenLifetime  = time.Hour
	defaultRefreshTokenLifetime = 30 * 24 * time.Hour
)

// Scopes a user can consent to (ADR-0022). They record what was consented to and
// can only narrow what the user's Permissions allow, never widen it.
const (
	ScopeCatalogRead  = "catalog:read"
	ScopeProgressRead = "progress:read"
	ScopeAdmin        = "admin"
)

// Account is the signed-in user as the consent screen shows them.
type Account struct {
	ID int
	// Name is what the screen calls the user (a name, or the email).
	Name string
}

// Config is everything the package needs from configuration.
type Config struct {
	// Issuer is the bare canonical origin from configuration, never the Host header.
	Issuer string
	// EncryptionKey is the secret access tokens are encrypted under and the ID
	// token signing key is derived from.
	EncryptionKey string
	// RegistrationLimit is how many registration attempts one client address may
	// make per RegistrationWindow; zero values mean the package defaults.
	RegistrationLimit  int
	RegistrationWindow time.Duration
	// AccessTokenLifetime and RefreshTokenLifetime default to an hour and thirty
	// days when zero.
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
}

// Deps are the collaborators of the Server.
type Deps struct {
	Config Config
	// DB reads; Transactor writes that must land together (a token and its refresh
	// token, a rotation).
	DB         *ent.Client
	Transactor store.Transactor
	// SignedIn names the user whose session cookie accompanies ctx. The consent
	// operations are guarded by the session scheme in the contract, so the package
	// only reads it, never checks it again. A function and not an interface: the
	// package may not import the application that implements it (ADR-0025).
	SignedIn func(ctx context.Context) (Account, bool)
	// ClientAddress names the caller of a request, for the registration rate
	// limit. The application supplies it from its client-IP middleware, so the
	// address is never read from a header here; without it every caller shares
	// one bucket.
	ClientAddress func(ctx context.Context) (string, bool)
	// Now is the clock; nil means time.Now. Tests move it to expire things.
	Now func() time.Time
}

// Server is the authorization server.
type Server struct {
	issuer   string
	resource string
	signedIn func(ctx context.Context) (Account, bool)
	limiter  *limiter
	// documents resolves client ids that are metadata document URLs; its own
	// limiter keeps fetches from spending the registration budget and the other
	// way round.
	documents      *documentClients
	metadataClient *http.Client
	// clientAddress is Deps.ClientAddress, never nil.
	clientAddress func(context.Context) (string, bool)
	storage       *storage
	provider      *op.Provider
	handler       http.Handler
}

// errEncryptionKeyRequired is what New answers when the key that seals the
// provider's tokens is missing.
var errEncryptionKeyRequired = errors.New("oauthserver: encryption key is required")

// New builds the Server. It fails on a configuration the provider rejects (an
// issuer that is not an origin, a missing key) so a misconfigured deployment
// refuses to start rather than answering with a broken metadata document.
func New(deps Deps, opts ...Option) (*Server, error) {
	if deps.Config.EncryptionKey == "" {
		return nil, errEncryptionKeyRequired
	}
	if deps.Now == nil {
		deps.Now = time.Now
	}
	if deps.Config.AccessTokenLifetime <= 0 {
		deps.Config.AccessTokenLifetime = defaultAccessTokenLifetime
	}
	if deps.Config.RefreshTokenLifetime <= 0 {
		deps.Config.RefreshTokenLifetime = defaultRefreshTokenLifetime
	}
	issuer, err := url.Parse(deps.Config.Issuer)
	if err != nil {
		return nil, errs.Wrapf(err, "parse issuer")
	}

	key := sha256.Sum256([]byte(deps.Config.EncryptionKey))
	signing := deriveSigningKey(deps.Config.EncryptionKey)
	st := &storage{
		db:       deps.DB,
		tx:       deps.Transactor,
		issuer:   deps.Config.Issuer,
		resource: deps.Config.Issuer + mcpPath,
		now:      deps.Now,
		signing:  signing,

		accessLifetime:  deps.Config.AccessTokenLifetime,
		refreshLifetime: deps.Config.RefreshTokenLifetime,
	}

	providerOpts := []op.Option{
		op.WithCustomEndpoints(
			op.NewEndpoint(authorizePath),
			op.NewEndpoint(tokenPath),
			op.NewEndpoint(userinfoPath),
			op.NewEndpoint(revokePath),
			op.NewEndpoint(endSessionPath),
			op.NewEndpoint(keysPath),
		),
		// The ID tokens are signed with EdDSA, and the library's verifier of an
		// id_token_hint accepts only RS256 unless told otherwise: without this
		// it would reject the server's own tokens.
		op.WithIDTokenHintVerifierOpts(op.WithSupportedIDTokenHintSigningAlgorithms(string(jose.EdDSA))),
		op.WithCustomIntrospectionEndpoint(op.NewEndpoint(introspectPath)),
		op.WithCustomDeviceAuthorizationEndpoint(op.NewEndpoint(deviceAuthzPath)),
		// Browser-based MCP clients call the token and metadata endpoints across
		// origins, but nothing here is authenticated by a cookie, so credentials
		// are never allowed: the library default allows them from any origin.
		op.WithCORSOptions(&cors.Options{
			AllowOriginFunc: func(string) bool { return true },
			AllowedMethods:  []string{http.MethodGet, http.MethodHead, http.MethodPost},
			AllowedHeaders:  []string{"Accept", "Authorization", "Content-Type"},
		}),
	}
	// http is accepted for the issuer only where the library itself would insist
	// on https: a development origin. Production configuration rejects it.
	if issuer.Scheme == "http" {
		providerOpts = append(providerOpts, op.WithAllowInsecure())
	}
	provider, err := op.NewProvider(&op.Config{
		CryptoKey:             key,
		CodeMethodS256:        true,
		GrantTypeRefreshToken: true,
		SupportedScopes: []string{
			ScopeCatalogRead, ScopeProgressRead, ScopeAdmin, "offline_access",
		},
	}, st, op.StaticIssuer(deps.Config.Issuer), providerOpts...)
	if err != nil {
		return nil, errs.Wrapf(err, "build provider")
	}

	clientAddress := deps.ClientAddress
	if clientAddress == nil {
		clientAddress = func(context.Context) (string, bool) { return "", false }
	}
	s := &Server{
		issuer: deps.Config.Issuer, resource: st.resource, signedIn: deps.SignedIn, storage: st, provider: provider,
		limiter:       newLimiter(deps.Config.RegistrationLimit, deps.Config.RegistrationWindow, deps.Now),
		clientAddress: clientAddress,
	}
	for _, opt := range opts {
		opt(s)
	}
	s.documents = newDocumentClients(s.metadataClient, newLimiter(deps.Config.RegistrationLimit, deps.Config.RegistrationWindow, deps.Now), deps.DB, deps.Now)
	s.handler = s.routes()
	return s, nil
}

// Handler serves every public path of the server: /oauth/* and the two
// well-known metadata documents. The router mounts it at those patterns only,
// so the library's own root-level routes (/healthz, /.well-known/...) are never
// reachable.
func (s *Server) Handler() http.Handler { return s.handler }

// Patterns are the ServeMux patterns the router mounts Handler at.
func Patterns() []string {
	return []string{"/oauth/", "GET " + metadataPath, "GET " + openIDMetadataURI}
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+metadataPath, s.metadata)
	mux.HandleFunc("GET "+openIDMetadataURI, s.metadata)
	// Introspection authenticates the caller as a confidential client, and every
	// client here is public; the resource server verifies tokens by lookup. The
	// library mounts it regardless, so it is closed here.
	mux.Handle(introspectPath, http.NotFoundHandler())
	// The device grant is not offered either: its storage is not implemented.
	mux.Handle(deviceAuthzPath, http.NotFoundHandler())
	// Authorization and token are the library's, behind the package's admission
	// policy (scopes, resource) and the iss parameter (policy.go).
	mux.HandleFunc(authorizePath, s.authorize)
	mux.HandleFunc(callbackPath, s.authorizeCallback)
	mux.HandleFunc(tokenPath, s.token)
	// Dynamic Client Registration is ours entirely (registration.go): the library
	// has none.
	mux.Handle(registerPath, s.registerHandler())
	mux.Handle("/oauth/", s.provider)
	return mux
}

// deriveSigningKey derives the ID-token signing key from the encryption secret,
// so every replica signs and verifies with the same key without a second secret
// to store. ID tokens are not used for authorization by anything here (MCP
// clients ignore them), but the library issues them and publishes the key.
func deriveSigningKey(secret string) ed25519.PrivateKey {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte("hexlet-basics oauth id-token signing key v1"))
	return ed25519.NewKeyFromSeed(mac.Sum(nil))
}
