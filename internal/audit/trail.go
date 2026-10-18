package audit

import (
	"errors"
	"time"

	"hexletbasics/ent"
	"hexletbasics/ent/auditrecord"
)

// The audit table is partitioned by month on occurred_at (ADR-0020, #846), so a
// query that does not bound occurred_at visits every partition and every index
// of two years of records. The listing therefore always carries a window, and
// these are the numbers it is held to.
const (
	// DefaultWindow is how far back the listing looks when the caller names no
	// start: a month, enough to answer "what happened recently" while touching
	// one or two partitions.
	DefaultWindow = 30 * 24 * time.Hour

	// MaxWindow caps the span of one query: three 31-day months at most, so at
	// most four partitions are visited whatever the filters. Looking further back
	// is done with several narrower windows, not with one wide scan.
	MaxWindow = 92 * 24 * time.Hour
)

var (
	// ErrWindowReversed is a window that ends before it starts.
	ErrWindowReversed = errors.New("the end of the window is before its start")
	// ErrWindowTooWide is a window longer than MaxWindow.
	ErrWindowTooWide = errors.New("the window is longer than the audit trail allows in one query")
	// ErrEntityIDWithoutType is an entity id given without its type. The index is
	// (entity_type, entity_id, occurred_at); an id alone could not use it.
	ErrEntityIDWithoutType = errors.New("an entity id needs its entity type")
)

// Filter is what the listing narrows the trail by. The zero value of each
// pointer means "no filter on it".
type Filter struct {
	ActorUserID *int
	Channel     *Channel
	Operation   *string
	EntityType  *string
	EntityID    *string
	// From is inclusive and To exclusive, both optional: To defaults to now and
	// From to DefaultWindow before To (see Window).
	From *time.Time
	To   *time.Time
}

// Window resolves the half-open [from, to) interval the query is bounded by.
func (f Filter) Window(now time.Time) (time.Time, time.Time, error) {
	to := now
	if f.To != nil {
		to = *f.To
	}
	from := to.Add(-DefaultWindow)
	if f.From != nil {
		from = *f.From
	}
	if !from.Before(to) {
		return time.Time{}, time.Time{}, ErrWindowReversed
	}
	if to.Sub(from) > MaxWindow {
		return time.Time{}, time.Time{}, ErrWindowTooWide
	}
	return from, to, nil
}

// Query builds the listing query: newest first, always bounded by the window,
// and narrowed by the filters. Rows are ordered by (occurred_at, id) descending;
// the id breaks the tie between records written in the same instant so a page
// boundary never repeats or skips a row. The window is what lets PostgreSQL
// prune partitions, and the equality filters line up with the existing indexes:
// occurred_at, (actor_user_id, occurred_at) and (entity_type, entity_id,
// occurred_at). Every call builds a fresh query so the count and the page read
// from independent builders.
func (f Filter) Query(db *ent.Client, now time.Time) (*ent.AuditRecordQuery, error) {
	if f.EntityID != nil && f.EntityType == nil {
		return nil, ErrEntityIDWithoutType
	}
	from, to, err := f.Window(now)
	if err != nil {
		return nil, err
	}
	q := db.AuditRecord.Query().
		Where(auditrecord.OccurredAtGTE(from), auditrecord.OccurredAtLT(to)).
		Order(ent.Desc(auditrecord.FieldOccurredAt), ent.Desc(auditrecord.FieldID))
	if f.ActorUserID != nil {
		q = q.Where(auditrecord.ActorUserID(*f.ActorUserID))
	}
	if f.Channel != nil {
		q = q.Where(auditrecord.Channel(string(*f.Channel)))
	}
	if f.Operation != nil {
		q = q.Where(auditrecord.Operation(*f.Operation))
	}
	if f.EntityType != nil {
		q = q.Where(auditrecord.EntityType(*f.EntityType))
	}
	if f.EntityID != nil {
		q = q.Where(auditrecord.EntityID(*f.EntityID))
	}
	return q, nil
}
