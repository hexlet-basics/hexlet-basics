package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/testsupport"
)

// A request that matches no route, or matches one with the wrong method, is
// Problem Details like every other non-2xx (ADR-0028), not plain text.
func TestAnUnmatchedRouteAnswersWithProblemDetails(t *testing.T) {
	t.Parallel()
	harness := testsupport.NewHarness(t)

	cases := []struct {
		name, method, path string
		status             int
	}{
		{"unknown path", http.MethodGet, "/api/no-such-route", http.StatusNotFound},
		{"wrong method", http.MethodDelete, "/api/languages", http.StatusMethodNotAllowed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			response := harness.Send(t, tc.method, tc.path)
			defer func() { _ = response.Body.Close() }()

			assert.Equal(t, tc.status, response.StatusCode)
			assert.Equal(t, "application/problem+json", response.Header.Get("Content-Type"))
			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			var problem struct {
				Status int    `json:"status"`
				Title  string `json:"title"`
			}
			require.NoError(t, json.Unmarshal(body, &problem))
			assert.Equal(t, tc.status, problem.Status)
			assert.NotEmpty(t, problem.Title)
		})
	}
}
