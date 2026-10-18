package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// qnaItemFields are the columns the two question/answer tables share: the
// foreign key of the parent, which the nested URL supplies and the input
// skips, and the question and answer, nullable in the legacy tables but
// required in the contract's QnaItemInput.
func qnaItemFields(parentField, parentColumn string) []ent.Field {
	return []ent.Field{
		field.Int(parentField).StorageKey(parentColumn).
			Annotations(AdminInputField{Skip: true}),
		field.String("question").Optional().Nillable().
			Annotations(AdminInputField{Required: true}),
		field.String("answer").Optional().Nillable().
			Annotations(AdminInputField{Required: true}),
	}
}
