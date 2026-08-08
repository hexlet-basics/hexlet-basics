package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// CourseModuleTranslation maps the legacy `language_module_version_infos`
// table: the localized (per-locale) name/description of one module version,
// sourced from the module dir's `description.<locale>.yml`. Written only by the
// exercise loader. `version_id` FKs the CourseModuleVersion (NOT NULL, as are
// language_id/language_version_id); the text columns are nullable. The Table
// annotation pins the legacy table, which the renamed type no longer derives.
type CourseModuleTranslation struct {
	ent.Schema
}

func (CourseModuleTranslation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_module_version_infos"},
	}
}

func (CourseModuleTranslation) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Nillable(),
		field.String("description").Optional().Nillable(),
		field.String("locale").Optional().Nillable(),
		field.Int("language_id"),
		field.Int("language_version_id"),
		field.Int("version_id"),
	}
}

func (CourseModuleTranslation) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
