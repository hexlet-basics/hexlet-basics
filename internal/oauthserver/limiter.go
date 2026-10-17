package oauthserver

import (
	"net/netip"
	"sync"
	"time"
)

const (
	// DefaultRegistrationLimit and DefaultRegistrationWindow are how many
	// registrations (successful or not) one address may attempt per window: ten
	// an hour is far more than a person connecting assistants needs, and far too
	// few to fill the table from one machine.
	DefaultRegistrationLimit  = 10
	DefaultRegistrationWindow = time.Hour

	// maxTrackedAddresses bounds the limiter's memory. A flood from this many
	// distinct addresses inside one window makes new addresses wait for the
	// window to pass rather than letting the map grow without limit.
	maxTrackedAddresses = 100_000
)

// limiter is a fixed-window counter per client address. It lives in memory, so
// with several replicas an address gets the limit once per replica; that is a
// deliberate trade for having no table to clean and no new dependency, and the
// limit is a brake on abuse, not an exact quota.
type limiter struct {
	limit  int
	window time.Duration
	now    func() time.Time

	mu      sync.Mutex
	windows map[string]*addressWindow
}

type addressWindow struct {
	start time.Time
	count int
}

func newLimiter(limit int, window time.Duration, now func() time.Time) *limiter {
	if limit <= 0 {
		limit = DefaultRegistrationLimit
	}
	if window <= 0 {
		window = DefaultRegistrationWindow
	}
	return &limiter{limit: limit, window: window, now: now, windows: map[string]*addressWindow{}}
}

// allow spends one attempt of an address. When it refuses, wait is how long
// until the address has its budget back.
func (l *limiter) allow(address string) (ok bool, wait time.Duration) {
	key := addressKey(address)
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()

	w := l.windows[key]
	if w != nil && now.Sub(w.start) >= l.window {
		delete(l.windows, key)
		w = nil
	}
	if w == nil {
		if len(l.windows) >= maxTrackedAddresses {
			l.pruneLocked(now)
			if len(l.windows) >= maxTrackedAddresses {
				return false, l.window
			}
		}
		w = &addressWindow{start: now}
		l.windows[key] = w
	}
	if w.count >= l.limit {
		return false, w.start.Add(l.window).Sub(now)
	}
	w.count++
	return true, 0
}

func (l *limiter) pruneLocked(now time.Time) {
	for key, w := range l.windows {
		if now.Sub(w.start) >= l.window {
			delete(l.windows, key)
		}
	}
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
