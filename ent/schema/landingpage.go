package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LandingPage maps the legacy `language_landing_pages` table — a localized
// catalog entry for a Course. The public catalog reads a slim projection
// (ToCatalogItem), while the admin surface reads/writes the full row
// (ToCourseLandingPage), so all editable columns are mapped; nullable columns
// mirror the schema with Optional().Nillable(). The Rails-owned timestamps are
// supplied by ent because admin create/update writes this table.
type LandingPage struct {
	ent.Schema
}

func (LandingPage) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_landing_pages"},
		// The input's outcomesImageAttachmentId maps to nothing — the asset
		// is deferred until the Attachments uploader lands; `locale` is not
		// part of the admin input.
		AdminInput{Type: "CourseLandingPageInput"},
	}
}

func (LandingPage) Fields() []ent.Field {
	return []ent.Field{
		field.Int("course_id").StorageKey("language_id").
			Annotations(AdminInputField{Rename: "CourseId"}),
		field.String("slug").Optional().Nillable(),
		field.String("header").Optional().Nillable(),
		field.String("name").Optional().Nillable(),
		field.String("locale").Optional().Nillable().
			Annotations(AdminInputField{Skip: true}),
		field.Bool("listed").Optional().Nillable(),
		field.Bool("main").Optional().Nillable(),
		field.String("state").Optional().Nillable(),
		// `order` is a free-form string key in the legacy schema, not an int.
		field.String("order").Optional().Nillable(),
		field.String("meta_title").Optional().Nillable(),
		field.String("meta_description").Optional().Nillable(),
		field.String("description").Optional().Nillable(),
		field.String("used_in_header").Optional().Nillable(),
		field.String("used_in_description").Optional().Nillable(),
		field.String("outcomes_header").Optional().Nillable(),
		field.String("outcomes_description").Optional().Nillable(),
		field.Bool("footer").Optional().Nillable(),
		field.String("footer_name").Optional().Nillable(),
		field.Int("landing_page_to_redirect_id").Optional().Nillable().
			Annotations(AdminInputField{Rename: "LandingPageToRedirectId"}),
	}
}

func (LandingPage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("course", Course.Type).
			Ref("landing_pages").
			Field("course_id").
			Unique().
			Required(),
	}
}

func (LandingPage) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
