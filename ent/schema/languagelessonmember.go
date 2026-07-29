package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LanguageLessonMember maps the legacy `language_lesson_members` table: a user's
// participation in a single lesson. The admin surface is read-only (a list), so
// only the columns the list serializer reads are declared. The *_id columns are
// NOT NULL FKs in the baseline (value fields); state/messages_count are
// nullable. No @Table annotation: the struct name's snake-plural already matches
// the table.
type LanguageLessonMember struct {
	ent.Schema
}

func (LanguageLessonMember) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("language_id"),
		field.Int("lesson_id"),
		field.String("state").Optional().Nillable(),
		field.Int("messages_count").Optional().Nillable(),
		field.Time("created_at").Immutable(),
	}
}

func (LanguageLessonMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("course", Course.Type).
			Field("language_id").
			Unique().
			Required(),
		edge.To("lesson", LanguageLesson.Type).
			Field("lesson_id").
			Unique().
			Required(),
	}
}
