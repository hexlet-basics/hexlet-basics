package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// LanguageModuleVersionInfo maps the legacy `language_module_version_infos`
// table: the localized (per-locale) name/description of one module version,
// sourced from the module dir's `description.<locale>.yml`. Written only by the
// exercise loader. `version_id` FKs the LanguageModuleVersion (NOT NULL, as are
// language_id/language_version_id); the text columns are nullable. The struct
// name's snake-plural already matches the table.
type LanguageModuleVersionInfo struct {
	ent.Schema
}

func (LanguageModuleVersionInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Nillable(),
		field.String("description").Optional().Nillable(),
		field.String("locale").Optional().Nillable(),
		field.Int("language_id"),
		field.Int("language_version_id"),
		field.Int("version_id"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
