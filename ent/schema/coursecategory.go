package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// CourseCategory maps the legacy `language_categories` table (the domain
// concept is a course category; `language` survives only in the table name for
// backward compat). Unlike the read-only catalog mappings, admin CRUD writes
// this table, so the Rails-owned `created_at`/`updated_at` (NOT NULL, no DB
// default) are supplied by ent on insert/update.
type CourseCategory struct {
	ent.Schema
}

func (CourseCategory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_categories"},
	}
}

func (CourseCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").Optional().Nillable(),
		field.String("name").Optional().Nillable(),
		field.String("header").Optional().Nillable(),
		field.String("description").Optional().Nillable(),
		field.String("locale").Optional().Nillable(),
	}
}

func (CourseCategory) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
