package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Course maps the legacy `languages` table (the domain concept is Course).
// The table is owned by the legacy Rails schema; ent only reads it, so nullable
// columns are mirrored exactly with Optional().Nillable() to avoid scan errors.
type Course struct {
	ent.Schema
}

func (Course) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "languages"},
	}
}

func (Course) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").Optional().Nillable(),
		field.String("name").Optional().Nillable(),
		field.String("learn_as").Optional().Nillable(),
		field.String("progress").Optional().Nillable(),
		field.String("hexlet_program_landing_page").Optional().Nillable(),
		// Counters maintained by the app; the baseline columns are NOT NULL
		// DEFAULT 0, so a newly-created course starts at 0 without the admin
		// create path having to supply them.
		field.Int("members_count").Default(0),
		field.Int("lessons_count").Default(0),
		field.Int("category_id").Optional().Nillable(),
		field.Int("current_version_id").Optional().Nillable(),
		field.Int("order").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		// Rails-owned timestamp; supplied by ent now that admin create/update
		// writes the table (NOT NULL, no DB default).
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Course) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("landing_pages", LandingPage.Type),
		// current_version is the specific version a course points at via its
		// `current_version_id` column (belongs-to), distinct from a course's
		// full version history. The FK lives on this (languages) table, so the
		// edge owns the field.
		edge.To("current_version", CourseVersion.Type).
			Field("current_version_id").
			Unique(),
	}
}
