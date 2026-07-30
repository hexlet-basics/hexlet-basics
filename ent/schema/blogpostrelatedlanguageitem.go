package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BlogPostRelatedLanguageItem maps the legacy `blog_post_related_language_items`
// join table: a course promoted on a blog post. The admin related-courses
// action replaces a post's whole set, so ent both reads and writes this table
// (Rails-owned timestamps via TimestampsMixin). `order` mirrors the legacy
// nullable position column; the Go write path fills it from the submitted
// courseIds order so the public post page can keep a stable course list.
type BlogPostRelatedLanguageItem struct {
	ent.Schema
}

func (BlogPostRelatedLanguageItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("blog_post_id"),
		field.Int("language_id"),
		field.Int("order").Optional().Nillable(),
	}
}

func (BlogPostRelatedLanguageItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("post", BlogPost.Type).
			Field("blog_post_id").
			Unique().
			Required(),
		edge.To("course", Course.Type).
			Field("language_id").
			Unique().
			Required(),
	}
}

func (BlogPostRelatedLanguageItem) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
