package testsupport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// DocumentsHost is the hostname Client ID Metadata Documents appear to be served
// from, one the test certificate is valid for. It resolves to nothing: the test client dials the local server whatever
// name it is asked for.
const DocumentsHost = "example.com"

// ClientDocuments serves Client ID Metadata Documents to the authorization
// server over TLS. The production client refuses loopback by design, so the
// server under test is given a client (oauthserver.WithMetadataHTTPClient) that
// reaches this host instead; what the server does with the answer is the real
// code. A test changes what a path answers to move a document on, and advances
// the OAuth clock to let the cached copy expire.
type ClientDocuments struct {
	client *http.Client

	mu       sync.Mutex
	handlers map[string]http.HandlerFunc
}

func newClientDocuments(t *testing.T) *ClientDocuments {
	t.Helper()
	d := &ClientDocuments{handlers: map[string]http.HandlerFunc{}}
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d.mu.Lock()
		h := d.handlers[r.URL.Path]
		d.mu.Unlock()
		if h == nil {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}))
	t.Cleanup(srv.Close)

	roots := x509.NewCertPool()
	roots.AddCert(srv.Certificate())
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	d.client = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, srv.Listener.Addr().String())
			},
			TLSClientConfig:   &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12},
			DisableKeepAlives: true,
		},
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	return d
}

// URL is the client id of a document at a path.
func (d *ClientDocuments) URL(path string) string { return "https://" + DocumentsHost + path }

// Publish serves doc as JSON at a path (replacing what was there) and returns
// the client id, which is the document's URL. headers are added to the response.
// The document's own client_id is set to the URL unless doc names one.
func (d *ClientDocuments) Publish(path string, doc map[string]any, headers ...http.Header) string {
	id := d.URL(path)
	if _, ok := doc["client_id"]; !ok {
		doc["client_id"] = id
	}
	d.Serve(path, func(w http.ResponseWriter, _ *http.Request) {
		for _, h := range headers {
			for k, v := range h {
				w.Header()[k] = v
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(doc)
	})
	return id
}

// Serve answers a path with a handler of the test's own, for a document that is
// malformed, slow, redirected or too large.
func (d *ClientDocuments) Serve(path string, h http.HandlerFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[path] = h
}

// ValidDocument is a document of a client that asks for the given redirect URI.
func ValidDocument(redirectURI string) map[string]any {
	return map[string]any{
		"client_name":   "Document client",
		"redirect_uris": []string{redirectURI},
		"grant_types":   []string{"authorization_code", "refresh_token"},
	}
}
