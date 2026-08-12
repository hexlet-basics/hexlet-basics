package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Enrollment maps the legacy `language_members` table: a learner's record of
// where they are in one Course. It carries no `progress` column — Completion is
// computed from the current Version's finished Lessons (ADR-0012), and the
// denormalized `finished_lessons_count` the legacy app maintained counts
// Lessons that later Versions may have dropped, which is why the legacy
// serializer had to clamp its percentage at 100.
//
// state is nullable in the baseline even though AASM always wrote it, so the
// conversion layer resolves a null to Started.
//
// The unique index below is the point of the schema: the legacy find-or-create
// was racy (two non-unique indexes and no constraint on the pair), so every
// later ticket in the progress epic can only insert conflict-tolerantly once
// PostgreSQL owns the invariant. The Table annotation pins the legacy table,
// which the domain type name no longer derives.
type Enrollment struct {
	ent.Schema
}

func (Enrollment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_members"},
	}
}

func (Enrollment) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("course_id").StorageKey("language_id"),
		field.String("state").Optional().Nillable(),
		field.Int("finished_lessons_count").Default(0),
	}
}

func (Enrollment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("course", Course.Type).
			Field("course_id").
			Unique().
			Required(),
		edge.To("user", User.Type).
			Field("user_id").
			Unique().
			Required(),
		edge.From("lesson_progress", LessonProgress.Type).
			Ref("enrollment"),
	}
}

// Indexes declares the (learner, course) uniqueness the migration adds. ent
// never creates it — atlas owns the schema — but declaring it here keeps the
// generated client's OnConflict helpers aware of the real constraint.
func (Enrollment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "course_id").Unique(),
	}
}

func (Enrollment) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
