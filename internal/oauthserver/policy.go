package oauthserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

// mcpPath is where the MCP endpoint lives. The canonical resource URL is the
// configured issuer plus this path: one value, built once from configuration,
// never from a request (ADR-0019 addendum, #841).
const mcpPath = "/mcp"

// Resource is the canonical URL of the MCP endpoint (RFC 8707): the one resource
// a token of this server is for, and so the audience of every access token. A
// resource server compares Token.Audience with it.
func (s *Server) Resource() string { return s.resource }

// customScopes are the Scopes a user can consent to.
var customScopes = []string{ScopeCatalogRead, ScopeProgressRead, ScopeAdmin}

// standardScopes are the OpenID/OAuth scopes that are accepted but carry no
// grant here: openid only asks for an ID token (which discloses the subject and
// nothing else), offline_access for a refresh token. The flow works without
// openid.
var standardScopes = []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess}

// authorizer is the provider with this server's admission policy added. The
// library decides what to do with a scope or resource it does not know by
// silently dropping the scope and ignoring the resource; the policy here
// refuses both with a clear OAuth error instead, so a client never believes it
// holds something it was not given.
type authorizer struct {
	*op.Provider
	s *Server
}

var _ op.AuthorizeValidator = authorizer{}

type resourceKey struct{}

// ValidateAuthRequest is the library's validation preceded by the scope and
// resource policy. The redirect URI is checked first, so the refusals that
// follow can be sent to the client at it (RFC 6749 section 4.1.2.1).
func (a authorizer) ValidateAuthRequest(
	ctx context.Context, req *oidc.AuthRequest, storage op.Storage, verifier *op.IDTokenHintVerifier,
) (string, error) {
	c, err := storage.GetClientByClientID(ctx, req.ClientID)
	if err != nil {
		return "", oidc.ErrInvalidRequestRedirectURI().WithDescription("unable to retrieve client by id").WithParent(err)
	}
	if err := op.ValidateAuthReqRedirectURI(c, req.RedirectURI, req.ResponseType); err != nil {
		return "", err
	}
	for _, scope := range req.Scopes {
		switch {
		case slices.Contains(standardScopes, scope):
		case slices.Contains(customScopes, scope) && c.IsScopeAllowed(scope):
		default:
			return "", oidc.ErrInvalidScope().WithDescription("The scope %q is not available to this client. "+
				"Supported scopes: %s.", scope, strings.Join(append(slices.Clone(customScopes), standardScopes...), ", "))
		}
	}
	if resources, _ := ctx.Value(resourceKey{}).([]string); !a.s.acceptsResources(resources) {
		return "", errInvalidTarget(a.s.resource)
	}
	return op.ValidateAuthRequestClient(ctx, req, c, verifier)
}

// acceptsResources reports whether every requested resource is the canonical
// MCP URL. No resource at all is accepted and means the same thing: this server
// has one resource, and a token is bound to it regardless.
func (s *Server) acceptsResources(resources []string) bool {
	return !slices.ContainsFunc(resources, func(r string) bool { return r != s.resource })
}

func errInvalidTarget(resource string) *oidc.Error {
	return oidc.ErrInvalidTarget().WithDescription("The resource must be %s: this server issues tokens for that resource only.", resource)
}

// authorize serves the authorization endpoint through the library, with the
// resource parameter made available to the validator and the `iss` parameter
// (RFC 9207) added to any redirect to the client.
func (s *Server) authorize(w http.ResponseWriter, r *http.Request) {
	ctx := op.ContextWithIssuer(r.Context(), s.issuer)
	if err := r.ParseForm(); err == nil {
		ctx = context.WithValue(ctx, resourceKey{}, r.Form["resource"])
	}
	op.Authorize(s.issuerWriter(w), r.WithContext(ctx), authorizer{Provider: s.provider, s: s})
}

// authorizeCallback issues the code after consent and redirects to the client.
func (s *Server) authorizeCallback(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(op.ContextWithIssuer(r.Context(), s.issuer))
	op.AuthorizeCallback(s.issuerWriter(w), r, authorizer{Provider: s.provider, s: s})
}

// token serves the token endpoint. Only a resource that is not ours is refused
// here (RFC 8707 section 2.2); the audience itself comes from the grant.
func (s *Server) token(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.ParseForm() == nil && !s.acceptsResources(r.PostForm["resource"]) {
		op.RequestError(w, r, errInvalidTarget(s.resource), nil)
		return
	}
	s.provider.ServeHTTP(w, r)
}

// issuerWriter adds `iss` to a redirect that leaves for the client, success or
// error alike: a client that talks to several authorization servers uses it to
// tell a response meant for another one (RFC 9207). The redirect to our own
// consent screen is not a response to the client and is left alone.
func (s *Server) issuerWriter(w http.ResponseWriter) http.ResponseWriter {
	return &issWriter{ResponseWriter: w, issuer: s.issuer, own: s.issuer + consentPagePath}
}

type issWriter struct {
	http.ResponseWriter
	issuer, own string
}

func (w *issWriter) WriteHeader(code int) {
	if loc := w.Header().Get("Location"); code == http.StatusFound && loc != "" && !strings.HasPrefix(loc, w.own) {
		if u, err := url.Parse(loc); err == nil {
			q := u.Query()
			q.Set("iss", w.issuer)
			u.RawQuery = q.Encode()
			w.Header().Set("Location", u.String())
		}
	}
	w.ResponseWriter.WriteHeader(code)
}

// withIssParameter adds authorization_response_iss_parameter_supported (RFC
// 9207) to a metadata document; the library's discovery type has no field for it.
func withIssParameter(doc any) (map[string]any, error) {
	raw, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	out["authorization_response_iss_parameter_supported"] = true
	return out, nil
}
