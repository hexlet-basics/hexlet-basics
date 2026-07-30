package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Review maps the legacy `reviews` table — a student testimonial for a course.
// The admin surface writes this table, so the Rails-owned `created_at`/
// `updated_at` (NOT NULL, no DB default) are supplied by ent. The table name
// already matches ent's default plural of `Review`, so no @Table annotation.
//
// A review belongs to a course (the FK column is `language_id`) and a user;
// both FKs are NOT NULL. The admin API embeds the full course and user, so the
// edges are loaded (WithCourse/WithUser) and converted at the API boundary.
// `first_name`/`last_name` are the review's own display name (the reviewer may
// differ from the account), distinct from the associated user's name.
type Review struct {
	ent.Schema
}

func (Review) Fields() []ent.Field {
	return []ent.Field{
		field.Int("language_id"),
		field.Int("user_id"),
		field.String("body").Optional().Nillable(),
		field.String("first_name").Optional().Nillable(),
		field.String("last_name").Optional().Nillable(),
		field.String("locale").Optional().Nillable(),
		field.String("state").Optional().Nillable(),
		field.Bool("pinned").Optional().Nillable(),
	}
}

func (Review) Edges() []ent.Edge {
	return []ent.Edge{
		// belongs-to: the FK lives on this (reviews) table, so the edge owns the
		// field. `course` binds the legacy `language_id` column to a Course.
		edge.To("course", Course.Type).
			Field("language_id").
			Unique().
			Required(),
		edge.To("user", User.Type).
			Field("user_id").
			Unique().
			Required(),
	}
}

func (Review) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
