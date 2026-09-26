package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BlogPostLike maps the legacy `blog_post_likes` join table (a user's like on a
// blog post). The read model needs only the per-post like count, so it is
// queried with a grouped count — no edge is defined (mirrors Lead's
// plain-FK-field approach). The table name matches ent's default plural of
// `BlogPostLike`.
//
// `user_id` is nullable in the baseline (legacy never enforced it), but every
// like the Go stack writes carries the signed-in user.
type BlogPostLike struct {
	ent.Schema
}

func (BlogPostLike) Fields() []ent.Field {
	return []ent.Field{
		field.Int("blog_post_id"),
		field.Int("user_id").Optional().Nillable(),
	}
}

// Indexes declares the (post, user) uniqueness the migration adds. Legacy
// counted a like once per session; the Go stack has no server session, so the
// once-per-post rule lives in the data. ent never creates the index — atlas owns
// the schema — but declaring it keeps the generated OnConflict helpers aware of
// the real constraint.
func (BlogPostLike) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("blog_post_id", "user_id").Unique(),
	}
}

// Mixin supplies the Rails-owned timestamps (NOT NULL, no DB default) now that
// the public like writes the table.
func (BlogPostLike) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
