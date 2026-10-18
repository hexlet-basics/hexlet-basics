package oauthserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/rs/cors"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"hexletbasics/internal/ident"
)

const (
	registerPath = "/oauth/register"

	// maxRegistrationBody is far above any legitimate registration (a few URIs
	// and a name) and far below what is worth parsing from an anonymous caller.
	maxRegistrationBody = 16 << 10
	maxRedirectURIs     = 10
	maxRedirectURILen   = 2048
	maxClientNameRunes  = 100

	defaultClientName = "Unnamed application"
)

// defaultRegistrationScopes is what a registration that names no scope may
// request later: reading, never admin. A client that wants more asks for it by
// name, and the user still consents and the user's Permissions still decide
// (ADR-0022 addendum, #840).
func defaultRegistrationScopes() []string {
	return []string{ScopeCatalogRead, ScopeProgressRead}
}

// registrationFields are the metadata (RFC 7591) a registration is judged on;
// ignoredRegistrationFields are harmless descriptive metadata that is accepted
// and dropped. Every other field is refused rather than ignored: this endpoint
// is anonymous, and a field it does not understand (an audience, a secret,
// skip_consent) is exactly what an attacker would try, so silence would let a
// future change give one of them meaning by accident.
var ()

// refusal is a validation failure whose text is the answer to the client (the
// error_description of an RFC 7591 response), so it is a message built for the
// reader rather than an error message with data in attributes (ADR-0018).
type refusal string

func (r refusal) Error() string { return string(r) }

// registrationError is the RFC 7591 error response.
type registrationError struct {
	status      int
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

func (e *registrationError) Error() string { return e.Code + ": " + e.Description }

func metadataError(format string, args ...any) *registrationError {
	return &registrationError{http.StatusBadRequest, "invalid_client_metadata", fmt.Sprintf(format, args...)}
}

func redirectError(format string, args ...any) *registrationError {
	return &registrationError{http.StatusBadRequest, "invalid_redirect_uri", fmt.Sprintf(format, args...)}
}

// registration is a validated request.
type registration struct {
	name         string
	redirectURIs []string
	// scopes are the custom Scopes the client row stores; offlineAccess is the
	// standard scope, admitted for every client and so not stored.
	scopes        []string
	offlineAccess bool
}

// register serves Dynamic Client Registration (RFC 7591): anonymous, so the
// request is judged against what a public client of this server may be and
// nothing else. The rate limit is spent first and by every attempt, a refused
// one included, so probing for a gap in the validation is not free.
func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		writeRegistrationError(w, &registrationError{http.StatusMethodNotAllowed, "invalid_request", "use POST"})

		return
	}

	address, _ := s.clientAddress(r.Context())
	if ok, wait := s.limiter.allow(address); !ok {
		w.Header().Set("Retry-After", strconv.Itoa(int(wait/time.Second)+1))
		writeRegistrationError(w, &registrationError{
			http.StatusTooManyRequests, "temporarily_unavailable", "too many registrations from this address, try again later",
		})

		return
	}

	reg, regErr := s.readRegistration(w, r)
	if regErr != nil {
		writeRegistrationError(w, regErr)

		return
	}

	clientID := ident.New()

	_, err := s.storage.db.OAuthClient.Create().
		SetClientID(clientID).
		SetName(reg.name).
		SetRedirectUris(reg.redirectURIs).
		SetScopes(reg.scopes).
		Save(r.Context())
	if err != nil {
		writeRegistrationError(w, &registrationError{http.StatusInternalServerError, "server_error", "could not register the client"})

		return
	}

	granted := slices.Clone(reg.scopes)
	if reg.offlineAccess {
		granted = append(granted, oidc.ScopeOfflineAccess)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]any{
		"client_id":                  clientID,
		"client_id_issued_at":        s.storage.now().Unix(),
		"client_name":                reg.name,
		"redirect_uris":              reg.redirectURIs,
		"grant_types":                []string{string(oidc.GrantTypeCode), string(oidc.GrantTypeRefreshToken)},
		"response_types":             []string{string(oidc.ResponseTypeCode)},
		"token_endpoint_auth_method": string(oidc.AuthMethodNone),
		"scope":                      strings.Join(granted, " "),
	})
}

