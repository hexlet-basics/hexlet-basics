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
		// Shares QnaItemInput with LandingPageQnaItem. The parent FK comes
		// from the URL (nested resource), so it is skipped; question/answer
		// are nullable legacy columns but required in the contract.
		AdminInput{Type: "QnaItemInput"},
	}
}

func (CategoryQnaItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("language_category_id").
			Annotations(AdminInputField{Skip: true}),
		field.String("question").Optional().Nillable().
			Annotations(AdminInputField{Required: true}),
		field.String("answer").Optional().Nillable().
			Annotations(AdminInputField{Required: true}),
	}
}

func (CategoryQnaItem) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
