package oauthserver

import (
	"net/url"
	"slices"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"hexletbasics/ent"
)

// client adapts a registered application to op.Client. All of them are public
// clients using the authorization-code grant, with refresh, and PKCE S256 (the
// last enforced in CreateAuthRequest).
type client struct {
	row *ent.OAuthClient
	// loginBase is the consent page; the library appends the request id.
	loginBase string
}

var _ op.Client = (*client)(nil)

func (c *client) GetID() string                       { return c.row.ClientID }
func (c *client) RedirectURIs() []string              { return c.row.RedirectUris }
func (c *client) PostLogoutRedirectURIs() []string    { return nil }
func (c *client) ApplicationType() op.ApplicationType { return op.ApplicationTypeNative }
func (c *client) AuthMethod() oidc.AuthMethod         { return oidc.AuthMethodNone }

// ResponseTypes is code only: the implicit flow has no place in OAuth 2.1.
func (c *client) ResponseTypes() []oidc.ResponseType {
	return []oidc.ResponseType{oidc.ResponseTypeCode}
}
func (c *client) GrantTypes() []oidc.GrantType {
	return []oidc.GrantType{oidc.GrantTypeCode, oidc.GrantTypeRefreshToken}
}

// LoginURL sends the user to the consent screen, an SSR route on the canonical
// origin that signs them in first when needed.
func (c *client) LoginURL(requestID string) string {
	return c.loginBase + "?id=" + url.QueryEscape(requestID)
}

// AccessTokenType is bearer: the opaque, encrypted id the database can verify by
// lookup, rather than a self-contained JWT that outlives its revocation.
func (c *client) AccessTokenType() op.AccessTokenType { return op.AccessTokenTypeBearer }
func (c *client) IDTokenLifetime() time.Duration      { return time.Hour }
func (c *client) DevMode() bool                       { return false }
func (c *client) RestrictAdditionalIdTokenScopes() func([]string) []string {
	return func(scopes []string) []string { return scopes }
}
func (c *client) RestrictAdditionalAccessTokenScopes() func([]string) []string {
	return func(scopes []string) []string { return scopes }
}

// IsScopeAllowed admits a custom scope only when the client was registered with
// it, so a client cannot be shown or granted more than it declared. The policy
// for what a registration may declare is a later step of ADR-0020.
func (c *client) IsScopeAllowed(scope string) bool     { return slices.Contains(c.row.Scopes, scope) }
func (c *client) IDTokenUserinfoClaimsAssertion() bool { return false }
func (c *client) ClockSkew() time.Duration             { return 0 }
