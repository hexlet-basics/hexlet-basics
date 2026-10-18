package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/paging"
)

type paginationContractHandler struct {
	api.UnimplementedHandler

	calls int
}

func (h *paginationContractHandler) ListPublicReviews(_ context.Context, params api.ListPublicReviewsParams) (*api.ReviewPage, error) {
	h.calls++
	page := paging.New(params.Page, params.PerPage)
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
	errorHandler := NewAPIErrorHandler(translator, slog.New(slog.DiscardHandler))
	server, err := api.NewServer(
		handler,
		nil,
		api.WithErrorHandler(errorHandler.Write),
		api.WithNotFound(NewNotFoundHandler(errorHandler)),
		api.WithMethodNotAllowed(NewMethodNotAllowedHandler(errorHandler)),
	)
	require.NoError(t, err)
	return translator.Middleware(server)
}

func TestPaginationContractRejectsInvalidQuery(t *testing.T) {
	t.Parallel()
	tests := []string{
		"/api/reviews?page=0",
		"/api/reviews?page=-100",
		"/api/reviews?perPage=0",
		"/api/reviews?perPage=-1",
		"/api/reviews?perPage=101",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			handler := &paginationContractHandler{}
			router := newPaginationRouter(t, handler)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil))

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Zero(t, handler.calls, "generated decoder must reject the request before the handler")
		})
	}
}

func TestPaginationContractErrorUsesRequestLocale(t *testing.T) {
	t.Parallel()
	router := newPaginationRouter(t, &paginationContractHandler{})
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/reviews?page=0", nil)
	req.Header.Set("Accept-Language", "ru")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var body api.ProblemDetails
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t, "Некорректный запрос", body.Title)
}

func TestRouterErrorsUseRequestLocale(t *testing.T) {
	t.Parallel()
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
			path:       "/api/reviews",
			status:     http.StatusMethodNotAllowed,
			body:       "Método no permitido",
			allowValue: "GET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), tt.method, tt.path, nil)
			req.Header.Set("Accept-Language", "es")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.status, rec.Code)
			assert.JSONEq(t, fmt.Sprintf(`{"type":"about:blank","title":%q,"status":%d}`, tt.body, tt.status), rec.Body.String())
			if tt.allowValue != "" {
				assert.Contains(t, rec.Header().Get("Allow"), tt.allowValue)
			}
		})
	}
}

func TestPaginationContractAcceptsBoundsAndResolvesDefaults(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		path        string
		wantPage    int32
		wantPerPage int32
	}{
		{name: "defaults", path: "/api/reviews", wantPage: 1, wantPerPage: 20},
		{name: "minimums", path: "/api/reviews?page=1&perPage=1", wantPage: 1, wantPerPage: 1},
		{name: "maximum page size", path: "/api/reviews?page=123&perPage=100", wantPage: 123, wantPerPage: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &paginationContractHandler{}
			router := newPaginationRouter(t, handler)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, tt.path, nil))

			require.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, 1, handler.calls)

			var page api.ReviewPage
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&page))
			assert.Equal(t, tt.wantPage, page.Page)
			assert.Equal(t, tt.wantPerPage, page.PerPage)
		})
	}
}
