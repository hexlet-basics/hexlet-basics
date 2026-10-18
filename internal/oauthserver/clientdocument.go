package oauthserver

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/singleflight"

	"hexletbasics/ent"
	"hexletbasics/ent/oauthclient"
)

// Client ID Metadata Documents (ADR-0022 addendum, #847). A client whose id is
// an https URL is registered by fetching that URL the first time its id is seen
// at the authorization endpoint. The fetch is triggered by an anonymous request
// and goes to an address the requester chose, so it is the most dangerous thing
// this package does: it is the classic SSRF shape. Everything in this file that
// is not plain bookkeeping exists to keep that fetch on the public internet.

const (
	// maxDocumentBody is how much of a document is read. A legitimate one is a
	// few hundred bytes; the cap stops a hostile host from streaming us memory.
	maxDocumentBody = 16 << 10
	// documentTimeout bounds the whole fetch (dial, TLS, headers, body), so a
	// host that stalls cannot hold a request open.
	documentTimeout = 5 * time.Second
	maxClientIDLen  = 2048

	// The cache lifetime follows the document's HTTP headers, clamped: never so
	// short that a client re-fetches on every request (the floor also covers
	// no-store, since the row has to be kept to run the flow at all), never so
	// long that a changed or withdrawn document is trusted for days.
	minDocumentTTL     = 5 * time.Minute
	defaultDocumentTTL = time.Hour
	maxDocumentTTL     = 24 * time.Hour
	// failureTTL is how long a failed fetch is remembered, so a URL that does not
	// work (or an attacker repeating one) cannot make us fetch it again on every
	// request.
	failureTTL = 5 * time.Minute
	// maxRememberedFailures bounds that memory; beyond it failures are simply not
	// remembered and the per-address limiter is what holds.
	maxRememberedFailures = 10_000
)

// blockedNetworks are the destinations a metadata document must never be fetched
// from, beyond what netip classifies as private, loopback, link-local, multicast
// or unspecified. The IPv6 entries are the forms that carry an IPv4 address
// inside them and so could name a private host indirectly.
var blockedNetworks = parsePrefixes(
	"0.0.0.0/8",       // "this network"
	"100.64.0.0/10",   // carrier-grade NAT, internal to many providers
	"192.0.0.0/24",    // IETF protocol assignments
	"192.0.2.0/24",    // documentation
	"198.18.0.0/15",   // benchmarking
	"198.51.100.0/24", // documentation
	"203.0.113.0/24",  // documentation
	"240.0.0.0/4",     // reserved, includes the broadcast address
	"64:ff9b::/96",    // NAT64: embeds an IPv4 address
	"64:ff9b:1::/48",  // local-use NAT64
	"100::/64",        // discard
	"2001::/32",       // Teredo: embeds an IPv4 address
	"2001:db8::/32",   // documentation
	"2002::/16",       // 6to4: embeds an IPv4 address
	"fec0::/10",       // deprecated site-local
)

func parsePrefixes(list ...string) []netip.Prefix {
	out := make([]netip.Prefix, len(list))
	for i, p := range list {
		out[i] = netip.MustParsePrefix(p)
	}

	return out
}

// publicAddress reports whether a metadata document may be fetched from ip. It
// is an allow-list in spirit: only a global unicast address outside every
// blocked range passes, so a range nobody thought of is refused, not trusted.
// An IPv4-mapped IPv6 address (::ffff:10.0.0.1) is judged as the IPv4 address it
// carries, since that is where the connection really goes.
func publicAddress(ip netip.Addr) bool {
	ip = ip.Unmap()
	if !ip.IsValid() || !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() || ip.IsMulticast() || ip.IsUnspecified() {
		return false
	}

	for _, p := range blockedNetworks {
		if p.Contains(ip) {
			return false
		}
	}

	return true
}

