package schema

import (
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
		// Explicit because the schema is wider than the admin surface: only
		// the marked fields map onto CourseInput. The input's repositoryUrl
		// and coverAttachmentId are handled elsewhere (derived on read /
		// deferred until the Attachments uploader), so they map to nothing.
		AdminInput{Type: "CourseInput", Explicit: true},
	}
}

func (Course) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").Optional().Nillable().
			Annotations(AdminInputField{}),
		field.String("name").Optional().Nillable(),
		field.String("learn_as").Optional().Nillable().
			Annotations(AdminInputField{}),
		field.String("progress").Optional().Nillable().
			Annotations(AdminInputField{}),
		field.String("hexlet_program_landing_page").Optional().Nillable().
			Annotations(AdminInputField{}),
		// Counters maintained by the app; the baseline columns are NOT NULL
		// DEFAULT 0, so a newly-created course starts at 0 without the admin
		// create path having to supply them.
		field.Int("members_count").Default(0),
		field.Int("lessons_count").Default(0),
		field.Int("category_id").Optional().Nillable(),
		field.Int("current_version_id").Optional().Nillable(),
		field.Int("order").Optional().Nillable(),
		// Rails-owned timestamp; supplied by ent now that admin create/update
		// writes the table (NOT NULL, no DB default).
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

func (Course) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
