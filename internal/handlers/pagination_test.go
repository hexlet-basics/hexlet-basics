package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/localization"
)

type paginationContractHandler struct {
	api.UnimplementedHandler
	calls int
}

func (h *paginationContractHandler) ListPublicReviews(_ context.Context, params api.ListPublicReviewsParams) (*api.ReviewPage, error) {
	h.calls++
	page := newPagination(params.Page, params.PerPage)
	return &api.ReviewPage{
		Items:   []api.Review{},
		Total:   0,
		Page:    page.Page,
		PerPage: page.PerPage,
	}, nil
}

func newPaginationRouter(t *testing.T, handler api.Handler) http.Handler {
	t.Helper()

	translator, err := localization.New()
	require.NoError(t, err)
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{})
	require.NoError(t, err)
	server, err := api.NewServer(
		handler,
		api.WithErrorHandler(NewAPIErrorHandler(
			translator,
			slog.New(slog.NewTextHandler(io.Discard, nil)),
			sentryClient,
		).Handle),
		api.WithNotFound(NewNotFoundHandler(translator)),
		api.WithMethodNotAllowed(NewMethodNotAllowedHandler(translator)),
	)
	require.NoError(t, err)
	return translator.Middleware(server)
}

func TestPaginationContractRejectsInvalidQuery(t *testing.T) {
	tests := []string{
		"/reviews?page=0",
		"/reviews?page=-100",
		"/reviews?perPage=0",
		"/reviews?perPage=-1",
		"/reviews?perPage=101",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			handler := &paginationContractHandler{}
			router := newPaginationRouter(t, handler)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Zero(t, handler.calls, "generated decoder must reject the request before the handler")
		})
	}
}

func TestPaginationContractErrorUsesRequestLocale(t *testing.T) {
	router := newPaginationRouter(t, &paginationContractHandler{})
	req := httptest.NewRequest(http.MethodGet, "/reviews?page=0", nil)
	req.Header.Set("Accept-Language", "ru")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var body struct {
		Error string `json:"error"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t, "Некорректный запрос", body.Error)
}

func TestRouterErrorsUseRequestLocale(t *testing.T) {
	router := newPaginationRouter(t, &paginationContractHandler{})
	tests := []struct {
		name       string
		method     string
		path       string
		status     int
		body       string
		allowValue string
	}{
		{
			name:   "not found",
			method: http.MethodGet,
			path:   "/missing",
			status: http.StatusNotFound,
			body:   "No encontrado",
		},
		{
			name:       "method not allowed",
			method:     http.MethodDelete,
			path:       "/reviews",
			status:     http.StatusMethodNotAllowed,
			body:       "Método no permitido",
			allowValue: "GET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Accept-Language", "es")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.status, rec.Code)
			assert.Equal(t, tt.body, strings.TrimSpace(rec.Body.String()))
			if tt.allowValue != "" {
				assert.Contains(t, rec.Header().Get("Allow"), tt.allowValue)
			}
		})
	}
}

func TestPaginationContractAcceptsBoundsAndResolvesDefaults(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantPage    int32
		wantPerPage int32
	}{
		{name: "defaults", path: "/reviews", wantPage: 1, wantPerPage: defaultPerPage},
		{name: "minimums", path: "/reviews?page=1&perPage=1", wantPage: 1, wantPerPage: 1},
		{name: "maximum page size", path: "/reviews?page=123&perPage=100", wantPage: 123, wantPerPage: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &paginationContractHandler{}
			router := newPaginationRouter(t, handler)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			require.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, 1, handler.calls)

			var page api.ReviewPage
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&page))
			assert.Equal(t, tt.wantPage, page.Page)
			assert.Equal(t, tt.wantPerPage, page.PerPage)
		})
	}
}
