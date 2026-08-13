package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CourseLessonTranslation maps the legacy `language_lesson_version_infos`
// table: the localized (per-locale) name/description/theory of one lesson
// version. The admin lesson list is built AROUND this table (legacy
// `Language::Lesson::Version::Info.current.with_locale`) — one row per lesson
// per locale for the course's current version — so name/description come from
// here while the lesson slug is joined in separately.
//
// Only the columns the admin lists read are declared. The *_id columns are NOT
// NULL FKs in the baseline, hence plain value fields (no Nillable); the text
// columns are nullable. The Table annotation pins the legacy table, which the
// renamed type no longer derives.
type CourseLessonTranslation struct {
	ent.Schema
}

func (CourseLessonTranslation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_lesson_version_infos"},
	}
}

func (CourseLessonTranslation) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Nillable(),
		field.String("description").Optional().Nillable(),
		field.String("locale").Optional().Nillable(),
		field.Int("course_id").StorageKey("language_id"),
		field.Int("course_lesson_id").StorageKey("language_lesson_id"),
		field.Int("course_version_id").StorageKey("language_version_id"),
		// Written by the loader from the lesson's <locale>/ dir: README.md →
		// theory, EXERCISE.md → instructions, data.yml tips/definitions serialized
		// as YAML arrays (Rails `serialize type: Array` compatibility). version_id
		// FKs the CourseLessonVersion this info belongs to (NOT NULL).
		field.String("theory").Optional().Nillable(),
		field.String("instructions").Optional().Nillable(),
		field.String("tips").Optional().Nillable(),
		field.String("definitions").Optional().Nillable(),
		field.Int("version_id"),
	}
}

func (CourseLessonTranslation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("lesson", CourseLesson.Type).
			Field("course_lesson_id").
			Unique().
			Required(),
		edge.To("course_version", CourseVersion.Type).
			Field("course_version_id").
			Unique().
			Required(),
		// The lesson version this row localizes. It is what carries the Position
		// (`natural_order`) the reads order by: the loader writes that number onto
		// the version row and never onto the lesson, so the lesson's own column is
		// a legacy leftover nothing maintains.
		edge.To("version", CourseLessonVersion.Type).
			Field("version_id").
			Unique().
			Required(),
	}
}

func (CourseLessonTranslation) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
