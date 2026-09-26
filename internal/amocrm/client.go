// Package amocrm adapts lead snapshots to amoCRM's unsorted-forms HTTP API.
package amocrm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"

	"github.com/samber/lo"

	"hexletbasics/internal/amocrm/generated"
	"hexletbasics/internal/events"
)

const (
	leadPipelineID       int64 = 9_614_774
	responsibleUserID    int64 = 7_877_026
	maxErrorMessageBytes       = 8 << 10
	source                     = "lead_form"
)

// Client sends lead-created snapshots to amoCRM.
type Client struct {
	baseURL string
	token   string
	api     *generated.Client
	http    *http.Client
	payload payloadBuilder
	initErr error
}

// NewClient builds the external integration adapter.
func NewClient(baseURL, token, ymCounter string) *Client {
	client := &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		http: &http.Client{
			Timeout: 15 * time.Second,
			// Redirects from an account-specific endpoint indicate configuration
			// drift. Do not turn one River attempt into multiple HTTP requests.
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		payload: newPayloadBuilder(ymCounter),
	}
	api, err := generated.NewClient(
		client.baseURL,
		bearerToken(token),
		generated.WithClient(statusClient{http: client.http}),
	)
	if err != nil {
		client.initErr = fmt.Errorf("initialize amoCRM client: %w", err)
		return client
	}
	client.api = api
	return client
}

// CreateLead creates an unsorted form lead. River owns retries around this
// network call; the generated client deliberately uses a plain http.Client.
func (c *Client) CreateLead(ctx context.Context, event events.LeadCreated) error {
	if c.baseURL == "" || c.token == "" {
		return newRequestError(0, "amoCRM is not configured", nil, false)
	}
	if c.initErr != nil {
		return newRequestError(0, "amoCRM client initialization failed", c.initErr, false)
	}
	body := []generated.UnsortedFormCreateItem{c.payload.build(event)}
	if _, err := c.api.UnsortedLeadsCreateForms(ctx, body); err != nil {
		return classifyRequestError(err)
	}
	return nil
}

type payloadBuilder struct {
	ymCounter string
}

func newPayloadBuilder(ymCounter string) payloadBuilder {
	return payloadBuilder{ymCounter: ymCounter}
}

