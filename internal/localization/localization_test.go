package localization_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/localization"
)

func TestTranslatorSelectsRequestLanguage(t *testing.T) {
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
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept-Language", tt.header)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expected, rec.Body.String())
			assert.Contains(t, rec.Header().Values("Vary"), "Accept-Language")
		})
	}
}

func TestTranslatorLocalizesHTTPStatus(t *testing.T) {
	translator, err := localization.New()
	require.NoError(t, err)

	handler := translator.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(translator.StatusText(r.Context(), http.StatusNotFound)))
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "es-MX")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "No encontrado", rec.Body.String())
}
