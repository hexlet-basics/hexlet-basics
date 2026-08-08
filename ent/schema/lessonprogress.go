package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LessonProgress maps the legacy `language_lesson_members` table: a user's
// participation in a single lesson. The admin surface is read-only (a list), so
// only the columns the list serializer reads are declared. The *_id columns are
// NOT NULL FKs in the baseline (value fields); state/messages_count are
// nullable. The Table annotation pins the legacy table, which the renamed type
// no longer derives.
type LessonProgress struct {
	ent.Schema
}

func (LessonProgress) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_lesson_members"},
	}
}

func (LessonProgress) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("language_id"),
		field.Int("lesson_id"),
		field.String("state").Optional().Nillable(),
		field.Int("messages_count").Optional().Nillable(),
		field.Time("created_at").Immutable(),
	}
}

func (LessonProgress) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("course", Course.Type).
			Field("language_id").
			Unique().
			Required(),
		edge.To("lesson", CourseLesson.Type).
			Field("lesson_id").
			Unique().
			Required(),
	}
}
