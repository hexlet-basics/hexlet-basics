package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Tag maps the legacy acts_as_taggable_on `tags` table. The Go stack only takes
// a tag off a user (the lead form clears `should_be_lead`), so just the name
// and the counter cache that removal keeps in step are mapped; the Rails
// timestamps stay unmapped because nothing here inserts a tag, and legacy's
// counter update did not touch `updated_at`. Atlas owns the table; the name is
// ent's default plural.
type Tag struct {
	ent.Schema
}

func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Nillable(),
		// acts_as_taggable_on's counter cache (`tags_counter`, on by default).
		field.Int("taggings_count").Optional().Nillable(),
	}
}
