package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LanguageLessonVersionInfo maps the legacy `language_lesson_version_infos`
// table: the localized (per-locale) name/description/theory of one lesson
// version. The admin lesson list is built AROUND this table (legacy
// `Language::Lesson::Version::Info.current.with_locale`) — one row per lesson
// per locale for the course's current version — so name/description come from
// here while the lesson slug is joined in separately.
//
// Only the columns the admin lists read are declared. The *_id columns are NOT
// NULL FKs in the baseline, hence plain value fields (no Nillable); the text
// columns are nullable. No @Table annotation: the struct name's snake-plural
// already matches the table.
type LanguageLessonVersionInfo struct {
	ent.Schema
}

func (LanguageLessonVersionInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Nillable(),
		field.String("description").Optional().Nillable(),
		field.String("locale").Optional().Nillable(),
		field.Int("language_id"),
		field.Int("language_lesson_id"),
		field.Int("language_version_id"),
		// Written by the loader from the lesson's <locale>/ dir: README.md →
		// theory, EXERCISE.md → instructions, data.yml tips/definitions serialized
		// as YAML arrays (Rails `serialize type: Array` compatibility). version_id
		// FKs the LanguageLessonVersion this info belongs to (NOT NULL).
		field.String("theory").Optional().Nillable(),
		field.String("instructions").Optional().Nillable(),
		field.String("tips").Optional().Nillable(),
		field.String("definitions").Optional().Nillable(),
		field.Int("version_id"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (LanguageLessonVersionInfo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("lesson", LanguageLesson.Type).
			Field("language_lesson_id").
			Unique().
			Required(),
		edge.To("course_version", CourseVersion.Type).
			Field("language_version_id").
			Unique().
			Required(),
	}
}