// newGuardedClient is the HTTP client documents are fetched with. The check on
// the destination runs in the dialer's Control hook, which sees the IP address
// the socket is about to connect to, after name resolution and for every
// candidate address. Checking a hostname beforehand would leave a gap for DNS
// rebinding (a name that resolves to a public address for the check and to
// 169.254.169.254 for the connection); checking at dial time has none.
// allow decides which addresses pass; production uses publicAddress.
//
// Redirects are refused outright rather than followed and re-validated: a
// metadata document has no business moving, and each hop followed is another
// chance to be wrong.
func newGuardedClient(allow func(netip.Addr) bool, timeout time.Duration) *http.Client {
	dialer := &net.Dialer{
		Timeout: timeout,
		Control: func(_, address string, _ syscall.RawConn) error {
			ap, err := netip.ParseAddrPort(address)
			if err != nil {
				return errs.With("address", address).Wrapf(err, "refusing to connect")
			}

			if !allow(ap.Addr()) {
				return errs.With("address", ap.Addr().String()).Errorf("refusing to connect: not a public address")
			}

			return nil
		},
	}

	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			// No proxy: one would be dialed instead of the target, and the target
			// would never reach the check above.
			Proxy:                  nil,
			DialContext:            dialer.DialContext,
			ForceAttemptHTTP2:      true,
			TLSClientConfig:        &tls.Config{MinVersion: tls.VersionTLS12},
			TLSHandshakeTimeout:    timeout,
			MaxResponseHeaderBytes: maxDocumentHeaderBytes,
			DisableKeepAlives:      true,
		},
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
}

// Option changes how New builds a Server beyond its Deps. Production passes
// none: every option exists for a test.
type Option func(*Server)

// WithMetadataHTTPClient replaces the client documents are fetched with. It is
// for tests, which need to reach a server on loopback that the production
// client refuses by design; nothing in the application wiring uses it.
func WithMetadataHTTPClient(c *http.Client) Option {
	return func(s *Server) { s.metadataClient = c }
}

// documentError is a refusal to use a metadata document, with the response the
// authorization endpoint gives for it.
type documentError struct {
	status     int
	retryAfter time.Duration
	reason     string
}

func (e *documentError) Error() string { return e.reason }

func invalidDocument(format string, args ...any) *documentError {
	return &documentError{status: http.StatusBadRequest, reason: fmt.Sprintf(format, args...)}
}

// documentClients resolves client ids that are URLs into client rows.
type documentClients struct {
	client  *http.Client
	limiter *limiter
	db      *ent.Client
	now     func() time.Time
	flight  singleflight.Group

	mu       sync.Mutex
	failures map[string]time.Time
}

func newDocumentClients(client *http.Client, lim *limiter, db *ent.Client, now func() time.Time) *documentClients {
	if client == nil {
		client = newGuardedClient(publicAddress, documentTimeout)
	}

	return &documentClients{client: client, limiter: lim, db: db, now: now, failures: map[string]time.Time{}}
}

// isDocumentID reports whether a client id is the kind this file handles. Any
// other id is left to the stored clients, so a typo does not start a fetch.
func isDocumentID(clientID string) bool {
	return len(clientID) >= len("https://") && strings.EqualFold(clientID[:len("https://")], "https://")
}

// What checkClientIDURL answers, in the order it checks.
var (
	errClientIDTooLong     = errors.New("the client_id URL is too long")
	errClientIDNotHTTPS    = errors.New("the client_id must be an absolute https URL")
	errClientIDCredentials = errors.New("the client_id URL must not carry credentials")
	errClientIDFragment    = errors.New("the client_id URL must not have a fragment")
	errClientIDPort        = errors.New("the client_id URL must use the default https port")
	errClientIDNoPath      = errors.New("the client_id URL must have a path")
	errClientIDDotSegments = errors.New("the client_id URL must not contain dot segments")
)

// checkClientIDURL holds a client id to what the specification requires of one,
// plus the rules of this server: https, a path, no credentials, no fragment, no
// dot segments, and the default port, so the URL names an ordinary public web
// resource and nothing else.
func checkClientIDURL(clientID string) error {
	if len(clientID) > maxClientIDLen {
		return errClientIDTooLong
	}

	u, err := url.Parse(clientID)
	if err != nil || u.Scheme != "https" || u.Opaque != "" || u.Hostname() == "" {
		return errClientIDNotHTTPS
	}

	if u.User != nil {
		return errClientIDCredentials
	}

	if u.Fragment != "" || strings.Contains(clientID, "#") {
		return errClientIDFragment
	}

	if port := u.Port(); port != "" && port != "443" {
		return errClientIDPort
	}

	if u.Path == "" || u.Path == "/" {
		return errClientIDNoPath
	}

	for segment := range strings.SplitSeq(u.Path, "/") {
		if segment == "." || segment == ".." {
			return errClientIDDotSegments
		}
	}

	return nil
}

