package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Tagging maps the legacy acts_as_taggable_on `taggings` table: a polymorphic
// join of a tag onto any record (`taggable_type` names the Rails class). Only
// the columns a tag removal matches on are mapped, as plain fields with no edge
// (the polymorphic target cannot be one). Atlas owns the table; the name is
// ent's default plural.
type Tagging struct {
	ent.Schema
}

func (Tagging) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tag_id").Optional().Nillable(),
		field.Int("taggable_id").Optional().Nillable(),
		field.String("taggable_type").Optional().Nillable(),
		// The tag list the tagging belongs to; `tag_list` is the "tags" context.
		field.String("context").Optional().Nillable(),
	}
}
