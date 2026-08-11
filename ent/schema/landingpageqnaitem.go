package schema

import (
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
		// Same mapping rationale as CategoryQnaItem (shared QnaItemInput).
		AdminInput{Type: "QnaItemInput"},
	}
}

func (LandingPageQnaItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("course_landing_page_id").StorageKey("language_landing_page_id").
			Annotations(AdminInputField{Skip: true}),
		field.String("question").Optional().Nillable().
			Annotations(AdminInputField{Required: true}),
		field.String("answer").Optional().Nillable().
			Annotations(AdminInputField{Required: true}),
	}
}

func (LandingPageQnaItem) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
