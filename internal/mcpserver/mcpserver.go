// Package mcpserver is the MCP endpoint at `/mcp` (ADR-0020): a stateless
// Streamable HTTP server from the official SDK, whose every request carries an
// OAuth bearer token issued by internal/oauthserver (ADR-0021).
//
// It serves tools that read the catalogue and the signed-in learner's own progress. A tool does not reimplement what the
// website shows: it calls the same handler code as the HTTP operation it mirrors
// (Catalogue is satisfied by handlers.Server), and describes its result with the
// operation's schema from the contract.
//
// The package does not import internal/handlers, which imports it to mount the
// endpoint, the way it mounts the OAuth server.
package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/modelcontextprotocol/go-sdk/oauthex"

	"hexletbasics/internal/api"
	"hexletbasics/internal/oauthserver"
)

const (
	// Path is where the endpoint lives. It is also the path of the canonical
	// resource URL, which oauthserver builds from the issuer (ADR-0021).
	Path = "/mcp"
	// metadataPath is the RFC 9728 well-known location of the resource's metadata,
	// and metadataPathForResource the one for a resource that has a path: the
	// well-known segment is inserted before it. A client may try either, so both
	// answer.
	metadataPath            = "/.well-known/oauth-protected-resource"
	metadataPathForResource = metadataPath + Path
)

// Patterns are the ServeMux patterns the router mounts Handler at. They carry no
// method: the SDK answers an unsupported one itself, and the metadata handler
// also answers the CORS preflight a browser-based client sends.
func Patterns() []string {
	return []string{Path, metadataPath, metadataPathForResource}
}

// Verifier is what the endpoint needs of the OAuth server: to accept a bearer
// token, and to know the canonical URL its tokens must be bound to.
type Verifier interface {
	VerifyAccessToken(ctx context.Context, bearer string) (oauthserver.Token, error)
	Resource() string
}

// Catalogue is the read side of the website's public catalogue, as the HTTP
// operations implement it. Tools call it rather than the database, so a tool and
// an HTTP read cannot answer differently.
//
// The methods are the generated handler's own, with its own result types, so
// handlers.Server satisfies it as it is and a tool shows exactly what the
// operation answers a visitor.
type Catalogue interface {
	ListCourses(ctx context.Context) ([]api.CourseCatalogItem, error)
	ListPublicCourseCategories(ctx context.Context) ([]api.CourseCategory, error)
	GetPublicCourseCategory(ctx context.Context, params api.GetPublicCourseCategoryParams) (api.GetPublicCourseCategoryRes, error)
	GetCourse(ctx context.Context, params api.GetCourseParams) (api.GetCourseRes, error)
	GetCourseLesson(ctx context.Context, params api.GetCourseLessonParams) (api.GetCourseLessonRes, error)
	ListBlogPosts(ctx context.Context, params api.ListBlogPostsParams) (*api.BlogPostPage, error)
	GetBlogPost(ctx context.Context, params api.GetBlogPostParams) (api.GetBlogPostRes, error)
}

// Token info keys the verifier fills and the admin tools read into an Audit
// Record's Source.
const (
	extraConnectedAppID = "connected_app_id"
	extraTokenID        = "token_id"
)

// Deps are the collaborators of the Server.
type Deps struct {
	Verifier  Verifier
	Catalogue Catalogue
	// API serves an admin tool call: the generated api.Server, behind the
	// locale middleware, so a call runs the code of the HTTP operation it mirrors,
	// security handler and audit middleware included (ADR-0020 addendum). It is an
	// http.Handler on purpose: it gives the package no way to reach an operation
	// without the checkpoint. Without it no admin tool is offered.
	API http.Handler
	// Learner serves the learner tools; the token's user is the only subject.
	Learner Learner
	// Issuer is the authorization server's origin, advertised in the resource
	// metadata so a client that was refused knows where to sign in.
	Issuer string
	// Logger receives verification failures that are ours, not the caller's.
	Logger *slog.Logger
	// Limits are the page, response, timeout and rate limits of the tools; a
	// zero field takes its default.
	Limits Limits
}

// Server is the MCP endpoint.
type Server struct {
	deps     Deps
	resource string
	limits   Limits
	meter    *meter
	handler  http.Handler
}

