package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BlogPost maps the legacy `blog_posts` table plus the Go-stack `rich_body`
// column. Blog HTML is intentionally trusted and returned as stored: the five
// production posts are migrated by hand, so there is no ActionText compatibility
// or synchronization layer. The cover remains in ActiveStorage and is read via
// its explicit polymorphic attachment query.
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
		field.Text("rich_body").Default(""),
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