// registerHandler is register behind CORS: browser-based MCP clients register
// from their own origin. Nothing here is cookie-authenticated, so credentials
// are never allowed.
func (s *Server) registerHandler() http.Handler {
	return cors.New(cors.Options{
		AllowOriginFunc: func(string) bool { return true },
		AllowedMethods:  []string{http.MethodPost},
		AllowedHeaders:  []string{"Content-Type"},
	}).Handler(http.HandlerFunc(s.register))
}

func writeRegistrationError(w http.ResponseWriter, e *registrationError) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(e.status)
	writeJSON(w, e)
}

// readRegistration reads and validates the body.
func (s *Server) readRegistration(w http.ResponseWriter, r *http.Request) (registration, *registrationError) {
	if mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err != nil || mt != "application/json" {
		return registration{}, &registrationError{
			http.StatusUnsupportedMediaType, "invalid_client_metadata", "the request must be application/json",
		}
	}

	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxRegistrationBody))
	if err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			return registration{}, &registrationError{
				http.StatusRequestEntityTooLarge, "invalid_client_metadata", "the request is larger than it may be",
			}
		}

		return registration{}, metadataError("the request body could not be read")
	}

	dec := json.NewDecoder(bytes.NewReader(raw))

	var fields map[string]json.RawMessage
	if err := dec.Decode(&fields); err != nil || fields == nil {
		return registration{}, metadataError("the body must be a JSON object")
	}

	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return registration{}, metadataError("the body must be a single JSON object")
	}

	return parseMetadata(fields)
}

// parseMetadata judges a JSON object of client metadata, from a registration
// request or from a Client ID Metadata Document, by the one allow-list: both are
// anonymous input describing a client, so both are held to the same constraints.
// Fields other than those this server understands are refused, never ignored.
func parseMetadata(fields map[string]json.RawMessage, extraFields ...string) (registration, *registrationError) {
	for name := range fields {
		if !slices.Contains(registrationFields(), name) && !slices.Contains(ignoredRegistrationFields(), name) &&
			!slices.Contains(extraFields, name) {
			return registration{}, metadataError("the field %q is not accepted", name)
		}
	}

	return parseRegistration(fields)
}

func parseRegistration(fields map[string]json.RawMessage) (registration, *registrationError) {
	var reg registration

	uris, err := stringList(fields, "redirect_uris")
	if err != nil {
		return reg, metadataError("redirect_uris must be an array of strings")
	}

	if len(uris) == 0 || len(uris) > maxRedirectURIs {
		return reg, redirectError("redirect_uris must list between 1 and %d URIs", maxRedirectURIs)
	}

	for _, uri := range uris {
		if err := checkRedirectURI(uri); err != nil {
			return reg, redirectError("%s", err)
		}
	}

	reg.redirectURIs = uris

	if reg.name, err = clientName(fields); err != nil {
		return reg, metadataError("%s", err)
	}

	// The grants a client row has are fixed (code and refresh), so a request is
	// accepted when it asks for no more than that, and refused otherwise instead
	// of being quietly narrowed: a client that wanted client_credentials must
	// learn it will not get it.
	grants, err := stringList(fields, "grant_types")
	if err != nil {
		return reg, metadataError("grant_types must be an array of strings")
	}

	for _, g := range grants {
		if g != string(oidc.GrantTypeCode) && g != string(oidc.GrantTypeRefreshToken) {
			return reg, metadataError("the grant type %q is not supported: only authorization_code and refresh_token", g)
		}
	}

	if len(grants) > 0 && !slices.Contains(grants, string(oidc.GrantTypeCode)) {
		return reg, metadataError("grant_types must include authorization_code")
	}

	responses, err := stringList(fields, "response_types")
	if err != nil {
		return reg, metadataError("response_types must be an array of strings")
	}

	for _, rt := range responses {
		if rt != string(oidc.ResponseTypeCode) {
			return reg, metadataError("the response type %q is not supported: only code", rt)
		}
	}

	if raw, ok := fields["token_endpoint_auth_method"]; ok {
		var method string
		if json.Unmarshal(raw, &method) != nil || method != string(oidc.AuthMethodNone) {
			return reg, metadataError("only public clients are registered: token_endpoint_auth_method must be none")
		}
	}

	if reg.scopes, reg.offlineAccess, err = requestedScopes(fields); err != nil {
		return reg, metadataError("%s", err)
	}

	return reg, nil
}

// stringList reads an optional array of strings; absent is an empty list.
func stringList(fields map[string]json.RawMessage, name string) ([]string, error) {
	raw, ok := fields[name]
	if !ok {
		return nil, nil
	}

	var list []string
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}

	return list, nil
}

