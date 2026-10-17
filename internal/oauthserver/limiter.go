package oauthserver

import (
	"net/netip"
	"time"

	"hexletbasics/internal/ratelimit"
)

const (
	// DefaultRegistrationLimit and DefaultRegistrationWindow are how many
	// registrations (successful or not) one address may attempt per window: ten
	// an hour is far more than a person connecting assistants needs, and far too
	// few to fill the table from one machine.
	DefaultRegistrationLimit  = 10
	DefaultRegistrationWindow = time.Hour
)

// limiter counts registration attempts per client address, on the fixed-window
// counter shared with the MCP endpoint (internal/ratelimit). It lives in
// memory, so with several replicas an address gets the limit once per replica;
// that is a deliberate trade for having no table to clean and no new
// dependency, and the limit is a brake on abuse, not an exact quota.
type limiter struct {
	limit   int
	counter *ratelimit.Limiter
}

func newLimiter(limit int, window time.Duration, now func() time.Time) *limiter {
	if limit <= 0 {
		limit = DefaultRegistrationLimit
	}
	if window <= 0 {
		window = DefaultRegistrationWindow
	}
	return &limiter{limit: limit, counter: ratelimit.New(window, now)}
}

// allow spends one attempt of an address. When it refuses, wait is how long
// until the address has its budget back.
func (l *limiter) allow(address string) (ok bool, wait time.Duration) {
	return l.counter.Allow(ratelimit.Charge{Key: addressKey(address), Limit: l.limit})
}

// addressKey names the bucket of an address. An IPv6 customer is handed a whole
// /64, so counting single addresses would let one machine mint a new budget per
// request; the prefix is the unit instead. Anything that is not an IP shares one
// bucket, which fails closed: a request whose address could not be resolved
// competes with the others like that instead of being exempt.
func addressKey(address string) string {
	ip, err := netip.ParseAddr(address)
	if err != nil {
		return "unknown"
	}
	ip = ip.Unmap()
	if ip.Is6() {
		if p, err := ip.Prefix(64); err == nil {
			return p.String()
		}
	}
	return ip.String()
}
