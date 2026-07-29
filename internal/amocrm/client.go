// Package amocrm adapts lead snapshots to amoCRM's unsorted-forms HTTP API.
package amocrm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"hexletbasics/internal/events"
)

const (
	leadPipelineID        = 9_614_774
	responsibleUserID     = 7_877_026
	maxErrorResponseBytes = 8 << 10
)

// Client sends lead-created snapshots to amoCRM.
type Client struct {
	baseURL   string
	token     string
	ymCounter string
	http      *http.Client
}

// NewClient builds the external integration adapter.
func NewClient(baseURL, token, ymCounter string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token, ymCounter: ymCounter,
		http: &http.Client{Timeout: 15 * time.Second},
	}
}

// CreateLead creates an unsorted form lead. River owns retries around this
// network call; a non-2xx response is returned with a bounded diagnostic body.
func (c *Client) CreateLead(ctx context.Context, event events.LeadCreated) error {
	if c.baseURL == "" || c.token == "" {
		return fmt.Errorf("amoCRM is not configured")
	}
	body, err := json.Marshal(c.payload(event))
	if err != nil {
		return fmt.Errorf("encode amoCRM lead: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/v4/leads/unsorted/forms",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create amoCRM request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("send amoCRM lead: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorResponseBytes))
		return fmt.Errorf("amoCRM returned %d: %s", resp.StatusCode, strings.TrimSpace(string(detail)))
	}
	return nil
}

func (c *Client) payload(event events.LeadCreated) []any {
	const source = "lead_form"
	contactFields := customFields(nil, map[string]*string{
		"EMAIL": event.Email,
		"PHONE": event.Phone,
	})
	leadFields := customFields(map[string]int{
		"UTM_CONTENT":  316_913,
		"UTM_MEDIUM":   316_915,
		"UTM_CAMPAIGN": 316_917,
		"UTM_SOURCE":   316_919,
		"UTM_TERM":     316_921,
		"_YM_UID":      316_941,
		"_YM_COUNTER":  316_943,
	}, map[string]*string{
		"UTM_CONTENT":  event.UTMContent,
		"UTM_MEDIUM":   event.UTMMedium,
		"UTM_CAMPAIGN": event.UTMCampaign,
		"UTM_SOURCE":   event.UTMSource,
		"UTM_TERM":     event.UTMTerm,
		"_YM_UID":      event.YMClientID,
		"_YM_COUNTER":  stringPtr(c.ymCounter),
	})
	leadName := firstNonEmpty(event.Email, stringPtr("Lead from "+source))
	contactName := firstNonEmpty(event.FirstName, event.Email, event.Phone, event.Telegram, event.WhatsApp, stringPtr("Unknown"))
	return []any{map[string]any{
		"source_uid":  fmt.Sprintf("%s-%d", source, event.LeadID),
		"source_name": source,
		"metadata": map[string]any{
			"form_id": source, "form_name": source,
			"form_sent_at": event.OccurredAt.Unix(),
		},
		"_embedded": map[string]any{
			"contacts": []any{map[string]any{
				"name": contactName, "first_name": value(event.FirstName),
				"last_name": value(event.LastName), "custom_fields_values": contactFields,
			}},
			"leads": []any{map[string]any{
				"name": leadName, "pipeline_id": leadPipelineID,
				"responsible_user_id":  responsibleUserID,
				"custom_fields_values": leadFields,
			}},
		},
	}}
}

func customFields(ids map[string]int, values map[string]*string) []any {
	result := make([]any, 0, len(values))
	for code, raw := range values {
		if raw == nil || *raw == "" {
			continue
		}
		field := map[string]any{
			"field_code": code,
			"values":     []any{map[string]any{"value": *raw}},
		}
		if id := ids[code]; id != 0 {
			field["field_id"] = id
		}
		result = append(result, field)
	}
	return result
}

func firstNonEmpty(values ...*string) string {
	for _, candidate := range values {
		if candidate != nil && *candidate != "" {
			return *candidate
		}
	}
	return ""
}

func value(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}

func stringPtr(value string) *string { return &value }
