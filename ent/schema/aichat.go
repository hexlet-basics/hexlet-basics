package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AiChat maps the legacy `ai_chats` table: one in-lesson assistant conversation
// per (user, lesson progress). Only the product columns are declared — the
// RubyLLM bookkeeping column (`ai_model_id`) stays in the table but out of the
// Go domain: which model serves a chat is configuration, not a row attribute
// (see the assistant port design). The assistant will write this table, so the
// Rails-owned timestamps come from TimestampsMixin. No @Table annotation: the
// struct name's snake-plural already matches the table.
type AiChat struct {
	ent.Schema
}

func (AiChat) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("language_lesson_member_id"),
	}
}

func (AiChat) Edges() []ent.Edge {
	return []ent.Edge{
		// belongs-to edges own their NOT NULL FK columns; the lesson progress
		// carries the course/lesson identity the admin read models enrich
		// through. The edge name is ours to choose — the pinned Field keeps the
		// legacy column, so renaming it moves no storage identifier.
		edge.To("lesson_progress", LessonProgress.Type).
			Field("language_lesson_member_id").
			Unique().
			Required(),
		edge.To("user", User.Type).
			Field("user_id").
			Unique().
			Required(),
	}
}

func (AiChat) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
