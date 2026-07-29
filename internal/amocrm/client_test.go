package amocrm

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/events"
)

func TestCreateLeadSendsUnsortedForm(t *testing.T) {
	var body []map[string]any
	client := NewClient("https://example.amocrm.test", "secret", "counter")
	client.http.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/api/v4/leads/unsorted/forms", r.URL.Path)
		assert.Equal(t, "Bearer secret", r.Header.Get("Authorization"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Header:     make(http.Header),
		}, nil
	})
	email := "ada@example.com"
	firstName := "Ada"

	err := client.CreateLead(t.Context(), events.LeadCreated{
		LeadID: 10, UserID: 20, Email: &email, FirstName: &firstName,
		OccurredAt: time.Unix(123, 0),
	})
	require.NoError(t, err)
	require.Len(t, body, 1)
	assert.Equal(t, "lead_form-10", body[0]["source_uid"])
	assert.Equal(t, "lead_form", body[0]["source_name"])
}

func TestCreateLeadReturnsRemoteError(t *testing.T) {
	client := NewClient("https://example.amocrm.test", "secret", "")
	client.http.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewBufferString("invalid lead")),
			Header:     make(http.Header),
		}, nil
	})
	err := client.CreateLead(
		t.Context(),
		events.LeadCreated{LeadID: 1},
	)
	require.ErrorContains(t, err, "amoCRM returned 400")
	require.ErrorContains(t, err, "invalid lead")
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
