package amocrm

import (
	"bytes"
	"encoding/json"
	"errors"
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
		const responseBody = `{"_total_items":1,"_embedded":{"unsorted":[]}}`
		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(bytes.NewBufferString(responseBody)),
			ContentLength: int64(len(responseBody)),
			Header:        http.Header{"Content-Type": []string{"application/hal+json"}},
		}, nil
	})
	email := "ada@example.com"
	firstName := "Ada"
	utmSource := "newsletter"

	err := client.CreateLead(t.Context(), events.LeadCreated{
		LeadID: 10, UserID: 20, Email: &email, FirstName: &firstName,
		UTMSource: &utmSource, OccurredAt: time.Unix(123, 0),
	})
	require.NoError(t, err)
	require.Len(t, body, 1)
	assert.Equal(t, "lead_form-10", body[0]["source_uid"])
	assert.Equal(t, "lead_form", body[0]["source_name"])

	metadata := requireMap(t, body[0]["metadata"])
	assert.Equal(t, "lead_form", metadata["form_id"])
	assert.Equal(t, float64(123), metadata["form_sent_at"])

	embedded := requireMap(t, body[0]["_embedded"])
	contacts := requireSlice(t, embedded["contacts"])
	contact := requireMap(t, contacts[0])
	assert.Equal(t, "Ada", contact["name"])
	assertCustomField(t, contact["custom_fields_values"], "EMAIL", 0, email)

	leads := requireSlice(t, embedded["leads"])
	lead := requireMap(t, leads[0])
	assert.Equal(t, float64(leadPipelineID), lead["pipeline_id"])
	assert.Equal(t, float64(responsibleUserID), lead["responsible_user_id"])
	assertCustomField(t, lead["custom_fields_values"], "UTM_SOURCE", 316_919, utmSource)
	assertCustomField(t, lead["custom_fields_values"], "_YM_COUNTER", 316_943, "counter")
}

func TestCreateLeadReturnsRemoteError(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		retryable bool
	}{
		{name: "bad request", status: http.StatusBadRequest, retryable: false},
		{name: "unauthorized", status: http.StatusUnauthorized, retryable: false},
		{name: "payment required", status: http.StatusPaymentRequired, retryable: false},
		{name: "forbidden", status: http.StatusForbidden, retryable: false},
		{name: "request timeout", status: http.StatusRequestTimeout, retryable: true},
		{name: "rate limited", status: http.StatusTooManyRequests, retryable: true},
		{name: "server error", status: http.StatusServiceUnavailable, retryable: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requests := 0
			client := NewClient("https://example.amocrm.test", "secret", "")
			client.http.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				requests++
				const body = `{"title":"invalid lead","detail":"email is malformed"}`
				return &http.Response{
					StatusCode:    tt.status,
					Body:          io.NopCloser(bytes.NewBufferString(body)),
					ContentLength: int64(len(body)),
					Header:        http.Header{"Content-Type": []string{"application/problem+json"}},
				}, nil
			})

			err := client.CreateLead(t.Context(), events.LeadCreated{LeadID: 1})

			require.ErrorContains(t, err, "invalid lead")
			assert.Equal(t, 1, requests, "one River attempt must make one HTTP request")
			var classified interface{ Retryable() bool }
			require.True(t, errors.As(err, &classified))
			assert.Equal(t, tt.retryable, classified.Retryable())
		})
	}
}

func TestCreateLeadReturnsRetryableTransportError(t *testing.T) {
	client := NewClient("https://example.amocrm.test", "secret", "")
	client.http.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("connection reset")
	})

	err := client.CreateLead(t.Context(), events.LeadCreated{LeadID: 1})

	require.ErrorContains(t, err, "connection reset")
	var classified interface{ Retryable() bool }
	require.True(t, errors.As(err, &classified))
	assert.True(t, classified.Retryable())
}

func TestCreateLeadTreatsMissingConfigurationAsPermanent(t *testing.T) {
	err := NewClient("", "", "").CreateLead(t.Context(), events.LeadCreated{LeadID: 1})

	var classified interface{ Retryable() bool }
	require.True(t, errors.As(err, &classified))
	assert.False(t, classified.Retryable())
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func requireMap(t *testing.T, value any) map[string]any {
	t.Helper()
	result, ok := value.(map[string]any)
	require.True(t, ok)
	return result
}

func requireSlice(t *testing.T, value any) []any {
	t.Helper()
	result, ok := value.([]any)
	require.True(t, ok)
	require.NotEmpty(t, result)
	return result
}

func assertCustomField(t *testing.T, value any, code string, id int, expected string) {
	t.Helper()
	fields := requireSlice(t, value)
	for _, rawField := range fields {
		field := requireMap(t, rawField)
		if field["field_code"] != code {
			continue
		}
		if id != 0 {
			assert.Equal(t, float64(id), field["field_id"])
		}
		values := requireSlice(t, field["values"])
		assert.Equal(t, expected, requireMap(t, values[0])["value"])
		return
	}
	t.Errorf("custom field %q not found", code)
}
