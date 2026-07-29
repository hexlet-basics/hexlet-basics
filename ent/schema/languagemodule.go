package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// LanguageModule maps the legacy `language_modules` table: a stable module (a
// group of lessons) within a course. Like LanguageLesson it is upserted by
// (language, slug) so its identity survives rebuilds; the per-build ordering and
// localized name/description live on LanguageModuleVersion and
// LanguageModuleVersionInfo. Only the loader-written columns are declared (the
// baseline `upload_id` is unused). The struct name's snake-plural already matches
// the table.
type LanguageModule struct {
	ent.Schema
}

func (LanguageModule) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").Optional().Nillable(),
		field.Int("language_id").Optional().Nillable(),
		field.Int("order").Optional().Nillable(),
		field.String("state").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (LanguageModule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("language_id", "slug").Unique(),
	}
}
