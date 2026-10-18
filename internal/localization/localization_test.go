package localization_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/localization"
)

func TestTranslatorSelectsRequestLanguage(t *testing.T) {
	t.Parallel()
	translator, err := localization.New()
	require.NoError(t, err)

	tests := []struct {
		name     string
		header   string
		expected string
		locale   string
	}{
		{name: "missing header", expected: "Wrong email or password", locale: "en"},
		{name: "Russian region", header: "ru-RU", expected: "Неверный адрес электронной почты или пароль", locale: "ru"},
		{name: "Spanish", header: "es", expected: "Correo electrónico o contraseña incorrectos", locale: "es"},
		{name: "quality priority", header: "es;q=0.7, ru-RU;q=0.9", expected: "Неверный адрес электронной почты или пароль", locale: "ru"},
		{name: "unsupported", header: "de-DE", expected: "Wrong email or password", locale: "en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := translator.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.locale, translator.Locale(r.Context()))
				_, _ = w.Write([]byte(translator.Text(r.Context(), localization.WrongCredentials)))
			}))
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
			req.Header.Set("Accept-Language", tt.header)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expected, rec.Body.String())
			assert.Contains(t, rec.Header().Values("Vary"), "Accept-Language")
		})
	}
}

func TestTranslatorLocalizesHTTPStatus(t *testing.T) {
	t.Parallel()
	translator, err := localization.New()
	require.NoError(t, err)

	handler := translator.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(translator.StatusText(r.Context(), http.StatusNotFound)))
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "es-MX")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "No encontrado", rec.Body.String())
}

func TestTranslatorTranslatesForExplicitLocale(t *testing.T) {
	t.Parallel()
	translator, err := localization.New()
	require.NoError(t, err)

	assert.Equal(t, "Ссылка для входа", translator.TextIn("ru", localization.MagicLinkSubject))
	assert.Equal(t, "Restablecer contraseña", translator.TextIn("es", localization.PasswordResetSubject))
}

func TestTranslatorSpeaksEnglishOutsideAnHTTPRequest(t *testing.T) {
	t.Parallel()
	translator, err := localization.New()
	require.NoError(t, err)
	ctx := context.Background()

	assert.Equal(t, "en", translator.Locale(ctx))
	assert.Equal(t, "Wrong email or password", translator.Text(ctx, localization.WrongCredentials))
}

func TestTranslatorAnswersAStatusItHasNoTextFor(t *testing.T) {
	t.Parallel()
	translator, err := localization.New()
	require.NoError(t, err)
	ctx := context.Background()

	// net/http knows no text for this status either, so the answer is the
	// generic server error, never an empty string or an internal message.
	const unknownStatus = 799
	assert.Equal(t,
		translator.Text(ctx, localization.InternalServerError),
		translator.StatusText(ctx, unknownStatus))
	assert.NotEmpty(t, translator.StatusText(ctx, unknownStatus))
}

func TestTranslatorKeepsNetHTTPTextForAStatusItHasNoCatalogEntryFor(t *testing.T) {
	t.Parallel()
	translator, err := localization.New()
	require.NoError(t, err)

	// A real status the catalog does not translate reads as net/http's English
	// text, rather than hiding behind the generic server error.
	assert.Equal(t, "I'm a teapot", translator.StatusText(context.Background(), http.StatusTeapot))
}
