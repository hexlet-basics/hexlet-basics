package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LandingPage maps the legacy `language_landing_pages` table — a localized
// catalog entry for a Course. Read-only mapping over the Rails-owned table;
// nullable columns mirror the schema with Optional().Nillable().
type LandingPage struct {
	ent.Schema
}

func (LandingPage) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_landing_pages"},
	}
}

func (LandingPage) Fields() []ent.Field {
	return []ent.Field{
		field.Int("language_id"),
		field.String("slug").Optional().Nillable(),
		field.String("header").Optional().Nillable(),
		field.String("name").Optional().Nillable(),
		field.String("locale").Optional().Nillable(),
		field.Bool("listed").Optional().Nillable(),
	}
}

func (LandingPage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("course", Course.Type).
			Ref("landing_pages").
			Field("language_id").
			Unique().
			Required(),
	}
}
