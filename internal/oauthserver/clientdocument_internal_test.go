package oauthserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicAddressAdmitsOnlyGlobalUnicast(t *testing.T) {
	for _, ip := range []string{"93.184.216.34", "1.1.1.1", "2606:2800:220:1::1", "::ffff:93.184.216.34"} {
		assert.True(t, publicAddress(netip.MustParseAddr(ip)), ip)
	}
	for _, ip := range []string{
		"127.0.0.1", "127.255.255.254", "::1", "0.0.0.0", "::", "10.1.2.3", "172.31.0.1", "192.168.0.1",
		"100.64.0.1", "169.254.169.254", "169.254.0.1", "fe80::1", "fd12:3456::1", "fc00::1",
		"224.0.0.1", "239.255.255.250", "ff02::1", "255.255.255.255", "240.0.0.1", "198.18.0.1", "192.0.2.1",
		"::ffff:127.0.0.1", "::ffff:10.0.0.1", "::ffff:169.254.169.254", "64:ff9b::a00:1", "2002:a00:1::1",
		"2001:0:4136:e378:8000:63bf:3fff:fdd2", "fec0::1",
	} {
		assert.False(t, publicAddress(netip.MustParseAddr(ip)), ip)
	}
}

// loopbackServer is an HTTPS server on loopback and a client trusting it.
func loopbackServer(t *testing.T, h http.HandlerFunc) (*httptest.Server, *tls.Config) {
	t.Helper()
	srv := httptest.NewTLSServer(h)
	t.Cleanup(srv.Close)
	roots := x509.NewCertPool()
	roots.AddCert(srv.Certificate())
	return srv, &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
}

func get(ctx context.Context, c *http.Client, target string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func TestTheProductionClientRefusesALoopbackServerAtDialTime(t *testing.T) {
	srv, tlsConfig := loopbackServer(t, func(w http.ResponseWriter, _ *http.Request) {
		t.Error("the server must never be reached")
	})
	client := newGuardedClient(publicAddress, time.Second)
	client.Transport.(*http.Transport).TLSClientConfig = tlsConfig

	res, err := get(t.Context(), client, srv.URL)

	if res != nil {
		_ = res.Body.Close()
	}
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a public address")
}

// A name is resolved by the dialer, and the address it gets is what is judged:
// "localhost" is the stand-in for a name an attacker points at a private host
// (DNS rebinding), which a check of the hostname beforehand cannot see.
func TestAHostnameThatResolvesToLoopbackIsRefusedAtDialTime(t *testing.T) {
	srv, tlsConfig := loopbackServer(t, func(http.ResponseWriter, *http.Request) {
		t.Error("the server must never be reached")
	})
	client := newGuardedClient(publicAddress, time.Second)
	client.Transport.(*http.Transport).TLSClientConfig = tlsConfig
	_, port, found := cutLast(srv.URL, ':')
	require.True(t, found)

	res, err := get(t.Context(), client, "https://localhost:"+port+"/")

	if res != nil {
		_ = res.Body.Close()
	}
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a public address")
}

func cutLast(s string, sep byte) (before, after string, found bool) {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == sep {
			return s[:i], s[i+1:], true
		}
	}
	return s, "", false
}

func TestTheGuardedClientDoesNotFollowRedirects(t *testing.T) {
	srv, tlsConfig := loopbackServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/final" {
			t.Error("a redirect must not be followed")
		}
		http.Redirect(w, r, "/final", http.StatusFound)
	})
	client := newGuardedClient(func(netip.Addr) bool { return true }, time.Second)
	client.Transport.(*http.Transport).TLSClientConfig = tlsConfig

	res, err := get(t.Context(), client, srv.URL+"/start")

	require.NoError(t, err)
	defer func() { _ = res.Body.Close() }()
	assert.Equal(t, http.StatusFound, res.StatusCode, "the redirect itself is the answer, and is then refused as not a 200")
}

func TestTheGuardedClientGivesUpOnAHostThatStalls(t *testing.T) {
	release := make(chan struct{})
	srv, tlsConfig := loopbackServer(t, func(http.ResponseWriter, *http.Request) { <-release })
	t.Cleanup(func() { close(release) })
	client := newGuardedClient(func(netip.Addr) bool { return true }, 100*time.Millisecond)
	client.Transport.(*http.Transport).TLSClientConfig = tlsConfig

	start := time.Now()
	res, err := get(t.Context(), client, srv.URL)

	if res != nil {
		_ = res.Body.Close()
	}
	require.Error(t, err)
	assert.Less(t, time.Since(start), 3*time.Second)
}