func (b payloadBuilder) build(event events.LeadCreated) generated.UnsortedFormCreateItem {
	contact := generated.ContactCreate{
		Name: generated.NewOptString(lo.CoalesceOrEmpty(
			lo.FromPtr(event.FirstName),
			lo.FromPtr(event.Email),
			lo.FromPtr(event.Phone),
			lo.FromPtr(event.Telegram),
			lo.FromPtr(event.WhatsApp),
			"Unknown",
		)),
		FirstName: optString(event.FirstName),
		LastName:  optString(event.LastName),
		CustomFieldsValues: customFields(
			customField{code: "EMAIL", value: event.Email},
			customField{code: "PHONE", value: event.Phone},
		),
	}
	landing := landingQuery(event.LandingPage)
	lead := generated.LeadCreate{
		Name:              generated.NewOptString(lo.CoalesceOrEmpty(lo.FromPtr(event.Email), "Lead from "+source)),
		PipelineID:        generated.NewOptInt64(leadPipelineID),
		ResponsibleUserID: generated.NewOptInt64(responsibleUserID),
		CustomFieldsValues: customFields(
			customField{id: 316_913, code: "UTM_CONTENT", value: event.UTMContent},
			customField{id: 316_915, code: "UTM_MEDIUM", value: event.UTMMedium},
			customField{id: 316_917, code: "UTM_CAMPAIGN", value: event.UTMCampaign},
			customField{id: 316_919, code: "UTM_SOURCE", value: event.UTMSource},
			customField{id: 316_921, code: "UTM_TERM", value: event.UTMTerm},
			customField{id: 316_923, code: "UTM_REFERRER", value: event.Referrer},
			customField{id: 316_927, code: "REFERRER", value: event.Referrer},
			customField{id: 316_937, code: "FROM", value: landing("from")},
			customField{id: 316_941, code: "_YM_UID", value: event.YMClientID},
			customField{id: 316_943, code: "_YM_COUNTER", value: lo.ToPtr(b.ymCounter)},
			customField{id: 316_945, code: "GCLID", value: landing("gclid")},
			// Legacy sent the Metrika client id here ahead of the landing
			// page's yclid; kept as is, since amoCRM reports are built on it.
			customField{id: 316_947, code: "YCLID", value: lo.CoalesceOrEmpty(event.YMClientID, landing("yclid"))},
			customField{id: 316_949, code: "FBCLID", value: landing("fbclid")},
			customField{id: 957_711, code: "GA_UTM", value: landing("ga_utm")},
			// The source-form field has no code in amoCRM, only an id.
			customField{id: 936_587, value: lo.ToPtr(source)},
		),
	}
	return generated.UnsortedFormCreateItem{
		SourceUID:  fmt.Sprintf("%s-%d", source, event.LeadID),
		SourceName: source,
		Metadata: generated.FormMetadata{
			FormID:     generated.NewOptFormMetadataFormID(generated.NewStringFormMetadataFormID(source)),
			FormName:   generated.NewOptString(source),
			FormPage:   optString(event.LandingPage),
			FormSentAt: generated.NewOptInt64(event.OccurredAt.Unix()),
			IP:         optString(ipv4(event.IP)),
			Referer:    optString(event.Referrer),
		},
		Embedded: generated.NewOptUnsortedEmbeddedCreate(generated.UnsortedEmbeddedCreate{
			Contacts: []generated.ContactCreate{contact},
			Leads:    []generated.LeadCreate{lead},
		}),
	}
}

// landingQuery reads the click ids legacy parsed out of the landing page's
// query string. A page that does not parse contributes nothing.
func landingQuery(landingPage *string) func(key string) *string {
	var query url.Values
	if parsed, err := url.Parse(lo.FromPtr(landingPage)); err == nil {
		query = parsed.Query()
	}
	return func(key string) *string {
		value := query.Get(key)
		if value == "" {
			return nil
		}
		return &value
	}
}

// ipv4 keeps only an IPv4 address: amoCRM rejects IPv6 in the form metadata.
func ipv4(ip *string) *string {
	addr, err := netip.ParseAddr(lo.FromPtr(ip))
	if err != nil || !addr.Unmap().Is4() {
		return nil
	}
	return lo.ToPtr(addr.Unmap().String())
}

type customField struct {
	id    int64
	code  string
	value *string
}

func customFields(fields ...customField) []generated.CustomFieldValue {
	return lo.FilterMap(fields, func(input customField, _ int) (generated.CustomFieldValue, bool) {
		value := lo.FromPtr(input.value)
		if value == "" {
			return generated.CustomFieldValue{}, false
		}
		field := generated.CustomFieldValue{
			Values: []generated.CustomFieldValueItem{
				{Value: generated.NewStringCustomFieldValueItemValue(value)},
			},
		}
		if input.code != "" {
			field.FieldCode = generated.NewOptString(input.code)
		}
		if input.id != 0 {
			field.FieldID = generated.NewOptInt64(input.id)
		}
		return field, true
	})
}

// bearerToken is the generated client's security source: amoCRM long-lived
// integration tokens are static, so every operation gets the same value.
type bearerToken string

func (t bearerToken) BearerAuth(context.Context, generated.OperationName) (generated.BearerAuth, error) {
	return generated.BearerAuth{Token: string(t)}, nil
}

// statusClient sits between the generated client and net/http so every
// failure is classified here, where the HTTP status is still known. ogen only
// decodes the contract's application/problem+json errors; any other non-2xx
// (an HTML 502 from a proxy, a plain-text 401, a redirect that is not
// followed) would otherwise reach the caller as a status-less content-type
// error.
type statusClient struct {
	http *http.Client
}

