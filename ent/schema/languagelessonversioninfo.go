package schema

import (
	"entgo.io/ent"
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
	}
}
