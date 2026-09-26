package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// BlogPostLike maps the legacy `blog_post_likes` join table (a user's like on a
// blog post). The read model needs only the per-post like count, so it is
// queried with a grouped count — no edge is defined (mirrors Lead's
// plain-FK-field approach). The table name matches ent's default plural of
// `BlogPostLike`.
//
// `user_id` is nullable in the baseline (legacy never enforced it), but every
// like the Go stack writes carries the signed-in user. There is no unique
// (post, user) index on purpose: legacy may hold duplicate pairs and, during
// the ADR-0015 rollback window, still writes them, so one like per user is a
// rule of the like handler, not of the table.
type BlogPostLike struct {
	ent.Schema
}

func (BlogPostLike) Fields() []ent.Field {
	return []ent.Field{
		field.Int("blog_post_id"),
		field.Int("user_id").Optional().Nillable(),
	}
}

// Mixin supplies the Rails-owned timestamps (NOT NULL, no DB default) now that
// the public like writes the table.
func (BlogPostLike) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
