package di

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The MCP endpoint ships dark, like the OAuth server it depends on: with its
// flag at the default there is nothing at /mcp or its metadata paths.
func TestMCPEndpointIsNotMountedByDefault(t *testing.T) {
	configureTestEnvironment(t)

	for _, path := range []string{
		"/.well-known/oauth-protected-resource",
		"/.well-known/oauth-protected-resource/mcp",
		"/mcp",
	} {
		assert.Equal(t, http.StatusNotFound, serveRouter(t, path).Code, path)
	}
}

func TestMCPEndpointIsMountedBehindItsFlag(t *testing.T) {
	configureTestEnvironment(t)
	t.Setenv("OAUTH_ENABLED", "true")
	t.Setenv("OAUTH_ISSUER", "https://code-basics.test")
	t.Setenv("OAUTH_ENCRYPTION_KEY", "test-oauth-key-test-oauth-key-test")
	t.Setenv("MCP_ENABLED", "true")

	rec := serveRouter(t, "/.well-known/oauth-protected-resource/mcp")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"resource":"https://code-basics.test/mcp"`)
	// The endpoint itself answers a request with no token by pointing at it.
	unauthenticated := serveRouter(t, "/mcp")
	assert.Equal(t, http.StatusUnauthorized, unauthenticated.Code)
	assert.Contains(t, unauthenticated.Header().Get("WWW-Authenticate"), "resource_metadata=")
}
