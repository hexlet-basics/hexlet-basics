package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BlogPost maps the legacy `blog_posts` table. The admin surface is READ-ONLY
// for now (list/get) — the new stack will edit the rich body with the Mantine
// editor and store plain HTML in a future phase, so no write path and no new
// column land here yet. ent only reads, so nullable columns are mirrored with
// Optional().Nillable() to avoid scan errors.
//
// The rich body and cover are NOT columns on this table: the body lives in the
// polymorphic `action_text_rich_texts` (ActionText) and the cover in
// ActiveStorage (`active_storage_attachments`/`active_storage_blobs`). Those are
// read via their own schemas + explicit record_type/record_id queries in the
// handler, since ent has no clean polymorphic edge. `related_language_items_count`
// is a real counter column (no join needed); likes are counted from
// `blog_post_likes`.
//
// The table name already matches ent's default plural of `BlogPost`, so no
// @Table annotation is needed.
type BlogPost struct {
	ent.Schema
}

func (BlogPost) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Nillable(),
		field.String("slug").Optional().Nillable(),
		field.String("description").Optional().Nillable(),
		field.String("locale").Optional().Nillable(),
		field.String("state").Optional().Nillable(),
		// creator_id is NOT NULL (FK to users); the read model embeds the creator.
		field.Int("creator_id"),
		field.Int("language_id").Optional().Nillable(),
		// App-maintained counter; NOT NULL DEFAULT 0 in the baseline.
		field.Int("related_language_items_count").Default(0),
		field.Time("created_at").Immutable(),
	}
}

func (BlogPost) Edges() []ent.Edge {
	return []ent.Edge{
		// belongs-to: the FK (`creator_id`) lives on this table, so the edge owns
		// the field. Loaded via WithCreator for the embedded creator User.
		edge.To("creator", User.Type).
			Field("creator_id").
			Unique().
			Required(),
	}
}