// Do marks transport failures retryable (the request may not have reached
// amoCRM) and turns undescribed non-2xx responses into status-classified
// errors, so the only errors left for the generated decoder to produce come
// from responses amoCRM accepted.
func (c statusClient) Do(request *http.Request) (*http.Response, error) {
	response, err := c.http.Do(request)
	if err != nil {
		return nil, newRequestError(0, "send amoCRM lead", err, true)
	}
	if response.StatusCode < http.StatusMultipleChoices || isProblem(response) {
		return response, nil
	}
	defer func() { _ = response.Body.Close() }()
	body, _ := io.ReadAll(io.LimitReader(response.Body, maxErrorMessageBytes))
	return nil, newRequestError(
		response.StatusCode,
		string(body),
		nil,
		isRetryableStatus(response.StatusCode),
	)
}

// isProblem mirrors the contract's error media type: those bodies are left to
// the generated decoder, which returns them as ProblemDetailsStatusCode.
func isProblem(response *http.Response) bool {
	mediaType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	return err == nil && mediaType == "application/problem+json"
}

type requestError struct {
	statusCode int
	message    string
	cause      error
	retryable  bool
}

func newRequestError(statusCode int, message string, cause error, retryable bool) *requestError {
	return &requestError{
		statusCode: statusCode,
		message:    truncate(message, maxErrorMessageBytes),
		cause:      cause,
		retryable:  retryable,
	}
}

func (e *requestError) Error() string {
	if e.statusCode != 0 {
		return fmt.Sprintf("amoCRM returned %d: %s", e.statusCode, e.message)
	}
	if e.cause != nil {
		return fmt.Sprintf("%s: %s", e.message, truncate(e.cause.Error(), maxErrorMessageBytes))
	}
	return e.message
}

func (e *requestError) Unwrap() error {
	return e.cause
}

// Retryable reports whether another River attempt can plausibly succeed.
func (e *requestError) Retryable() bool {
	return e.retryable
}

// classifyRequestError maps the generated client's errors onto River's retry
// decision. Transport and undescribed-status errors were already classified by
// statusClient.
func classifyRequestError(err error) error {
	var classified *requestError
	if errors.As(err, &classified) {
		return classified
	}
	var problem *generated.ProblemDetailsStatusCode
	if errors.As(err, &problem) {
		return newRequestError(
			problem.StatusCode,
			joinProblem(problem.Response.Title, problem.Response.Detail),
			err,
			isRetryableStatus(problem.StatusCode),
		)
	}
	// What remains is either a request that could not be encoded (a bug) or a
	// 2xx whose body did not decode. amoCRM has accepted the lead in the second
	// case, so a retry would create a duplicate: fail permanently instead.
	return newRequestError(0, "read amoCRM response", err, false)
}

func isRetryableStatus(statusCode int) bool {
	return statusCode == http.StatusRequestTimeout ||
		statusCode == http.StatusTooManyRequests ||
		statusCode >= http.StatusInternalServerError && statusCode <= 599
}

func joinProblem(title, detail generated.OptString) string {
	parts := make([]string, 0, 2)
	if value, ok := title.Get(); ok && value != "" {
		parts = append(parts, value)
	}
	if value, ok := detail.Get(); ok && value != "" {
		parts = append(parts, value)
	}
	if len(parts) == 0 {
		return "request failed"
	}
	return strings.Join(parts, ": ")
}

func truncate(value string, limit int) string {
	if len(value) <= limit {
		return value
	}
	return value[:limit]
}

// optString keeps nil out of the request: ogen has no pointer constructor for
// its Opt types, and an unset Opt field is omitted from the JSON body.
func optString(value *string) generated.OptString {
	if value == nil {
		return generated.OptString{}
	}
	return generated.NewOptString(*value)
}
