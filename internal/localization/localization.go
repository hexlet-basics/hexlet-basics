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
	"io/fs"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFiles embed.FS

type localizerContextKey struct{}

// Message is a backend message registered with the localization module.
//
// Its fields are deliberately private: callers choose one of the messages
// below instead of inventing ad-hoc IDs or fallback text in a handler.
type Message struct {
	id      string
	english string
}

// Translator owns the process-wide translation bundle. A request-specific
// go-i18n localizer is derived by Middleware and carried through context.
type Translator struct {
	bundle *i18n.Bundle
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
		if err := validateCatalog(localeFiles, name); err != nil {
			return nil, err
		}
		if _, err := bundle.LoadMessageFileFS(localeFiles, name); err != nil {
			return nil, fmt.Errorf("load localization catalog %s: %w", name, err)
		}
	}

	return &Translator{bundle: bundle}, nil
}

// Middleware resolves Accept-Language once per request and makes the resulting
// localizer available to generated and hand-written HTTP handlers alike.
func (t *Translator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept-Language")
		localizer := i18n.NewLocalizer(t.bundle, r.Header.Get("Accept-Language"))
		ctx := context.WithValue(r.Context(), localizerContextKey{}, localizer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
		MessageID: message.id,
		DefaultMessage: &i18n.Message{
			ID:    message.id,
			Other: message.english,
		},
	})
	if err != nil {
		return message.english
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

type catalogMessage struct {
	Other string `json:"other"`
}

func validateCatalog(files fs.FS, name string) error {
	data, err := fs.ReadFile(files, name)
	if err != nil {
		return fmt.Errorf("read localization catalog %s: %w", name, err)
	}

	var catalog map[string]catalogMessage
	if err := json.Unmarshal(data, &catalog); err != nil {
		return fmt.Errorf("decode localization catalog %s: %w", name, err)
	}
	for _, message := range allMessages {
		entry, ok := catalog[message.id]
		if !ok || entry.Other == "" {
			return fmt.Errorf("localization catalog %s is missing %q", name, message.id)
		}
	}
	return nil
}
