package handlers_test

import (
	"context"
	"net/http"
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
