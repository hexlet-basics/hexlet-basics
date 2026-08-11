package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CourseLesson maps the legacy `language_lessons` table (a single lesson in a
// course). The admin surface reads it to enrich other lists (a review's or a
// lesson progress' lesson slug and natural order); the exercise loader also
// WRITES it, upserting one stable row per (language, slug) so learner progress
// (language_lesson_members, keyed to the lesson id) survives rebuilds — the
// per-build code/order lives on CourseLessonVersion, not here.
//
// Only the columns the loader sets or the reads select are declared; atlas owns
// the schema, so the denormalized columns the loader does NOT write (order,
// original_code, prepared_code, test_code, path_to_code) are simply never
// selected. The Table annotation pins the legacy table, which the renamed type
// no longer derives.
//
// slug/language_id/module_id are nullable in the baseline even though the app
// validates slug presence, hence Optional().Nillable().
type CourseLesson struct {
	ent.Schema
}

func (CourseLesson) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_lessons"},
	}
}

func (CourseLesson) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").Optional().Nillable(),
		field.Int("natural_order").Optional().Nillable(),
		// Written by the loader: the owning course and the module the lesson
		// currently belongs to (a lesson can move modules across rebuilds).
		field.Int("course_id").StorageKey("language_id").Optional().Nillable(),
		field.Int("module_id").Optional().Nillable(),
		field.String("state").Optional().Nillable(),
	}
}

func (CourseLesson) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("infos", CourseLessonTranslation.Type).
			Ref("lesson"),
	}
}

func (CourseLesson) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("course_id", "slug").Unique(),
	}
}

func (CourseLesson) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
