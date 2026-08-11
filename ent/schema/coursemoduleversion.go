package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// CourseModuleVersion maps the legacy `language_module_versions` table: a
// module's immutable snapshot within one course version, carrying that build's
// ordering. Written only by the exercise loader. All *_id columns are NOT NULL
// FKs in the baseline, hence plain value fields. `order` is nullable. The Table
// annotation pins the legacy table, which the renamed type no longer derives.
type CourseModuleVersion struct {
	ent.Schema
}

func (CourseModuleVersion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_module_versions"},
	}
}

func (CourseModuleVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("order").Optional().Nillable(),
		field.Int("course_id").StorageKey("language_id"),
		field.Int("course_version_id").StorageKey("language_version_id"),
		field.Int("module_id"),
	}
}

func (CourseModuleVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
