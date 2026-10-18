package mcpserver

import (
	"fmt"
	"strconv"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"hexletbasics/internal/ratelimit"
)

// Class is the kind of work a tool does, which decides how often one caller may
// ask for it (ADR-0021): reading the public catalogue is cheap and harmless,
// reading a learner's own data is more sensitive, and a write changes the site.
// A tool names its Class when it is registered, so the tickets that add learner
// and admin tools only pick one.
type Class int

const (
	// ClassCatalogue is a read of data a signed-out visitor sees on the website.
	ClassCatalogue Class = iota
	// ClassPersonal is a read of the signed-in user's own data (progress,
	// enrollment): stricter than catalogue reads.
	ClassPersonal
	// ClassWrite is any tool that changes data: the strictest.
	ClassWrite
)

func (c Class) String() string {
	switch c {
	case ClassCatalogue:
		return "catalogue"
	case ClassPersonal:
		return "personal"
	case ClassWrite:
		return "write"
	default:
		return "class" + strconv.Itoa(int(c))
	}
}

// Rates is how many calls of each Class one subject may make per window.
type Rates struct {
	Catalogue, Personal, Write int
}

func (r Rates) of(c Class) int {
	switch c {
	case ClassPersonal:
		return r.Personal
	case ClassWrite:
		return r.Write
	default:
		return r.Catalogue
	}
}

// Limits are what protects the server from one assistant (ADR-0021). The zero
// value of a field means its default, so a caller sets only what it changes.
type Limits struct {
	// MaxPageSize is the most records one result holds.
	MaxPageSize int
	// MaxResponseBytes is the most one result's structured content and text hold.
	MaxResponseBytes int
	// ToolTimeout is how long a call may run before its context is cancelled.
	ToolTimeout time.Duration
	// RateWindow is the period the rates count over.
	RateWindow time.Duration
	// User is the budget of one user, Client the budget of one Connected App (an
	// OAuth client, shared by every user who connected it). A call spends one
	// call of both.
	User, Client Rates
	// Now is the clock of the rate limiter; a test moves it.
	Now func() time.Time
}

// DefaultLimits are the values config applies when nothing is set; they are
// repeated in internal/config's env defaults, which production reads.
func DefaultLimits() Limits {
	return Limits{
		MaxPageSize:      50,
		MaxResponseBytes: 100 * 1024,
		ToolTimeout:      30 * time.Second,
		RateWindow:       time.Minute,
		User:             Rates{Catalogue: 60, Personal: 30, Write: 10},
		Client:           Rates{Catalogue: 600, Personal: 300, Write: 100},
		Now:              time.Now,
	}
}

func (l Limits) withDefaults() Limits {
	d := DefaultLimits()
	fill := func(v *int, def int) {
		if *v <= 0 {
			*v = def
		}
	}
	fill(&l.MaxPageSize, d.MaxPageSize)
	fill(&l.MaxResponseBytes, d.MaxResponseBytes)
	fill(&l.User.Catalogue, d.User.Catalogue)
	fill(&l.User.Personal, d.User.Personal)
	fill(&l.User.Write, d.User.Write)
	fill(&l.Client.Catalogue, d.Client.Catalogue)
	fill(&l.Client.Personal, d.Client.Personal)
	fill(&l.Client.Write, d.Client.Write)
	if l.ToolTimeout <= 0 {
		l.ToolTimeout = d.ToolTimeout
	}
	if l.RateWindow <= 0 {
		l.RateWindow = d.RateWindow
	}
	if l.Now == nil {
		l.Now = d.Now
	}
	return l
}

// meter counts calls per user and per Connected App, per Class, on the limiter
// shared with the OAuth server's registration limit. It is in memory and per
// replica: with several replicas a caller gets its budget once per replica, a
// brake on a runaway agent rather than an exact quota.
type meter struct {
	limits  Limits
	counter *ratelimit.Limiter
}

func newMeter(limits Limits) *meter {
	return &meter{limits: limits, counter: ratelimit.New(limits.RateWindow, limits.Now)}
}

// allow spends one call of the Class from the user's budget and from the
// Connected App's, or neither: a call refused because the app is out of budget
// does not also cost the user one. When it refuses, wait is the time until the
// caller may try again.
func (m *meter) allow(class Class, userID int, clientID string) (bool, time.Duration) {
	return m.counter.Allow(
		ratelimit.Charge{Key: class.String() + "|user|" + strconv.Itoa(userID), Limit: m.limits.User.of(class)},
		ratelimit.Charge{Key: class.String() + "|client|" + clientID, Limit: m.limits.Client.of(class)},
	)
}

// admit is the checkpoint every tool passes before it does any work: the
// connection must hold the Scope and be a user (requireScope), and the caller
// must be within the rate limit of the tool's Class. A call refused for its
// Scope spends no budget, so a connection that was never allowed to use a tool
// cannot use up what the same user's other connections may call.
//
// The refusal is a tool error rather than HTTP 429 (decision in ADR-0021): the
// Class is known only from the tool name inside the JSON-RPC body, and a 429 on
// the endpoint would also stop `initialize` and `tools/list`, which cost
// nothing, and an SDK client reads it as a broken transport, not as a message
// to the agent. The error names the wait, in the text and in `_meta`.
func (s *Server) admit(req *mcp.CallToolRequest, scope string, class Class) (int, string, *mcp.CallToolResult) {
	userID, denied := requireScope(req, scope)
	if denied != nil {
		return 0, "", denied
	}
	clientID, denied := s.spend(req, userID, class)
	return userID, clientID, denied
}

// spend counts one call of the Class against the user and the Connected App and
// refuses it when either budget is gone. Admin tools call it without admit's
// Scope check, because their Scope is judged by the generated server's security
// handler, which records a refusal.
func (s *Server) spend(req *mcp.CallToolRequest, userID int, class Class) (string, *mcp.CallToolResult) {
	var clientID string
	if client, ok := req.GetExtra().TokenInfo.Extra["client_id"].(string); ok {
		clientID = client
	}
	if ok, wait := s.meter.allow(class, userID, clientID); !ok {
		seconds := int(wait.Round(time.Second) / time.Second)
		seconds = max(seconds, 1)
		refusal := toolError(fmt.Sprintf(
			"rate limit exceeded for %s tools: retry in %d seconds (not before %s)",
			class, seconds, s.limits.Now().Add(wait).UTC().Format(time.RFC3339)))
		refusal.Meta = mcp.Meta{"retryAfterSeconds": seconds, "rateLimitClass": class.String()}
		return "", refusal
	}
	return clientID, nil
}
