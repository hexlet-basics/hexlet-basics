package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CourseLessonReview maps the legacy `language_lesson_reviews` table: an
// AI-generated summary of student questions for one lesson version, per locale.
// The admin list filters to the current locale and to non-empty summaries
// (legacy `where(locale:).with_summary`); the review worker upserts rows by
// (course, lesson, locale), so the Rails-owned timestamps come from
// TimestampsMixin. All columns the serializer reads are NOT NULL in the
// baseline, so they are plain value fields. The Table annotation pins the legacy
// table, which the renamed type no longer derives.
type CourseLessonReview struct {
	ent.Schema
}

func (CourseLessonReview) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_lesson_reviews"},
	}
}

func (CourseLessonReview) Fields() []ent.Field {
	return []ent.Field{
		field.String("locale"),
		field.String("summary"),
		field.Int("language_id"),
		field.Int("language_lesson_id"),
		field.Int("language_lesson_version_id"),
		field.Int("language_lesson_version_info_id"),
	}
}

func (CourseLessonReview) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}

func (CourseLessonReview) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("course", Course.Type).
			Field("language_id").
			Unique().
			Required(),
		edge.To("lesson", CourseLesson.Type).
			Field("language_lesson_id").
			Unique().
			Required(),
	}
}
