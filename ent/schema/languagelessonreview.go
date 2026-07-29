package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LanguageLessonReview maps the legacy `language_lesson_reviews` table: an
// AI-generated summary of student questions for one lesson version, per locale.
// The admin surface is read-only (a list filtered to the current locale and to
// non-empty summaries — legacy `where(locale:).with_summary`). All columns the
// serializer reads are NOT NULL in the baseline, so they are plain value fields.
// No @Table annotation: the struct name's snake-plural already matches the
// table.
type LanguageLessonReview struct {
	ent.Schema
}

func (LanguageLessonReview) Fields() []ent.Field {
	return []ent.Field{
		field.String("locale"),
		field.String("summary"),
		field.Int("language_id"),
		field.Int("language_lesson_id"),
		field.Int("language_lesson_version_id"),
		field.Int("language_lesson_version_info_id"),
		field.Time("created_at").Immutable(),
	}
}

func (LanguageLessonReview) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("course", Course.Type).
			Field("language_id").
			Unique().
			Required(),
		edge.To("lesson", LanguageLesson.Type).
			Field("language_lesson_id").
			Unique().
			Required(),
	}
}
