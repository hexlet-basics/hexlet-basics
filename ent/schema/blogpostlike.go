package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// BlogPostLike maps the legacy `blog_post_likes` join table (a user's like on a
// blog post). The read model needs only the per-post like count, so this schema
// maps just `blog_post_id` and is queried with a grouped count — no edge is
// defined (mirrors Lead's plain-FK-field approach). The table name matches ent's
// default plural of `BlogPostLike`.
type BlogPostLike struct {
	ent.Schema
}

func (BlogPostLike) Fields() []ent.Field {
	return []ent.Field{
		field.Int("blog_post_id"),
	}
}
