package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// With the OAuth server switched off nothing is wired behind the consent
// operations, and they answer as if the request had never existed: a deployment
// that has not enabled the feature exposes no way to probe it.
func TestConsentOperationsAnswerNotFoundWhileTheOAuthServerIsOff(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	got, err := h.Client.GetOAuthConsent(ctx, api.GetOAuthConsentParams{RequestId: "any"})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, got)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())

	approved, err := h.Client.ApproveOAuthConsent(ctx, api.ApproveOAuthConsentParams{RequestId: "any"})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, approved)

	denied, err := h.Client.DenyOAuthConsent(ctx, api.DenyOAuthConsentParams{RequestId: "any"})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, denied)
}

// OAUTH_ENABLED and MCP_ENABLED default to off, and the ingress sends their
// paths to this server regardless (the chart carries no per-flag routing). With
// both flags off the router mounts nothing for them, so each path must fall
// through to the generated server's 404 rather than reach a half-wired handler.
func TestOAuthAndMCPPathsAnswerNotFoundWhileTheirFlagsAreOff(t *testing.T) {
	router, _, _ := newAttachmentRouterStack(t, false)

	for _, target := range []string{
		"/mcp",
		"/oauth/token",
		"/oauth/register",
		"/.well-known/oauth-authorization-server",
		"/.well-known/openid-configuration",
		"/.well-known/oauth-protected-resource",
		"/.well-known/oauth-protected-resource/mcp",
	} {
		t.Run(target, func(t *testing.T) {
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, target, nil))

			assert.Equal(t, http.StatusNotFound, rec.Code)
		})
	}
}

// Without the OAuth server there can be no Connected App: the list is empty (the
// account page shows its empty state) and a revoke finds nothing.
func TestConnectedAppsAreEmptyWhileTheOAuthServerIsOff(t *testing.T) {
	h := testsupport.NewHarnessAs(t, "user_admin_learner")
	ctx := context.Background()

	list, err := h.Client.ListConnectedApps(ctx)
	require.NoError(t, err)
	apps, ok := list.(*api.ConnectedApps)
	require.True(t, ok, "got %T", list)
	assert.NotNil(t, apps.Items)
	assert.Empty(t, apps.Items)

	revoked, err := h.Client.RevokeConnectedApp(ctx, api.RevokeConnectedAppParams{ID: 1})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, revoked)
}