// ensure makes the row of a metadata-document client current: it fetches the
// document when the id is new or the cached copy has expired, and does nothing
// otherwise (including for any id that is not a document URL). Only the
// authorization request calls it. The token endpoint uses the row as it is, so a
// flow already under way never depends on a remote host being up, and a code
// that was issued cannot be stranded by an outage.
//
// When a refresh fails the row is left alone, and the authorization request is
// refused: an application whose document can no longer be verified is not
// trusted to start a new grant, while the grants it already holds keep working
// until the user revokes them.
func (d *documentClients) ensure(ctx context.Context, clientID, address string) *documentError {
	if !isDocumentID(clientID) {
		return nil
	}

	if err := checkClientIDURL(clientID); err != nil {
		return invalidDocument("%s", err)
	}

	now := d.now()

	row, err := d.db.OAuthClient.Query().Where(oauthclient.ClientID(clientID)).Only(ctx)
	switch {
	case err == nil:
		if row.MetadataExpiresAt == nil || now.Before(*row.MetadataExpiresAt) {
			return nil
		}
	case !ent.IsNotFound(err):
		return &documentError{status: http.StatusInternalServerError, reason: "could not load the client"}
	}

	if d.failed(clientID, now) {
		return invalidDocument("the client metadata document could not be used, try again later")
	}
	// The budget is spent only when a fetch would really leave the server, so a
	// cached or remembered-as-failed id costs nothing and the limit is on the
	// thing worth limiting.
	if ok, wait := d.limiter.allow(address); !ok {
		return &documentError{http.StatusTooManyRequests, wait, "too many client metadata fetches from this address, try again later"}
	}

	_, err, _ = d.flight.Do(clientID, func() (any, error) { return nil, d.refresh(ctx, clientID) })
	if err != nil {
		if de, ok := errors.AsType[*documentError](err); ok {
			d.remember(clientID, now)

			return de
		}

		return &documentError{status: http.StatusInternalServerError, reason: "could not store the client"}
	}

	return nil
}

func (d *documentClients) failed(clientID string, now time.Time) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	until, ok := d.failures[clientID]
	if ok && !now.Before(until) {
		delete(d.failures, clientID)

		return false
	}

	return ok
}

func (d *documentClients) remember(clientID string, now time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.failures) >= maxRememberedFailures {
		for id, until := range d.failures {
			if !now.Before(until) {
				delete(d.failures, id)
			}
		}

		if len(d.failures) >= maxRememberedFailures {
			return
		}
	}

	d.failures[clientID] = now.Add(failureTTL)
}

// refresh fetches, validates and stores the document. A *documentError is a
// problem with the document or its host; anything else is ours.
func (d *documentClients) refresh(ctx context.Context, clientID string) error {
	reg, ttl, derr := d.fetch(ctx, clientID)
	if derr != nil {
		return derr
	}

	expires := d.now().Add(ttl)

	n, err := d.db.OAuthClient.Update().Where(oauthclient.ClientID(clientID)).
		SetName(reg.name).SetRedirectUris(reg.redirectURIs).SetScopes(reg.scopes).SetMetadataExpiresAt(expires).
		Save(ctx)
	if err != nil {
		return errs.Wrapf(err, "update client from metadata document")
	}

	if n > 0 {
		return nil
	}

	_, err = d.db.OAuthClient.Create().
		SetClientID(clientID).SetName(reg.name).SetRedirectUris(reg.redirectURIs).SetScopes(reg.scopes).
		SetMetadataExpiresAt(expires).Save(ctx)
	if ent.IsConstraintError(err) {
		// Another replica created it between our update and our insert.
		_, err = d.db.OAuthClient.Update().Where(oauthclient.ClientID(clientID)).
			SetName(reg.name).SetRedirectUris(reg.redirectURIs).SetScopes(reg.scopes).SetMetadataExpiresAt(expires).
			Save(ctx)
	}

	if err != nil {
		return errs.Wrapf(err, "create client from metadata document")
	}

	return nil
}

