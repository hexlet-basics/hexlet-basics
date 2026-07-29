// Package localization owns backend-emitted, human-readable text.
//
// UI copy remains in React i18next and localized course content remains in the
// database. This package only translates strings produced by the Go backend.
package localization

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFiles embed.FS

type (
	localizerContextKey struct{}
	localeContextKey    struct{}
)

// Message is a backend message registered with the localization module.
//
// Its fields are deliberately private: callers choose one of the messages
// below instead of inventing ad-hoc IDs or fallback text in a handler.
type Message struct {
	value i18n.Message
}

// Translator owns the process-wide translation bundle. A request-specific
// go-i18n localizer is derived by Middleware and carried through context.
type Translator struct {
	bundle  *i18n.Bundle
	matcher language.Matcher
}

// New loads and validates every embedded locale catalog.
func New() (*Translator, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, name := range []string{
		"locales/active.en.json",
		"locales/active.ru.json",
		"locales/active.es.json",
	} {
		if _, err := bundle.LoadMessageFileFS(localeFiles, name); err != nil {
			return nil, fmt.Errorf("load localization catalog %s: %w", name, err)
		}
	}

	return &Translator{
		bundle:  bundle,
		matcher: language.NewMatcher(bundle.LanguageTags()),
	}, nil
}

// Middleware resolves Accept-Language once per request and makes the resulting
// localizer available to generated and hand-written HTTP handlers alike.
func (t *Translator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept-Language")
		header := r.Header.Get("Accept-Language")
		localizer := i18n.NewLocalizer(t.bundle, header)
		ctx := context.WithValue(r.Context(), localizerContextKey{}, localizer)
		tags, _, err := language.ParseAcceptLanguage(header)
		if err != nil || len(tags) == 0 {
			tags = []language.Tag{language.English}
		}
		tag, _, _ := t.matcher.Match(tags...)
		base, _ := tag.Base()
		ctx = context.WithValue(ctx, localeContextKey{}, base.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Locale returns the canonical request locale used by domain facts. Contexts
// outside HTTP intentionally use English, matching Text's fallback.
func (t *Translator) Locale(ctx context.Context) string {
	if locale, ok := ctx.Value(localeContextKey{}).(string); ok {
		return locale
	}
	return language.English.String()
}

// Text translates message for the locale attached to ctx. Contexts that did
// not pass through Middleware intentionally fall back to English, which keeps
// direct handler tests and non-HTTP callers deterministic.
func (t *Translator) Text(ctx context.Context, message Message) string {
	localizer, ok := ctx.Value(localizerContextKey{}).(*i18n.Localizer)
	if !ok {
		localizer = i18n.NewLocalizer(t.bundle, language.English.String())
	}

	text, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:      message.value.ID,
		DefaultMessage: &message.value,
	})
	if err != nil {
		return message.value.Other
	}
	return text
}

// StatusText returns the localized generic text for an HTTP status. Unknown
// statuses retain net/http's English text rather than exposing internal errors.
func (t *Translator) StatusText(ctx context.Context, status int) string {
	if message, ok := statusMessages[status]; ok {
		return t.Text(ctx, message)
	}
	if text := http.StatusText(status); text != "" {
		return text
	}
	return t.Text(ctx, InternalServerError)
}
