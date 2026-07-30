package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// CategoryQnaItem maps the legacy `language_category_qna_items` table — a
// question/answer pair shown on a course category's page. The admin surface
// writes it (nested under a category), so the Rails-owned timestamps are
// supplied by ent. The default plural doesn't match the legacy table name, so
// it is pinned with an @Table annotation.
type CategoryQnaItem struct {
	ent.Schema
}

func (CategoryQnaItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "language_category_qna_items"},
	}
}

func (CategoryQnaItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("language_category_id"),
		field.String("question").Optional().Nillable(),
		field.String("answer").Optional().Nillable(),
	}
}

func (CategoryQnaItem) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