// fetch retrieves a document and judges it by the rules of dynamic
// registration. The reason for a refusal is shown to the caller, who chose the
// URL and can see the document, and carries nothing from our side.
func (d *documentClients) fetch(ctx context.Context, clientID string) (registration, time.Duration, *documentError) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, clientID, nil) //nolint:gosec // G704: the client_id URL is caller input by design; the client dials through the private-address-refusing Control hook, see the file comment
	if err != nil {
		return registration{}, 0, invalidDocument("the client_id URL could not be requested")
	}

	req.Header.Set("Accept", "application/json")

	res, err := d.client.Do(req) //nolint:gosec // G704: guarded by the dialer Control hook that refuses private, loopback and link-local addresses
	if err != nil {
		return registration{}, 0, invalidDocument("the client metadata document could not be fetched")
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return registration{}, 0, invalidDocument("the client metadata document answered %d: redirects are not followed and only 200 is accepted", res.StatusCode)
	}

	if mt, _, err := mime.ParseMediaType(res.Header.Get("Content-Type")); err != nil || mt != "application/json" {
		return registration{}, 0, invalidDocument("the client metadata document must be application/json")
	}

	raw, err := io.ReadAll(io.LimitReader(res.Body, maxDocumentBody+1))
	if err != nil {
		return registration{}, 0, invalidDocument("the client metadata document could not be read")
	}

	if len(raw) > maxDocumentBody {
		return registration{}, 0, invalidDocument("the client metadata document is larger than %d bytes", maxDocumentBody)
	}

	dec := json.NewDecoder(bytes.NewReader(raw))

	var fields map[string]json.RawMessage
	if err := dec.Decode(&fields); err != nil || fields == nil {
		return registration{}, 0, invalidDocument("the client metadata document must be a JSON object")
	}

	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return registration{}, 0, invalidDocument("the client metadata document must be a single JSON object")
	}
	// The document vouches for the URL it was fetched from: its own client_id
	// must be that URL, character for character, or anyone could publish a
	// document naming somebody else's id.
	var documented string
	if json.Unmarshal(fields["client_id"], &documented) != nil || documented != clientID {
		return registration{}, 0, invalidDocument("the client_id of the document must equal the URL it was fetched from")
	}

	delete(fields, "client_id")

	reg, regErr := parseMetadata(fields)
	if regErr != nil {
		return registration{}, 0, invalidDocument("%s", regErr.Description)
	}

	return reg, d.cacheTTL(res.Header), nil
}

// cacheTTL is how long a fetched document is trusted, from its HTTP caching
// headers (s-maxage, max-age, then Expires), clamped to the bounds above.
// no-store and no-cache cannot be honoured literally, because the row must
// exist to run the flow, so they mean the shortest allowed time.
func (d *documentClients) cacheTTL(h http.Header) time.Duration {
	ttl := defaultDocumentTTL

	var maxAge, sMaxAge *time.Duration

	immediate := false

	for directive := range strings.SplitSeq(h.Get("Cache-Control"), ",") {
		name, value, _ := strings.Cut(strings.ToLower(strings.TrimSpace(directive)), "=")
		switch name {
		case "no-store", "no-cache":
			immediate = true
		case "max-age", "s-maxage":
			if secs, err := strconv.ParseInt(strings.Trim(value, `"`), 10, 64); err == nil && secs >= 0 {
				v := time.Duration(min(secs, int64(maxDocumentTTL/time.Second))) * time.Second
				if name == "max-age" {
					maxAge = &v
				} else {
					sMaxAge = &v
				}
			}
		}
	}

	switch {
	case immediate:
		ttl = 0
	case sMaxAge != nil:
		ttl = *sMaxAge
	case maxAge != nil:
		ttl = *maxAge
	default:
		if expires, err := http.ParseTime(h.Get("Expires")); err == nil {
			ttl = expires.Sub(d.now())
		}
	}

	return min(max(ttl, minDocumentTTL), maxDocumentTTL)
}

// writeDocumentError answers an authorization request whose client id could not
// be resolved. The redirect URI of an unverified client is not trusted, so the
// answer is direct, never a redirect.
func writeDocumentError(w http.ResponseWriter, e *documentError) {
	code := "invalid_request"
	if e.status == http.StatusTooManyRequests {
		code = "temporarily_unavailable"

		w.Header().Set("Retry-After", strconv.Itoa(int(e.retryAfter/time.Second)+1))
	}

	if e.status >= http.StatusInternalServerError {
		code = "server_error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(e.status)
	writeJSON(w, map[string]string{"error": code, "error_description": e.reason})
}

// maxDocumentHeaderBytes bounds the response headers of a client metadata
// document.
const maxDocumentHeaderBytes = 16 << 10
