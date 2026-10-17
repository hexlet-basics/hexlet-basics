package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	"hexletbasics/ent"
)

// ownedFields are maintained by ent itself (the timestamps mixin), so a change
// of them says nothing the operation did.
var ownedFields = []string{"created_at", "updated_at"}

// Hook notes every row change made while a Recorder is in the context. It is
// installed on each client the store builds, so no handler or module has audit
// code of its own: whatever an operation writes through ent is seen here.
//
// An update of one row reads the old values first (ent only allows that for
// UpdateOne), so it can say "name: A to B" and leave out the fields the input
// set to what they already were. A delete of one row records the entity and id
// but not the old values: ent offers no generic read of the row being deleted.
func Hook(next ent.Mutator) ent.Mutator {
	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		rec := recorderFrom(ctx)
		// The record being written must not describe itself.
		if rec == nil || m.Type() == ent.TypeAuditRecord {
			return next.Mutate(ctx, m)
		}

		change := Change{Entity: m.Type(), Op: opName(m.Op())}
		secret := ent.AuditSecretFields[m.Type()]

		var old map[string]any
		if m.Op().Is(ent.OpUpdateOne) {
			old = make(map[string]any, len(m.Fields()))
			for _, name := range slices.Concat(m.Fields(), m.ClearedFields()) {
				v, err := m.OldField(ctx, name)
				if err != nil {
					return nil, fmt.Errorf("audit: read old %s.%s: %w", m.Type(), name, err)
				}
				old[name] = v
			}
		}

		value, err := next.Mutate(ctx, m)
		if err != nil {
			return value, err
		}

		change.ID = entityID(m)
		if !m.Op().Is(ent.OpDelete | ent.OpDeleteOne) {
			change.Fields = fieldChanges(m, old, secret)
		}
		rec.add(change)
		return value, nil
	})
}

func opName(op ent.Op) string {
	switch {
	case op.Is(ent.OpCreate):
		return "create"
	case op.Is(ent.OpDelete | ent.OpDeleteOne):
		return "delete"
	default:
		return "update"
	}
}

// entityID is the id of a single-row mutation as text, empty for a bulk one.
// ent's typed mutations expose ID(); both id kinds the schema uses are handled.
func entityID(m ent.Mutation) string {
	switch m := m.(type) {
	case interface{ ID() (int, bool) }:
		if id, ok := m.ID(); ok {
			return strconv.Itoa(id)
		}
	case interface{ ID() (string, bool) }:
		if id, ok := m.ID(); ok {
			return id
		}
	}
	return ""
}

// fieldChanges lists what the mutation set or cleared. old is nil for a create
// and for a bulk update, where there is no single row to read.
func fieldChanges(m ent.Mutation, old map[string]any, secret []string) []FieldChange {
	var out []FieldChange
	for _, name := range slices.Concat(m.Fields(), m.ClearedFields()) {
		if slices.Contains(ownedFields, name) {
			continue
		}
		var newValue any
		if v, ok := m.Field(name); ok {
			newValue = v
		}
		oldValue := old[name]

		if slices.Contains(secret, name) {
			// Only the name: whether it changed is the fact worth keeping.
			if old != nil && sameJSON(oldValue, newValue) {
				continue
			}
			out = append(out, FieldChange{Name: name, Secret: true})
			continue
		}
		if old != nil && sameJSON(oldValue, newValue) {
			continue
		}
		fc := FieldChange{Name: name}
		if old != nil {
			fc.Old = toJSON(oldValue)
		}
		fc.New = toJSON(newValue)
		out = append(out, fc)
	}
	return out
}

// toJSON encodes a value for the trail; nil becomes absent rather than "null".
func toJSON(v any) json.RawMessage {
	if v == nil {
		return nil
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(strconv.Quote(fmt.Sprint(v)))
	}
	if bytes.Equal(raw, []byte("null")) {
		return nil
	}
	return raw
}

// sameJSON compares by encoding, so a *string and a string holding the same
// text, or a nil pointer and a cleared field, are equal.
func sameJSON(a, b any) bool {
	return bytes.Equal(toJSON(a), toJSON(b))
}
