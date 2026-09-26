package amocrm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
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

func TestCreateLeadSendsFirstVisitAttribution(t *testing.T) {
	landing := "https://code-basics.com/ru/languages/python?from=vk&gclid=g-1&yclid=y-1&fbclid=f-1&ga_utm=ga-1"
	referrer := "https://ya.ru/"
	ip := "203.0.113.7"

	body := sendLead(t, events.LeadCreated{
		LeadID: 10, LandingPage: &landing, Referrer: &referrer, IP: &ip,
		OccurredAt: time.Unix(123, 0),
	})

	metadata := requireMap(t, body["metadata"])
	assert.Equal(t, landing, metadata["form_page"])
	assert.Equal(t, referrer, metadata["referer"])
	assert.Equal(t, ip, metadata["ip"])

	fields := requireMap(t, requireSlice(t, requireMap(t, body["_embedded"])["leads"])[0])["custom_fields_values"]
	assertCustomField(t, fields, "UTM_REFERRER", 316_923, referrer)
	assertCustomField(t, fields, "REFERRER", 316_927, referrer)
	assertCustomField(t, fields, "FROM", 316_937, "vk")
	assertCustomField(t, fields, "GCLID", 316_945, "g-1")
	assertCustomField(t, fields, "YCLID", 316_947, "y-1")
	assertCustomField(t, fields, "FBCLID", 316_949, "f-1")
	assertCustomField(t, fields, "GA_UTM", 957_711, "ga-1")
	assertCustomFieldByID(t, fields, 936_587, "lead_form")
}

func TestCreateLeadPrefersTheMetrikaClientIDAsYclid(t *testing.T) {
	landing := "https://code-basics.com/?yclid=y-1"
	clientID := "1700000000123456789"

	body := sendLead(t, events.LeadCreated{LeadID: 10, LandingPage: &landing, YMClientID: &clientID})

	fields := requireMap(t, requireSlice(t, requireMap(t, body["_embedded"])["leads"])[0])["custom_fields_values"]
	assertCustomField(t, fields, "YCLID", 316_947, clientID)
	assertCustomField(t, fields, "_YM_UID", 316_941, clientID)
}

func TestCreateLeadDropsAnIPv6Address(t *testing.T) {
	ip := "2001:db8::1"

	body := sendLead(t, events.LeadCreated{LeadID: 10, IP: &ip})

	assert.NotContains(t, requireMap(t, body["metadata"]), "ip")
}

// sendLead posts the event through a stubbed transport and returns the one
// unsorted form it sent.
func sendLead(t *testing.T, event events.LeadCreated) map[string]any {
	t.Helper()
	var body []map[string]any
	client := NewClient("https://example.amocrm.test", "secret", "counter")
	client.http.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		const responseBody = `{"_total_items":1,"_embedded":{"unsorted":[]}}`
		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(bytes.NewBufferString(responseBody)),
			ContentLength: int64(len(responseBody)),
			Header:        http.Header{"Content-Type": []string{"application/hal+json"}},
		}, nil
	})
	require.NoError(t, client.CreateLead(t.Context(), event))
	require.Len(t, body, 1)
	return body[0]
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

func TestCreateLeadClassifiesNonProblemErrorByStatus(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		retryable bool
	}{
		{name: "unauthorized", status: http.StatusUnauthorized, retryable: false},
		{name: "bad gateway", status: http.StatusBadGateway, retryable: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient("https://example.amocrm.test", "secret", "")
			client.http.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				const body = "<html>upstream failure</html>"
				return &http.Response{
					StatusCode:    tt.status,
					Body:          io.NopCloser(bytes.NewBufferString(body)),
					ContentLength: int64(len(body)),
					Header:        http.Header{"Content-Type": []string{"text/html"}},
				}, nil
			})

			err := client.CreateLead(t.Context(), events.LeadCreated{LeadID: 1})

			require.ErrorContains(t, err, strconv.Itoa(tt.status))
			var classified interface{ Retryable() bool }
			require.True(t, errors.As(err, &classified))
			assert.Equal(t, tt.retryable, classified.Retryable())
		})
	}
}

func TestCreateLeadDoesNotRetryUnreadableSuccess(t *testing.T) {
	requests := 0
	client := NewClient("https://example.amocrm.test", "secret", "")
	client.http.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		requests++
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Body:       io.NopCloser(bytes.NewBufferString("")),
			Header:     http.Header{},
		}, nil
	})

	err := client.CreateLead(t.Context(), events.LeadCreated{LeadID: 1})

	require.Error(t, err)
	assert.Equal(t, 1, requests)
	var classified interface{ Retryable() bool }
	require.True(t, errors.As(err, &classified))
	assert.False(t, classified.Retryable(), "amoCRM accepted the lead; a retry would duplicate it")
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

func assertCustomFieldByID(t *testing.T, value any, id int, expected string) {
	t.Helper()
	for _, rawField := range requireSlice(t, value) {
		field := requireMap(t, rawField)
		if field["field_id"] != float64(id) {
			continue
		}
		assert.NotContains(t, field, "field_code")
		values := requireSlice(t, field["values"])
		assert.Equal(t, expected, requireMap(t, values[0])["value"])
		return
	}
	t.Errorf("custom field %d not found", id)
}
