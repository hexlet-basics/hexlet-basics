package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CourseLessonVersion maps the legacy `language_lesson_versions` table: a
// lesson's immutable snapshot within one course version. It carries the build's
// code (original/prepared/test), the container path used by the runtime exercise
// runner, and the global `natural_order` that sequences lessons across all
// modules. Written only by the exercise loader. All *_id columns are NOT NULL FKs
// in the baseline (plain value fields); order/natural_order and the code columns
// are nullable. The Table annotation pins the legacy table, which the renamed
// type no longer derives.
type CourseLessonVersion struct {
	ent.Schema
}

func (CourseLessonVersion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_lesson_versions"},
	}
}

func (CourseLessonVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("natural_order").Optional().Nillable(),
		field.Int("order").Optional().Nillable(),
		field.String("original_code").Optional().Nillable(),
		field.String("prepared_code").Optional().Nillable(),
		field.String("test_code").Optional().Nillable(),
		field.String("path_to_code").Optional().Nillable(),
		field.Int("course_id").StorageKey("language_id"),
		field.Int("course_version_id").StorageKey("language_version_id"),
		field.Int("lesson_id"),
		field.Int("module_version_id"),
	}
}

// Edges own the existing module_version_id and lesson_id columns (no
// migration), so a read can walk module version -> lesson versions -> lesson
// slugs with eager loading instead of stitching id maps by hand.
func (CourseLessonVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("module_version", CourseModuleVersion.Type).
			Ref("lesson_versions").
			Field("module_version_id").
			Unique().
			Required(),
		edge.To("lesson", CourseLesson.Type).
			Field("lesson_id").
			Unique().
			Required(),
	}
}

func (CourseLessonVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
