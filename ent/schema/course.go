package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Course maps the legacy `languages` table (the domain concept is Course).
// The table is owned by the legacy Rails schema; ent only reads it, so nullable
// columns are mirrored exactly with Optional().Nillable() to avoid scan errors.
type Course struct {
	ent.Schema
}

func (Course) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "languages"},
	}
}

func (Course) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").Optional().Nillable(),
		field.String("name").Optional().Nillable(),
		field.String("learn_as").Optional().Nillable(),
		field.String("progress").Optional().Nillable(),
		field.Int("members_count"),
		field.Int("lessons_count"),
		field.Int("category_id").Optional().Nillable(),
		field.Int("order").Optional().Nillable(),
	}
}

func (Course) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("landing_pages", LandingPage.Type),
	}
}