// clientName reads the display name the consent screen shows. It is stored as
// plain text with control characters removed; escaping is the screen's job, but
// a name that cannot hold a NUL or a newline cannot break a log line or a layout
// either.
func clientName(fields map[string]json.RawMessage) (string, error) {
	raw, ok := fields["client_name"]
	if !ok {
		return defaultClientName, nil
	}

	var name string
	if err := json.Unmarshal(raw, &name); err != nil {
		return "", refusal("client_name must be a string")
	}

	name = strings.TrimSpace(strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}

		return r
	}, name))
	if utf8.RuneCountInString(name) > maxClientNameRunes {
		return "", refusal(fmt.Sprintf("client_name may be at most %d characters", maxClientNameRunes))
	}

	if name == "" {
		return defaultClientName, nil
	}

	return name, nil
}

// requestedScopes validates the space-separated scope string against the
// offered list: the three Scopes and offline_access, nothing else (not even
// openid or profile, which a client row cannot use). No scope at all means the
// read-only default.
func requestedScopes(fields map[string]json.RawMessage) ([]string, bool, error) {
	raw, ok := fields["scope"]
	if !ok {
		return slices.Clone(defaultRegistrationScopes()), false, nil
	}

	var scope string
	if json.Unmarshal(raw, &scope) != nil {
		return nil, false, refusal("scope must be a space-separated string")
	}

	names := strings.Fields(scope)
	if len(names) == 0 {
		return slices.Clone(defaultRegistrationScopes()), false, nil
	}
	// Never nil: the column is NOT NULL, and a client that asked only for
	// offline_access holds no custom Scope.
	custom := []string{}

	var offline bool

	for _, name := range names {
		switch {
		case name == oidc.ScopeOfflineAccess:
			offline = true
		case slices.Contains(customScopes(), name):
			if !slices.Contains(custom, name) {
				custom = append(custom, name)
			}
		default:
			return nil, false, refusal(fmt.Sprintf("the scope %q is not offered. Offered scopes: %s",
				name, strings.Join(append(slices.Clone(customScopes()), oidc.ScopeOfflineAccess), ", ")))
		}
	}

	return custom, offline, nil
}

// checkRedirectURI admits an https URI, or an http one on a loopback host (RFC
// 8252: a native client listens on its own machine). Everything else is
// refused: plain http elsewhere is readable on the wire, custom schemes cannot
// be told from a hijacker's, and a wildcard would let any path or subdomain
// receive the code. Userinfo is refused because it is how a lookalike host such
// as http://127.0.0.1@evil.example passes a prefix check, and fragments because
// RFC 6749 forbids them in a redirect URI.
func checkRedirectURI(raw string) error {
	if len(raw) > maxRedirectURILen {
		return refusal("a redirect URI is too long")
	}

	lower := strings.ToLower(raw)
	if strings.Contains(raw, "*") || strings.Contains(lower, "%2a") {
		return refusal(fmt.Sprintf("wildcards are not allowed in a redirect URI: %q", raw))
	}

	u, err := url.Parse(raw)
	if err != nil || !u.IsAbs() || u.Opaque != "" || u.Host == "" || u.Hostname() == "" {
		return refusal(fmt.Sprintf("a redirect URI must be an absolute URL: %q", raw))
	}

	if u.User != nil {
		return refusal(fmt.Sprintf("a redirect URI must not carry credentials: %q", raw))
	}

	if u.Fragment != "" || strings.Contains(raw, "#") {
		return refusal(fmt.Sprintf("a redirect URI must not have a fragment: %q", raw))
	}

	switch strings.ToLower(u.Scheme) {
	case "https":
		return nil
	case "http":
		if isLoopbackHost(u.Hostname()) {
			return nil
		}

		return refusal(fmt.Sprintf("plain http is allowed only for a loopback redirect: %q", raw))
	default:
		return refusal(fmt.Sprintf("a redirect URI must use https, or http on a loopback address: %q", raw))
	}
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}

	ip := net.ParseIP(host)

	return ip != nil && ip.IsLoopback()
}

func registrationFields() []string {
	return []string{
		"redirect_uris", "client_name", "grant_types", "response_types",
		"token_endpoint_auth_method", "scope",
	}
}

func ignoredRegistrationFields() []string {
	return []string{
		"client_uri", "logo_uri", "tos_uri", "policy_uri", "contacts",
		"software_id", "software_version", "application_type",
	}
}
