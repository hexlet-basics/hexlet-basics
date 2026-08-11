package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CourseModule maps the legacy `language_modules` table: a stable module (a
// group of lessons) within a course. Like CourseLesson it is upserted by
// (language, slug) so its identity survives rebuilds; the per-build ordering and
// localized name/description live on CourseModuleVersion and
// CourseModuleTranslation. Only the loader-written columns are declared (the
// baseline `upload_id` is unused). The Table annotation pins the legacy table,
// which the renamed type no longer derives.
type CourseModule struct {
	ent.Schema
}

func (CourseModule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_modules"},
	}
}

func (CourseModule) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").Optional().Nillable(),
		field.Int("course_id").StorageKey("language_id").Optional().Nillable(),
		field.Int("order").Optional().Nillable(),
		field.String("state").Optional().Nillable(),
	}
}

func (CourseModule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("course_id", "slug").Unique(),
	}
}

func (CourseModule) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
