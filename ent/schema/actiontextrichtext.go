package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// ActionTextRichText maps the legacy `action_text_rich_texts` table — Rails
// ActionText's single polymorphic store for every rich-text body in the app. A
// row is keyed by (record_type, record_id, name); a BlogPost's body is the row
// with record_type='BlogPost', name='rich_body'. ent has no polymorphic edge, so
// the handler queries this schema directly filtered by those columns.
//
// Only the columns the read side needs are mapped; `body` is nullable (a record
// may have no rich text yet). The @Table annotation pins the exact legacy name.
type ActionTextRichText struct {
	ent.Schema
}

func (ActionTextRichText) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "action_text_rich_texts"},
	}
}

func (ActionTextRichText) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Int("record_id"),
		field.String("record_type"),
		field.Text("body").Optional().Nillable(),
	}
}
