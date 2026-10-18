// Package audit writes the Audit Records of ADR-0020: one append-only row per
// administrative operation, in the same transaction as the change.
//
// Three small pieces make up the single shared place, and none of them knows
// about a particular operation:
//
//   - Hook is an ent hook installed on every client the store builds. While a
//     Recorder is in the context it notes each row change (type, id, and the old
//     and new value of each changed field, secrets by name only); otherwise it
//     does nothing.
//   - Recorder collects those changes for one operation.
//   - Write turns a Recorder plus the request's identity into the AuditRecord
//     row, through the ent client of the operation's transaction.
//
// The unit of work that opens the transaction and calls Write lives with the
// transport (handlers.NewAuditMiddleware for the website; the MCP server calls
// the same operations through the same middleware, ADR-0021), because only it
// knows the request. This package depends on ent only.
package audit

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"hexletbasics/ent"
)

// Channel is how the actor reached the operation.
type Channel string

const (
	// ChannelWeb is the website's admin UI.
	ChannelWeb Channel = "web"
	// ChannelMCP is a tool call from a Connected App over MCP.
	ChannelMCP Channel = "mcp"
)

// Outcome is how the operation ended.
type Outcome string

const (
	// OutcomeSucceeded is a change that committed.
	OutcomeSucceeded Outcome = "succeeded"
	// OutcomeForbidden is a refusal by the permission check (ticket #845).
	OutcomeForbidden Outcome = "forbidden"
	// OutcomeFailed is an operation that started and did not complete (#845).
	OutcomeFailed Outcome = "failed"
)

// Source says through which channel an operation arrived and, for MCP, on
// behalf of which Connected App and token. The zero value is the website.
type Source struct {
	Channel        Channel
	ConnectedAppID *int
	TokenID        *string
}

type sourceKey struct{}

// WithSource tells the operation layer how the caller arrived; the MCP server
// sets it before dispatching a tool call. Without it an operation is the website's.
func WithSource(ctx context.Context, s Source) context.Context {
	return context.WithValue(ctx, sourceKey{}, s)
}

// SourceFrom is the source WithSource set, the website when none was.
func SourceFrom(ctx context.Context) Source {
	if s, ok := ctx.Value(sourceKey{}).(Source); ok && s.Channel != "" {
		return s
	}

	return Source{
		Channel:        ChannelWeb,
		ConnectedAppID: nil,
		TokenID:        nil,
	}
}

// Entry is everything an Audit Record says except the changes, which come from
// the Recorder.
type Entry struct {
	Actor     *int
	Source    Source
	RequestID string
	Address   string
	UserAgent string
	Operation string
	Locale    string
	Outcome   Outcome

	// The rest is for a read of personal data, which is recorded from a domain
	// event some time after the request (ADR-0020). OccurredAt is when the read
	// happened, not when the record was written; EntityType and EntityID file the
	// record under what was read (the id only for a read of one row); Reads says
	// how much was touched. A change leaves all four empty: its time is now, and
	// its entity comes from the Recorder.
	OccurredAt time.Time
	EntityType string
	EntityID   string
	Reads      *Reads
}

// MaxReadIDs caps the ids a read record keeps. A page is far smaller, and a
// record that grew with the table would turn a bulk read into a bulk write.
const MaxReadIDs = 200

// Reads is what a read of personal data touched: the number of rows, and their
// ids up to MaxReadIDs. Identifiers and counts only, never a field value.
type Reads struct {
	Count int      `json:"count"`
	IDs   []string `json:"ids,omitempty"`
}

// FieldChange is the old and new value of one field. A secret field keeps only
// its name (Secret true, no values), so password material and tokens never
// reach the trail. Old is absent on create, New on delete or clear.
type FieldChange struct {
	Name   string          `json:"name"`
	Old    json.RawMessage `json:"old,omitempty"`
	New    json.RawMessage `json:"new,omitempty"`
	Secret bool            `json:"secret,omitempty"`
}

// Change is one mutation of one entity type inside the operation. A single-row
// change carries the id; a bulk update or delete (Where-based) does not.
type Change struct {
	Entity string        `json:"entity"`
	ID     string        `json:"id,omitempty"`
	Op     string        `json:"op"`
	Fields []FieldChange `json:"fields,omitempty"`
}

// Recorder collects the changes of one operation. It is safe for concurrent use
// because a handler may fan out inside the operation.
type Recorder struct {
	mu      sync.Mutex
	changes []Change
}

type recorderKey struct{}

// WithRecorder returns a context in which row changes are noted, and the
// Recorder that holds them.
func WithRecorder(ctx context.Context) (context.Context, *Recorder) {
	r := &Recorder{
		mu:      sync.Mutex{},
		changes: nil,
	}

	return context.WithValue(ctx, recorderKey{}, r), r
}

func recorderFrom(ctx context.Context) *Recorder {
	r, _ := ctx.Value(recorderKey{}).(*Recorder)

	return r
}

// Changes returns what was noted so far, in order.
func (r *Recorder) Changes() []Change {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]Change(nil), r.changes...)
}

func (r *Recorder) add(c Change) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.changes = append(r.changes, c)
}

// primary is the entity the record is filed under: the first single-row change.
// An operation that deletes children by predicate and then the post therefore
// files under the post, and an operation that changed no row files under none.
// It returns the entity name and the row id.
func primary(changes []Change) (string, string) {
	for _, c := range changes {
		if c.ID != "" {
			return c.Entity, c.ID
		}
	}

	return "", ""
}

// Write inserts the Audit Record for an operation through client, which must be
// the client of the operation's transaction: the row then commits or rolls back
// with the change it describes. An error means the record could not be written,
// and the caller must not commit.
func Write(ctx context.Context, client *ent.Client, e Entry, rec *Recorder) error {
	var changes []Change
	if rec != nil {
		changes = rec.Changes()
	}

	occurredAt := e.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}

	create := client.AuditRecord.Create().
		SetOccurredAt(occurredAt).
		SetChannel(string(e.Source.Channel)).
		SetOperation(e.Operation).
		SetOutcome(string(e.Outcome)).
		SetNillableActorUserID(e.Actor).
		SetNillableConnectedAppID(e.Source.ConnectedAppID).
		SetNillableTokenID(e.Source.TokenID)
	if e.RequestID != "" {
		create.SetRequestID(e.RequestID)
	}

	if e.Address != "" {
		create.SetAddress(e.Address)
	}

	if e.UserAgent != "" {
		create.SetUserAgent(e.UserAgent)
	}

	if e.Locale != "" {
		create.SetLocale(e.Locale)
	}

	entity, id := e.EntityType, e.EntityID
	if entity == "" {
		entity, id = primary(changes)
	}

	if entity != "" {
		create.SetEntityType(entity)
	}

	if id != "" {
		create.SetEntityID(id)
	}

	if e.Reads != nil {
		raw, err := json.Marshal(e.Reads)
		if err != nil {
			return errs.Wrapf(err, "encode audit reads")
		}

		create.SetReads(raw)
	}

	if len(changes) > 0 {
		raw, err := json.Marshal(changes)
		if err != nil {
			return errs.Wrapf(err, "encode audit changes")
		}

		create.SetChanges(raw)
	}

	if err := create.Exec(ctx); err != nil {
		return errs.With("operation", e.Operation).Wrapf(err, "write audit record")
	}

	return nil
}
