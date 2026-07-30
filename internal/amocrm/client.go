// Package amocrm adapts lead snapshots to amoCRM's unsorted-forms HTTP API.
package amocrm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	kiota "github.com/microsoft/kiota-abstractions-go"
	kiotaauth "github.com/microsoft/kiota-abstractions-go/authentication"
	kiotaserialization "github.com/microsoft/kiota-abstractions-go/serialization"
	kiotahttp "github.com/microsoft/kiota-http-go"
	kiotajson "github.com/microsoft/kiota-serialization-json-go"
	"github.com/samber/lo"

	"hexletbasics/internal/amocrm/generated"
	"hexletbasics/internal/amocrm/generated/models"
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
	api     *generated.APIClient
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
	}
	// Kiota's registries normalize vendor JSON media types such as
	// application/hal+json and application/problem+json to application/json.
	parseNodeFactory := kiotaserialization.NewParseNodeFactoryRegistry()
	parseNodeFactory.ContentTypeAssociatedFactories["application/json"] = kiotajson.NewJsonParseNodeFactory()
	writerFactory := kiotaserialization.NewSerializationWriterFactoryRegistry()
	writerFactory.ContentTypeAssociatedFactories["application/json"] = kiotajson.NewJsonSerializationWriterFactory()
	adapter, err := kiotahttp.NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		bearerAuthProvider{token: token},
		parseNodeFactory,
		writerFactory,
		client.http,
	)
	if err != nil {
		client.initErr = fmt.Errorf("initialize amoCRM request adapter: %w", err)
		return client
	}
	adapter.SetBaseUrl(client.baseURL)
	client.api = generated.NewAPIClient(adapter)
	client.payload = newPayloadBuilder(ymCounter)
	return client
}

// CreateLead creates an unsorted form lead. River owns retries around this
// network call; the Kiota adapter deliberately uses a plain http.Client.
func (c *Client) CreateLead(ctx context.Context, event events.LeadCreated) error {
	if c.baseURL == "" || c.token == "" {
		return newRequestError(0, "amoCRM is not configured", nil, false)
	}
	if c.initErr != nil {
		return newRequestError(0, "amoCRM client initialization failed", c.initErr, false)
	}
	body := []models.UnsortedFormCreateItemable{c.payload.build(event)}
	_, err := c.api.Api().V4().Leads().Unsorted().Forms().Post(ctx, body, nil)
	if err != nil {
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

func (b payloadBuilder) build(event events.LeadCreated) models.UnsortedFormCreateItemable {
	metadata := models.NewFormMetadata()
	formID := models.NewFormMetadata_FormMetadata_form_id()
	formID.SetString(lo.ToPtr(source))
	metadata.SetFormId(formID)
	metadata.SetFormName(lo.ToPtr(source))
	sentAt := event.OccurredAt.Unix()
	metadata.SetFormSentAt(&sentAt)

	contact := models.NewContactCreate()
	contact.SetName(lo.ToPtr(firstNonEmpty(
		event.FirstName,
		event.Email,
		event.Phone,
		event.Telegram,
		event.WhatsApp,
		lo.ToPtr("Unknown"),
	)))
	contact.SetFirstName(event.FirstName)
	contact.SetLastName(event.LastName)
	contact.SetCustomFieldsValues(customFields(
		customField{code: "EMAIL", value: event.Email},
		customField{code: "PHONE", value: event.Phone},
	))

	lead := models.NewLeadCreate()
	lead.SetName(lo.ToPtr(firstNonEmpty(event.Email, lo.ToPtr("Lead from "+source))))
	lead.SetPipelineId(lo.ToPtr(leadPipelineID))
	lead.SetResponsibleUserId(lo.ToPtr(responsibleUserID))
	lead.SetCustomFieldsValues(customFields(
		customField{id: 316_913, code: "UTM_CONTENT", value: event.UTMContent},
		customField{id: 316_915, code: "UTM_MEDIUM", value: event.UTMMedium},
		customField{id: 316_917, code: "UTM_CAMPAIGN", value: event.UTMCampaign},
		customField{id: 316_919, code: "UTM_SOURCE", value: event.UTMSource},
		customField{id: 316_921, code: "UTM_TERM", value: event.UTMTerm},
		customField{id: 316_941, code: "_YM_UID", value: event.YMClientID},
		customField{id: 316_943, code: "_YM_COUNTER", value: lo.ToPtr(b.ymCounter)},
	))

	embedded := models.NewUnsortedEmbeddedCreate()
	embedded.SetContacts([]models.ContactCreateable{contact})
	embedded.SetLeads([]models.LeadCreateable{lead})

	item := models.NewUnsortedFormCreateItem()
	item.SetSourceUid(lo.ToPtr(fmt.Sprintf("%s-%d", source, event.LeadID)))
	item.SetSourceName(lo.ToPtr(source))
	item.SetMetadata(metadata)
	item.SetEmbedded(embedded)
	return item
}

type customField struct {
	id    int64
	code  string
	value *string
}

func customFields(fields ...customField) []models.CustomFieldValueable {
	result := make([]models.CustomFieldValueable, 0, len(fields))
	for _, input := range fields {
		if input.value == nil || *input.value == "" {
			continue
		}
		value := models.NewCustomFieldValueItem_CustomFieldValueItem_value()
		value.SetString(input.value)
		valueItem := models.NewCustomFieldValueItem()
		valueItem.SetValue(value)

		field := models.NewCustomFieldValue()
		field.SetFieldCode(&input.code)
		if input.id != 0 {
			field.SetFieldId(&input.id)
		}
		field.SetValues([]models.CustomFieldValueItemable{valueItem})
		result = append(result, field)
	}
	return result
}

type bearerAuthProvider struct {
	token string
}

var _ kiotaauth.AuthenticationProvider = bearerAuthProvider{}

func (p bearerAuthProvider) AuthenticateRequest(
	_ context.Context,
	request *kiota.RequestInformation,
	_ map[string]any,
) error {
	request.Headers.TryAdd("Authorization", "Bearer "+p.token)
	return nil
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

func classifyRequestError(err error) error {
	var problem *models.ProblemDetails
	if errors.As(err, &problem) {
		statusCode := problem.GetStatusCode()
		return newRequestError(
			statusCode,
			joinProblem(problem.GetTitle(), problem.GetDetail()),
			err,
			isRetryableStatus(statusCode),
		)
	}
	var apiError kiota.ApiErrorable
	if errors.As(err, &apiError) {
		statusCode := apiError.GetStatusCode()
		return newRequestError(statusCode, err.Error(), err, isRetryableStatus(statusCode))
	}
	return newRequestError(0, "send amoCRM lead", err, true)
}

func isRetryableStatus(statusCode int) bool {
	return statusCode == http.StatusRequestTimeout ||
		statusCode == http.StatusTooManyRequests ||
		statusCode >= http.StatusInternalServerError && statusCode <= 599
}

func joinProblem(title, detail *string) string {
	parts := make([]string, 0, 2)
	if title != nil && *title != "" {
		parts = append(parts, *title)
	}
	if detail != nil && *detail != "" {
		parts = append(parts, *detail)
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

func firstNonEmpty(values ...*string) string {
	for _, candidate := range values {
		if candidate != nil && *candidate != "" {
			return *candidate
		}
	}
	return ""
}
