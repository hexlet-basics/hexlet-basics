package oauthserver

import (
	"encoding/json"
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

// metadata serves the authorization server metadata (RFC 8414), at the
// well-known path and, for clients that fall back to it, at the OpenID one: both
// answer with the same document.
//
// It starts from the library's discovery document and narrows it to what this
// server does. The library advertises the implicit flow, ID-token response
// types, a device endpoint and client authentication methods that do not exist
// here, and a client that believed them would fail in confusing ways. The issuer
// comes from configuration through the context, never from the request.
func (s *Server) metadata(w http.ResponseWriter, r *http.Request) {
	ctx := op.ContextWithIssuer(r.Context(), s.issuer)
	doc := op.CreateDiscoveryConfig(ctx, s.provider, s.provider.Storage())

	doc.ResponseTypesSupported = []string{string(oidc.ResponseTypeCode)}
	doc.GrantTypesSupported = []oidc.GrantType{oidc.GrantTypeCode, oidc.GrantTypeRefreshToken}
	doc.CodeChallengeMethodsSupported = []oidc.CodeChallengeMethod{oidc.CodeChallengeMethodS256}
	doc.TokenEndpointAuthMethodsSupported = []oidc.AuthMethod{oidc.AuthMethodNone}
	doc.RevocationEndpointAuthMethodsSupported = []oidc.AuthMethod{oidc.AuthMethodNone}
	doc.TokenEndpointAuthSigningAlgValuesSupported = nil
	doc.RevocationEndpointAuthSigningAlgValuesSupported = nil
	doc.IntrospectionEndpoint = ""
	doc.IntrospectionEndpointAuthMethodsSupported = nil
	doc.IntrospectionEndpointAuthSigningAlgValuesSupported = nil
	doc.DeviceAuthorizationEndpoint = ""
	doc.CheckSessionIframe = ""
	doc.ClaimsSupported = nil
	doc.RequestParameterSupported = false
	doc.RequestObjectSigningAlgValuesSupported = nil

	out, err := withIssParameter(doc)
	if err != nil {
		http.Error(w, "metadata", http.StatusInternalServerError)
		return
	}
	// RFC 7591: where a client registers itself. The library's discovery type has
	// no field for it.
	out["registration_endpoint"] = s.issuer + registerPath
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_ = json.NewEncoder(w).Encode(out)
}
