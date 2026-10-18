package mcpserver_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// secondClient is a second Connected App.
const secondClient = "second-test-mcp-client"

// rejected asserts a call was refused for the rate limit and returns the wait
// the refusal names.
func rejected(t *testing.T, result *mcp.CallToolResult) int {
	t.Helper()
	require.True(t, result.IsError)
	assert.Contains(t, textOf(t, result), "rate limit exceeded")
	assert.Contains(t, textOf(t, result), "retry in")
	assert.Nil(t, result.StructuredContent)
	seconds, ok := result.Meta["retryAfterSeconds"].(float64)
	require.True(t, ok, "the wait is machine readable too: %v", result.Meta)
	assert.Positive(t, seconds)

	return int(seconds)
}

func accepted(t *testing.T, session *mcp.ClientSession) {
	t.Helper()
	result := callTool(t, session, "listCourses", nil)
	require.False(t, result.IsError, "%+v", result.Content)
}

func TestCallerOverTheUserLimitIsToldWhenToRetryAndGetsItsBudgetBack(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{
		RateWindow: 90 * time.Second,
		User:       mcpserver.Rates{Catalogue: 2},
	}))
	session := connect(t, o, alice(), oauthserver.ScopeCatalogRead)

	accepted(t, session)
	accepted(t, session)
	wait := rejected(t, callTool(t, session, "listCourses", nil))
	assert.Equal(t, 90, wait)

	// Another tool of the same class draws on the same budget.
	rejected(t, callTool(t, session, "listBlogPosts", nil))

	// Only the tools are limited: the session itself still works.
	_, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	o.Advance(91 * time.Second)
	accepted(t, session)
}

func TestOneUsersBudgetIsNotAnotherUsers(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{User: mcpserver.Rates{Catalogue: 1}}))
	first := connect(t, o, alice(), oauthserver.ScopeCatalogRead)
	second := o.MCPClient(o.ConnectAs(testsupport.OAuthClientID, bob(), oauthserver.ScopeCatalogRead).AccessToken, russian)

	accepted(t, first)
	rejected(t, callTool(t, first, "listCourses", nil))

	accepted(t, second)
	rejected(t, callTool(t, second, "listCourses", nil))
}

func TestOneConnectedAppsBudgetIsNotAnotherAppsAndIsSharedByItsUsers(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{Client: mcpserver.Rates{Catalogue: 2}}))
	// Alice and Bob through the app Connect registers, Alice also through another.
	aliceHere := connect(t, o, alice(), oauthserver.ScopeCatalogRead)
	bobHere := o.MCPClient(o.ConnectAs(testsupport.OAuthClientID, bob(), oauthserver.ScopeCatalogRead).AccessToken, russian)
	aliceThere := o.MCPClient(o.ConnectAs(secondClient, alice(), oauthserver.ScopeCatalogRead).AccessToken, russian)

	accepted(t, aliceHere)
	accepted(t, bobHere)
	rejected(t, callTool(t, aliceHere, "listCourses", nil))
	rejected(t, callTool(t, bobHere, "listCourses", nil))

	accepted(t, aliceThere)
	accepted(t, aliceThere)
	rejected(t, callTool(t, aliceThere, "listCourses", nil))
}

func TestCallRefusedForTheAppsLimitDoesNotSpendTheUsersBudget(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{
		User:   mcpserver.Rates{Catalogue: 2},
		Client: mcpserver.Rates{Catalogue: 1},
	}))
	here := connect(t, o, alice(), oauthserver.ScopeCatalogRead)
	there := o.MCPClient(o.ConnectAs(secondClient, alice(), oauthserver.ScopeCatalogRead).AccessToken, russian)

	accepted(t, here)
	rejected(t, callTool(t, here, "listCourses", nil))
	rejected(t, callTool(t, here, "listCourses", nil))

	// Alice has made one call and was refused twice: one call of hers is left.
	accepted(t, there)
	rejected(t, callTool(t, there, "listCourses", nil))
}

func TestCallRefusedForItsScopeSpendsNoBudget(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{User: mcpserver.Rates{Catalogue: 1}}))
	withoutScope := connect(t, o, alice(), oauthserver.ScopeProgressRead)
	withScope := o.MCPClient(o.ConnectAs(secondClient, alice(), oauthserver.ScopeCatalogRead).AccessToken, russian)

	for range 3 {
		assert.Contains(t, textOf(t, callTool(t, withoutScope, "listCourses", nil)), oauthserver.ScopeCatalogRead)
	}

	accepted(t, withScope)
}

// slowCatalogue blocks until its context ends, to show the timeout reaches the
// handler's context.
type slowCatalogue struct {
	mcpserver.Catalogue

	cancelled chan error
}

func (c slowCatalogue) ListCourses(ctx context.Context) ([]api.CourseCatalogItem, error) {
	<-ctx.Done()

	c.cancelled <- ctx.Err()

	return nil, ctx.Err()
}

func TestTimeoutCancelsTheHandlersContext(t *testing.T) {
	catalogue := slowCatalogue{cancelled: make(chan error, 1)}
	server, err := mcpserver.New(mcpserver.Deps{
		Verifier: verifier{token: oauthserver.Token{
			UserID: alice(), ClientID: "c", Scopes: []string{oauthserver.ScopeCatalogRead}, Audience: []string{testsupport.MCPEndpoint},
		}},
		Catalogue: catalogue,
		Issuer:    testsupport.OAuthIssuer,
		Limits:    mcpserver.Limits{ToolTimeout: 50 * time.Millisecond},
	})
	require.NoError(t, err)

	ts := httptest.NewServer(server.Handler())
	t.Cleanup(ts.Close)

	client := mcp.NewClient(&mcp.Implementation{Name: "t", Version: "0"}, nil)
	session, err := client.Connect(t.Context(), &mcp.StreamableClientTransport{
		Endpoint:             ts.URL + "/mcp",
		DisableStandaloneSSE: true,
		HTTPClient:           &http.Client{Transport: bearerTransport{}},
		MaxRetries:           -1,
	}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })

	result := callTool(t, session, "listCourses", nil)

	assert.True(t, result.IsError)
	assert.Contains(t, textOf(t, result), "took longer than 50ms")
	assert.ErrorIs(t, <-catalogue.cancelled, context.DeadlineExceeded, "the handler's context was cancelled by the timeout")
}

type bearerTransport struct{}

func (bearerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer anything")

	return http.DefaultTransport.RoundTrip(r)
}
