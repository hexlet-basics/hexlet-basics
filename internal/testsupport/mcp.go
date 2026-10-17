package testsupport

import (
	"net/http"
	"net/http/httptest"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

// MCPEndpoint is the canonical MCP URL of the test server: the resource every
// token the OAuth helper issues is bound to.
const MCPEndpoint = OAuthIssuer + "/mcp"

// This file is the top-of-stack test seam for everything MCP (ADR-0018): the
// full production router in process, driven by a real MCP client from the
// official SDK, authenticated with a token from the OAuth helper. A test reads
// like a connected assistant: it never builds a JSON-RPC message or calls a
// handler directly, so what it proves is what a client would see, including the
// token check, the Scope check and the contract-derived schemas.
//
// How to use it from a tool ticket (catalogue, learner, admin):
//
//	o := testsupport.NewOAuth(t)
//	tokens := o.Connect(userID, oauthserver.ScopeCatalogRead)
//	session := o.MCPClient(tokens.AccessToken)
//	result, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "listCourses"})
//
// Connect runs the real authorization-code flow for a fixture user, so the token
// carries that user as its subject and exactly the Scopes asked for; ask for
// fewer to test a refusal. Use MCPTransport instead of MCPClient when a test is
// about a request the SDK client would not send (no token, a cookie session, a
// token for another audience): it is an http.Client whose transport serves the
// request through the router without a socket, like the other helpers here.
// Database state is the test's rolled-back transaction (o.DB), as everywhere.
//
// Only what NewOAuth wires is reachable: the router is built with the handler
// dependencies the OAuth flow needs, so an operation or tool that reads
// something it leaves out (the dashboard, for one) panics on a nil dependency.
// A ticket whose tool needs more adds it to the handlers.Deps in NewOAuth.

// MCPTransport is an http.Client that serves requests through the router in
// process. When bearer is not empty it sends it as the Authorization header of
// every request; mutate edits each request after that (an Accept-Language, say).
func (o *OAuth) MCPTransport(bearer string, mutate ...func(*http.Request)) *http.Client {
	return &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		r = r.Clone(r.Context())
		if bearer != "" {
			r.Header.Set("Authorization", "Bearer "+bearer)
		}
		for _, m := range mutate {
			m(r)
		}
		rec := httptest.NewRecorder()
		o.Router.ServeHTTP(rec, r)
		res := rec.Result()
		res.Request = r
		return res, nil
	})}
}

// MCPClient connects a real SDK client to the endpoint with the bearer token and
// returns the session, closed when the test ends. It fails the test when the
// connection is refused; a test of a refusal uses MCPTransport or MCPConnect.
func (o *OAuth) MCPClient(bearer string, mutate ...func(*http.Request)) *mcp.ClientSession {
	o.t.Helper()
	session, err := o.MCPConnect(bearer, mutate...)
	require.NoError(o.t, err)
	return session
}

// MCPConnect is MCPClient returning the error, for a test that expects the
// connection to fail.
func (o *OAuth) MCPConnect(bearer string, mutate ...func(*http.Request)) (*mcp.ClientSession, error) {
	o.t.Helper()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0"}, nil)
	session, err := client.Connect(o.t.Context(), &mcp.StreamableClientTransport{
		Endpoint: MCPEndpoint,
		// The server is stateless and offers no stream to listen on.
		DisableStandaloneSSE: true,
		HTTPClient:           o.MCPTransport(bearer, mutate...),
		MaxRetries:           -1,
	}, nil)
	if err != nil {
		return nil, err
	}
	o.t.Cleanup(func() { _ = session.Close() })
	return session, nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
