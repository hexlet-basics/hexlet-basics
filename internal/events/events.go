// Package events owns the domain-event contracts and their delivery (ADR-0031).
//
// A domain event is an immutable fact a business module raises. Publishing it
// inserts, in the caller's transaction, one River job for each consumer that
// reacts to it; River owns retries, the dead-letter rows and execution history.
// An event no consumer reacts to leaves nothing behind, so a fact is added here
// together with the consumer that needs it.
package events

import (
	"time"
)

const (
	leadCreatedName = "lead_created"
	staffReadName   = "staff_read_personal_data"
)

// Event is closed to this package so arbitrary payloads cannot be published as
// domain facts without first defining their stable contract here.
type Event interface {
	eventName() string
}

// LeadCreated is the integration snapshot consumed by the amoCRM job.
type LeadCreated struct {
	LeadID      int     `json:"lead_id"`
	UserID      int     `json:"user_id"`
	UserName    string  `json:"user_name"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	YMClientID  *string `json:"ym_client_id"`
	UTMSource   *string `json:"utm_source"`
	UTMMedium   *string `json:"utm_medium"`
	UTMCampaign *string `json:"utm_campaign"`
	UTMTerm     *string `json:"utm_term"`
	UTMContent  *string `json:"utm_content"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Telegram    *string `json:"telegram"`
	WhatsApp    *string `json:"whatsapp"`
	// LandingPage, Referrer and IP are what the amoCRM payload used to read off
	// ahoy's visit (form page, referer, IP, and the click ids parsed from the
	// landing URL). The first two come from the browser's first-visit cookie,
	// the IP from the request (ADR-0015). IP is the raw client address; the
	// amoCRM adapter decides what it can send.
	LandingPage *string   `json:"landing_page"`
	Referrer    *string   `json:"referrer"`
	IP          *string   `json:"ip"`
	OccurredAt  time.Time `json:"occurred_at"`
}

func (LeadCreated) eventName() string { return leadCreatedName }

// StaffReadPersonalData is the fact that a staff member read personal data
// through the admin surface or over MCP (ADR-0020). The request publishes it and
// a worker job writes the Audit Record, so recording never slows the read.
//
// It carries identifiers and a count, never a field value: EntityIDs are the ids
// of the rows the response held, up to a cap (the whole count is in Count), so
// the record can say which learners were looked at without repeating what was
// seen. The rest is who and from where, copied from the request because the
// job has no request.
type StaffReadPersonalData struct {
	ActorUserID    *int      `json:"actor_user_id"`
	Channel        string    `json:"channel"`
	ConnectedAppID *int      `json:"connected_app_id"`
	TokenID        *string   `json:"token_id"`
	Operation      string    `json:"operation"`
	EntityType     string    `json:"entity_type"`
	EntityIDs      []string  `json:"entity_ids"`
	Count          int       `json:"count"`
	RequestID      string    `json:"request_id"`
	Address        string    `json:"address"`
	UserAgent      string    `json:"user_agent"`
	Locale         string    `json:"locale"`
	OccurredAt     time.Time `json:"occurred_at"`
}

func (StaffReadPersonalData) eventName() string { return staffReadName }
