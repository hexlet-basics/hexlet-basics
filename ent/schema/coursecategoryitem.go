package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CourseCategoryItem maps the legacy `language_category_items` join table: a
// landing page grouped under a course category. This is what the public
// category page lists (legacy `Language::Category#language_landing_pages`
// goes through it) — not the `language_category_id` column that also sits on
// `language_landing_pages`. The table already exists (baseline migration), so
// this schema needs no migration. Rails validated one row per (category,
// landing page) but the database has no unique index, which is why the page
// query filters landing pages by EXISTS rather than joining the rows.
type CourseCategoryItem struct {
	ent.Schema
}

func (CourseCategoryItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_category_items"},
	}
}

func (CourseCategoryItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("course_category_id").StorageKey("language_category_id"),
		field.Int("landing_page_id").StorageKey("language_landing_page_id"),
	}
}

func (CourseCategoryItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("category", CourseCategory.Type).
			Field("course_category_id").
			Unique().
			Required(),
		edge.From("landing_page", LandingPage.Type).
			Ref("category_items").
			Field("landing_page_id").
			Unique().
			Required(),
	}
}

func (CourseCategoryItem) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
