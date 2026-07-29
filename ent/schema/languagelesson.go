package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// LanguageLesson maps the legacy `language_lessons` table (a single lesson in a
// course). The admin surface reads it to enrich other lists (a review's or
// member's lesson slug and natural order); the exercise loader also WRITES it,
// upserting one stable row per (language, slug) so learner progress
// (language_lesson_members, keyed to the lesson id) survives rebuilds — the
// per-build code/order lives on LanguageLessonVersion, not here.
//
// Only the columns the loader sets or the reads select are declared; atlas owns
// the schema, so the denormalized columns the loader does NOT write (order,
// original_code, prepared_code, test_code, path_to_code) are simply never
// selected. The struct name's snake-plural already matches the table.
//
// slug/language_id/module_id are nullable in the baseline even though the app
// validates slug presence, hence Optional().Nillable().
type LanguageLesson struct {
	ent.Schema
}

func (LanguageLesson) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").Optional().Nillable(),
		field.Int("natural_order").Optional().Nillable(),
		// Written by the loader: the owning course and the module the lesson
		// currently belongs to (a lesson can move modules across rebuilds).
		field.Int("language_id").Optional().Nillable(),
		field.Int("module_id").Optional().Nillable(),
		field.String("state").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (LanguageLesson) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("infos", LanguageLessonVersionInfo.Type).
			Ref("lesson"),
	}
}

func (LanguageLesson) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("language_id", "slug").Unique(),
	}
}
