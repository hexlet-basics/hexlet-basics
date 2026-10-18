// Package health serves the process probes Kubernetes polls: /healthz (the
// process is alive) and /readyz (it can take traffic). They are plain HTTP on a
// mux around the generated handler, outside the OpenAPI contract and not under
// /api (ADR-0023).
package health

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"
)

// readyTimeout bounds one readiness check, so a hung database turns into a 503
// before the probe's own timeout (3 s in the chart) gives up on the pod.
const readyTimeout = 2 * time.Second

// Probes answers /healthz and /readyz. Check reports whether the process's
// dependencies work; Drain marks the process as leaving.
type Probes struct {
	check    func(context.Context) error
	draining atomic.Bool
}

// New returns Probes whose readiness is the result of check.
func New(check func(context.Context) error) *Probes {
	return &Probes{check: check}
}

// Drain makes /readyz answer 503 from now on. The server calls it on SIGTERM,
// before the HTTP drain begins, so the endpoint is taken out of rotation while
// in-flight requests finish. /healthz is unaffected: a draining process is
// still alive and must not be restarted.
func (p *Probes) Drain() { p.draining.Store(true) }

// Wrap serves the probes and hands every other request to next.
func (p *Probes) Wrap(next http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /readyz", p.ready)
	mux.Handle("/", next)
	return mux
}

func (p *Probes) ready(w http.ResponseWriter, r *http.Request) {
	if p.draining.Load() {
		http.Error(w, "shutting down", http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), readyTimeout)
	defer cancel()
	if err := p.check(ctx); err != nil {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
