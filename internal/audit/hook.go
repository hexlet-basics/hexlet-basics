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
func ownedFields() []string {
	return []string{"created_at", "updated_at"}
}

// RowLoader reads one row of any schema type by id. *ent.Client satisfies it
// through the generated AuditRow (ent/template/auditrow.tmpl).
type RowLoader interface {
	AuditRow(ctx context.Context, typeName string, id any) (map[string]any, error)
}

// Hook notes every row change made while a Recorder is in the context. It is
// installed on each client the store builds, so no handler or module has audit
// code of its own: whatever an operation writes through ent is seen here.
//
// An update of one row reads the old values first (ent only allows that for
// UpdateOne), so it can say "name: A to B" and leave out the fields the input
// set to what they already were. A delete of one row loads the row through
// loader before deleting it, with the context of the deletion and so inside the
// same transaction, and records every field it held as an old value (secrets by
// name only): what a delete destroyed is what an investigator needs. A bulk
// update or delete (a Where) has no single row to load and records no values.
func Hook(loader RowLoader) ent.Hook {
	return func(next ent.Mutator) ent.Mutator { return hook(loader, next) }
}

func hook(loader RowLoader, next ent.Mutator) ent.Mutator {
	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		rec := recorderFrom(ctx)
		// The record being written must not describe itself.
		if rec == nil || m.Type() == ent.TypeAuditRecord {
			return next.Mutate(ctx, m)
		}

		change := Change{
			Entity: m.Type(),
			Op:     opName(m.Op()),
			ID:     "",
			Fields: nil,
		}
		secret := ent.AuditSecretFields[m.Type()]

		var old map[string]any
		if m.Op().Is(ent.OpUpdateOne) {
			old = make(map[string]any, len(m.Fields()))
			for _, name := range slices.Concat(m.Fields(), m.ClearedFields()) {
				v, err := m.OldField(ctx, name)
				if err != nil {
					return nil, errs.With("type", m.Type(), "field", name).Wrapf(err, "audit: read old value")
				}

				old[name] = v
			}
		}

		var deleted []FieldChange

		missing := false

		if m.Op().Is(ent.OpDeleteOne) {
			if key := entityKey(m); key != nil {
				row, err := loader.AuditRow(ctx, m.Type(), key)
				switch {
				case ent.IsNotFound(err):
					// Nothing to keep and nothing deleted: ent reports the
					// missing row after the mutation, and no change happened.
					missing = true
				case err != nil:
					return nil, errs.With("type", m.Type(), "id", key).Wrapf(err, "audit: read before delete")
				default:
					deleted = deletedFields(row, secret)
				}
			}
		}

		value, err := next.Mutate(ctx, m)
		if err != nil {
			return value, errs.Wrapf(err, "run mutation")
		}

		change.ID = entityID(m)
		if m.Op().Is(ent.OpDelete | ent.OpDeleteOne) {
			change.Fields = deleted
		} else {
			change.Fields = fieldChanges(m, old, secret)
		}

		if !missing {
			rec.add(change)
		}

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

// entityKey is the id of a single-row mutation in the Go type of the schema's
// id, nil for a bulk one. ent's typed mutations expose ID(); both id kinds the
// schema uses are handled.
func entityKey(m ent.Mutation) any {
	switch m := m.(type) {
	case interface{ ID() (int, bool) }:
		if id, ok := m.ID(); ok {
			return id
		}
	case interface{ ID() (string, bool) }:
		if id, ok := m.ID(); ok {
			return id
		}
	}

	return nil
}

// entityID is entityKey as text, empty for a bulk mutation.
func entityID(m ent.Mutation) string {
	switch key := entityKey(m).(type) {
	case int:
		return strconv.Itoa(key)
	case string:
		return key
	}

	return ""
}

// deletedFields lists what a deleted row held as old values, ordered by name. A
// secret field keeps only its name, ent's own timestamps are skipped as in an
// update, and a column that was empty says nothing worth keeping.
func deletedFields(row map[string]any, secret []string) []FieldChange {
	names := make([]string, 0, len(row))
	for name := range row {
		if !slices.Contains(ownedFields(), name) {
			names = append(names, name)
		}
	}

	slices.Sort(names)

	var out []FieldChange

	for _, name := range names {
		value := toJSON(row[name])
		if value == nil {
			continue
		}

		if slices.Contains(secret, name) {
			out = append(out, FieldChange{
				Name:   name,
				Secret: true,
				Old:    nil,
				New:    nil,
			})

			continue
		}

		out = append(out, FieldChange{
			Name:   name,
			Old:    value,
			New:    nil,
			Secret: false,
		})
	}

	return out
}

// fieldChanges lists what the mutation set or cleared. old is nil for a create
// and for a bulk update, where there is no single row to read.
func fieldChanges(m ent.Mutation, old map[string]any, secret []string) []FieldChange {
	var out []FieldChange

	for _, name := range slices.Concat(m.Fields(), m.ClearedFields()) {
		if slices.Contains(ownedFields(), name) {
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

			out = append(out, FieldChange{
				Name:   name,
				Secret: true,
				Old:    nil,
				New:    nil,
			})

			continue
		}

		if old != nil && sameJSON(oldValue, newValue) {
			continue
		}

		fc := FieldChange{
			Name:   name,
			Old:    nil,
			New:    nil,
			Secret: false,
		}
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
