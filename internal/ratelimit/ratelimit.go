// Package ratelimit is a fixed-window counter keyed by an opaque string. The
// OAuth server counts client registrations per address with it (ADR-0019) and
// the MCP endpoint counts tool calls per user and per Connected App (ADR-0018).
//
// It lives in memory, per replica: with several replicas a key gets its limit
// once per replica. That is a deliberate trade for having no table to clean and
// no new dependency; the limit is a brake on abuse and a runaway agent, not an
// exact quota.
package ratelimit

import (
	"sync"
	"time"
)

// maxTracked bounds the limiter's memory. A flood of this many distinct keys
// inside one window makes new keys wait for the window to pass rather than
// letting the map grow without limit.
const maxTracked = 100_000

// Charge is one attempt against one key. The limit travels with it, so one
// Limiter (one window, one map) can hold counters with different budgets, such
// as the per-user and per-client budgets of each call class.
type Charge struct {
	Key   string
	Limit int
}

// Limiter is a fixed-window counter per key. It is in memory and per replica.
type Limiter struct {
	window time.Duration
	now    func() time.Time

	mu      sync.Mutex
	windows map[string]*keyWindow
}

type keyWindow struct {
	start time.Time
	count int
}

// New builds a limiter whose counters reset window after their first attempt.
// now is the clock, injected so a test can move time without sleeping.
func New(window time.Duration, now func() time.Time) *Limiter {
	return &Limiter{window: window, now: now, windows: map[string]*keyWindow{}}
}

// Allow spends one attempt of every charge, or none: the attempt is refused as
// a whole when any key is out of budget, so a caller refused by one dimension
// (the Connected App) does not also lose budget on another (the user). When it
// refuses, wait is how long until the key that is out of budget has it back
// (the longest, if several are).
func (l *Limiter) Allow(charges ...Charge) (ok bool, wait time.Duration) {
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()

	current := make([]*keyWindow, len(charges))
	for i, c := range charges {
		w := l.windows[c.Key]
		if w != nil && now.Sub(w.start) >= l.window {
			delete(l.windows, c.Key)
			w = nil
		}
		current[i] = w
		switch {
		case w != nil && w.count >= c.Limit:
			wait = max(wait, w.start.Add(l.window).Sub(now))
		case w == nil && len(l.windows) >= maxTracked:
			l.pruneLocked(now)
			if len(l.windows) >= maxTracked {
				wait = max(wait, l.window)
			}
		}
	}
	if wait > 0 {
		return false, wait
	}
	for i, c := range charges {
		if current[i] == nil {
			current[i] = &keyWindow{start: now}
			l.windows[c.Key] = current[i]
		}
		current[i].count++
	}
	return true, 0
}

func (l *Limiter) pruneLocked(now time.Time) {
	for key, w := range l.windows {
		if now.Sub(w.start) >= l.window {
			delete(l.windows, key)
		}
	}
}
