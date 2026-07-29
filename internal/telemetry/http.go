package telemetry

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

// NewSentryHTTPHandler adds Sentry's HTTP integration around the complete
// application handler. Each request starts with a clone of the process client's
// hub so request data and scope mutations cannot leak between concurrent calls.
//
// Repanicking preserves net/http's normal panic behavior after Sentry records
// the failure; swallowing the panic could otherwise turn an unwritten response
// into an empty 200.
func NewSentryHTTPHandler(client *sentry.Client, next http.Handler) http.Handler {
	rootHub := sentry.NewHub(client, sentry.NewScope())
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	}).Handle(next)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub := rootHub.Clone()
		ctx := sentry.SetHubOnContext(r.Context(), hub)
		sentryHandler.ServeHTTP(w, r.WithContext(ctx))
	})
}
