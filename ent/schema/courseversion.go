package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// CourseVersion maps the legacy `language_versions` table. The domain concept is
// a Course version; `language` survives only in the table name. It is read-only
// here (surfaced as Course.currentVersion), so only the columns the API exposes
// are declared — ent scans just these, ignoring the build-metadata columns.
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
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}
