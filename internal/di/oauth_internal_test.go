package di

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func serveRouter(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()

	injector := newServerContainer()

	t.Cleanup(func() { assert.Empty(t, injector.Shutdown().Errors) })

	router, err := do.InvokeNamed[http.Handler](injector, routerServiceName)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil))

	return rec
}

// The OAuth server ships dark: with the flag at its default the router has no
// /oauth surface and no well-known document, and nothing asks for an issuer.
func TestOAuthServerIsNotMountedByDefault(t *testing.T) {
	configureTestEnvironment(t)

	for _, path := range []string{
		"/.well-known/oauth-authorization-server",
		"/.well-known/openid-configuration",
		"/oauth/token",
	} {
		assert.Equal(t, http.StatusNotFound, serveRouter(t, path).Code, path)
	}
}

func TestOAuthServerIsMountedBehindItsFlag(t *testing.T) {
	configureTestEnvironment(t)
	t.Setenv("OAUTH_ENABLED", "true")
	t.Setenv("OAUTH_ISSUER", "https://code-basics.test")
	t.Setenv("OAUTH_ENCRYPTION_KEY", "test-oauth-key-test-oauth-key-test")

	rec := serveRouter(t, "/.well-known/oauth-authorization-server")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"issuer":"https://code-basics.test"`)
	assert.Contains(t, rec.Body.String(), `"token_endpoint":"https://code-basics.test/oauth/token"`)
}
