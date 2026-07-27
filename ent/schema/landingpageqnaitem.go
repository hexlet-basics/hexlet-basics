package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// LandingPageQnaItem maps the legacy `language_landing_page_qna_items` table — a
// question/answer pair shown on a course landing page. Same shape as
// CategoryQnaItem, differing only in the parent FK column. The default plural
// doesn't match the legacy table name, so it is pinned with an @Table annotation.
type LandingPageQnaItem struct {
	ent.Schema
}

func (LandingPageQnaItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_landing_page_qna_items"},
	}
}

func (LandingPageQnaItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("language_landing_page_id"),
		field.String("question").Optional().Nillable(),
		field.String("answer").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
