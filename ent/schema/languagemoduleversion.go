package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// LanguageModuleVersion maps the legacy `language_module_versions` table: a
// module's immutable snapshot within one course version, carrying that build's
// ordering. Written only by the exercise loader. All *_id columns are NOT NULL
// FKs in the baseline, hence plain value fields. `order` is nullable. The struct
// name's snake-plural already matches the table.
type LanguageModuleVersion struct {
	ent.Schema
}

func (LanguageModuleVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("order").Optional().Nillable(),
		field.Int("language_id"),
		field.Int("language_version_id"),
		field.Int("module_id"),
	}
}

func (LanguageModuleVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