// New builds the endpoint. It fails when a tool's output schema cannot be built
// from the contract, so a contract change that breaks a tool stops the process
// from starting instead of surfacing at the first call.
func New(deps Deps) (*Server, error) {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	limits := deps.Limits.withDefaults()
	s := &Server{deps: deps, resource: deps.Verifier.Resource(), limits: limits, meter: newMeter(limits)}

	mcpServer := mcp.NewServer(&mcp.Implementation{Name: "hexlet-basics", Title: "Hexlet Basics"}, nil)
	if err := s.addCatalogueTools(mcpServer); err != nil {
		return nil, err
	}
	if err := s.addLearnerTools(mcpServer); err != nil {
		return nil, err
	}
	if deps.API != nil {
		if err := s.addAdminTools(mcpServer); err != nil {
			return nil, err
		}
	}

	// Stateless: one request, one response, no sessions and no long-lived
	// streams, so the load is that of any other API call (ADR-0020). JSON
	// responses, because nothing here streams.
	streamable := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return mcpServer },
		&mcp.StreamableHTTPOptions{
			Stateless: true, JSONResponse: true,
			// Large enough for the upload tool's base64 argument (see upload.go).
			MaxRequestBodyBytes: requestBodyLimit(deps.API != nil && slices.Contains(adminOperations, "adminUploadAttachment")),
		},
	)
	metadataURL := deps.Issuer + metadataPathForResource
	protected := auth.RequireBearerToken(s.verify, &auth.RequireBearerTokenOptions{
		ResourceMetadataURL: metadataURL,
		// The verifier has already refused an expired token by looking it up
		// (ADR-0021), and the library does not carry an expiry.
		AllowMissingExpiration: true,
	})(streamable)

	resourceMetadata := auth.ProtectedResourceMetadataHandler(&oauthex.ProtectedResourceMetadata{
		Resource:               s.resource,
		AuthorizationServers:   []string{deps.Issuer},
		ScopesSupported:        []string{oauthserver.ScopeCatalogRead, oauthserver.ScopeProgressRead, oauthserver.ScopeAdmin},
		BearerMethodsSupported: []string{"header"},
		ResourceName:           "Hexlet Basics",
	})

	mux := http.NewServeMux()
	mux.Handle(Path, protected)
	mux.Handle(metadataPath, resourceMetadata)
	mux.Handle(metadataPathForResource, resourceMetadata)
	s.handler = mux
	return s, nil
}

// Handler serves the endpoint and its metadata, for the router to mount at Patterns.
func (s *Server) Handler() http.Handler { return s.handler }

// verify accepts a bearer token only when it is ours, live, bound to this
// endpoint and issued to a user.
//
// VerifyAccessToken decrypts and looks the token up but does not judge the
// audience (it has one resource and says what it is), so this is where a token
// meant for another resource is refused. A subject that is not a user (a client)
// is refused too: nothing here may act as anything but a person (ADR-0021).
// Every refusal is the same ErrInvalidToken, which answers 401 with the
// challenge; only a failure of ours becomes a 500, and its detail stays in the
// log.
func (s *Server) verify(ctx context.Context, bearer string, _ *http.Request) (*auth.TokenInfo, error) {
	token, err := s.deps.Verifier.VerifyAccessToken(ctx, bearer)
	if err != nil {
		if errors.Is(err, oauthserver.ErrInvalidToken) {
			return nil, auth.ErrInvalidToken
		}
		s.deps.Logger.ErrorContext(ctx, "verify MCP access token", slog.Any("error", err))
		return nil, errors.New("token verification is unavailable")
	}
	if !slices.Contains(token.Audience, s.resource) || token.UserID <= 0 {
		return nil, auth.ErrInvalidToken
	}
	return &auth.TokenInfo{
		Scopes: token.Scopes,
		UserID: strconv.Itoa(token.UserID),
		Extra: map[string]any{
			"client_id":         token.ClientID,
			extraConnectedAppID: token.ConnectedAppID,
			extraTokenID:        token.ID,
		},
	}, nil
}

// requireScope is the checkpoint of a catalogue tool: the connection must hold
// the Scope, and the token's subject must be a user (ADR-0020). The endpoint
// middleware has checked the subject already; it is read again from the context
// of the call so a tool cannot run without one.
func requireScope(req *mcp.CallToolRequest, scope string) (userID int, denied *mcp.CallToolResult) {
	info := req.GetExtra().TokenInfo
	if info == nil {
		return 0, toolError("this connection is not signed in")
	}
	userID, err := strconv.Atoi(info.UserID)
	if err != nil || userID <= 0 {
		return 0, toolError("this connection is not signed in as a user")
	}
	if !slices.Contains(info.Scopes, scope) {
		return 0, toolError(fmt.Sprintf("this connection was not granted the %s Scope", scope))
	}
	return userID, nil
}

// toolError is a tool result that tells the agent what went wrong in words, the
// way MCP reports a failed call (a protocol error would end the exchange).
func toolError(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: message}}}
}

// summarize is the short text of a list result.
func summarize(noun string, names []string) string {
	return fmt.Sprintf("%d %s: %s", len(names), noun, strings.Join(names, ", "))
}
