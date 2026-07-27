package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Lead maps the legacy `leads` table (sales leads captured from the site). The
// admin surface is read-only (list), so this schema maps only the columns the
// list serializer reads; unmapped baseline columns (state, ahoy_visit_id,
// ym_client_id, updated_at) are simply not selected — atlas owns the schema, so
// omitting them here never drops anything. The table name already matches ent's
// default plural of `Lead`, so no @Table annotation is needed.
//
// `full_name` is not a column: the legacy serializer derives it from the
// associated user. Until the User schema lands, the converter returns it null
// (contract-valid: `fullName: string | null`).
type Lead struct {
	ent.Schema
}

func (Lead) Fields() []ent.Field {
	return []ent.Field{
		// NOT NULL in the baseline (FK to users), hence a value field.
		field.Int("user_id"),
		field.String("email").Optional().Nillable(),
		field.String("phone").Optional().Nillable(),
		field.String("whatsapp").Optional().Nillable(),
		field.String("telegram").Optional().Nillable(),
		field.String("survey_answers_data").Optional().Nillable(),
		field.String("courses_data").Optional().Nillable(),
		// Rails-owned timestamp; immutable since admin never writes leads.
		field.Time("created_at").Immutable(),
	}
}
