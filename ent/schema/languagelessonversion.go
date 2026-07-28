package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// LanguageLessonVersion maps the legacy `language_lesson_versions` table: a
// lesson's immutable snapshot within one course version. It carries the build's
// code (original/prepared/test), the container path used by the runtime exercise
// runner, and the global `natural_order` that sequences lessons across all
// modules. Written only by the exercise loader. All *_id columns are NOT NULL FKs
// in the baseline (plain value fields); order/natural_order and the code columns
// are nullable. The struct name's snake-plural already matches the table.
type LanguageLessonVersion struct {
	ent.Schema
}

func (LanguageLessonVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("natural_order").Optional().Nillable(),
		field.Int("order").Optional().Nillable(),
		field.String("original_code").Optional().Nillable(),
		field.String("prepared_code").Optional().Nillable(),
		field.String("test_code").Optional().Nillable(),
		field.String("path_to_code").Optional().Nillable(),
		field.Int("language_id"),
		field.Int("language_version_id"),
		field.Int("lesson_id"),
		field.Int("module_version_id"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
