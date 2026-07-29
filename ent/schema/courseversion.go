package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CourseVersion maps the legacy `language_versions` table. The domain concept is
// a Course version; `language` survives only in the table name. It is surfaced
// read-only as Course.currentVersion, but the exercise loader (ADR: course build
// auto-promotes to live) also WRITES it: each build creates a fresh version row,
// fills the build-metadata columns from the repo's spec.yml, and drives the
// created→building→built/failed state column. The write-only columns are declared
// here so the loader can set them; the read converter maps only result/state.
type CourseVersion struct {
	ent.Schema
}

func (CourseVersion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_versions"},
	}
}

func (CourseVersion) Fields() []ent.Field {
	return []ent.Field{
		field.String("result").Optional().Nillable(),
		field.String("state").Optional().Nillable(),
		// Build metadata copied from spec.yml's `language:` map (all nullable in
		// the baseline). These drive the runtime exercise runner (docker_image,
		// exercise_filename/exercise_test_filename tell it which files to run) and
		// mirror the version-level snapshot of the course's identity.
		field.String("name").Optional().Nillable(),
		field.String("progress").Optional().Nillable(),
		field.String("learn_as").Optional().Nillable(),
		field.String("extension").Optional().Nillable(),
		field.String("docker_image").Optional().Nillable(),
		field.String("exercise_filename").Optional().Nillable(),
		field.String("exercise_test_filename").Optional().Nillable(),
		// counter_culture-maintained in Rails; the column is NOT NULL DEFAULT 0, so
		// a fresh version starts at 0 and the loader bumps it per lesson version.
		field.Int("lessons_count").Default(0),
		// language_id is the owning course (NOT NULL FK). A version is always built
		// for a specific course, so it is a required value field, not nillable.
		field.Int("language_id"),
		field.Time("created_at").Default(time.Now).Immutable(),
		// Rails-owned timestamp (NOT NULL, no DB default); supplied by ent now that
		// the loader writes this table.
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (CourseVersion) Edges() []ent.Edge {
	return []ent.Edge{
		// A version can be selected as current by a course through the FK stored
		// on languages.current_version_id. The inverse edge lets read queries
		// express "info belongs to a completed course's current version"
		// without first materializing version ids in the handler.
		edge.From("current_courses", Course.Type).
			Ref("current_version"),
	}
}
