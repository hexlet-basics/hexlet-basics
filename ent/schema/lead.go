package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Lead maps the legacy `leads` table (sales leads captured from the site). The
// public lead form writes it (createLead) and the admin list reads it.
//
// Two baseline columns stay unmapped on purpose: `ahoy_visit_id`, because ahoy
// is not ported (ADR-0015) and the first-visit attribution travels in the
// LeadCreated event instead, and `state`, which legacy never set on this path
// and nothing reads. Atlas owns the schema, so omitting them drops nothing.
// The table name already matches ent's default plural of `Lead`, so no @Table
// annotation is needed.
//
// `full_name` is not a column: the legacy serializer derives it from the
// associated user.
type Lead struct {
	ent.Schema
}

// Mixin supplies Rails' application-side timestamps: both columns are NOT NULL
// with no DB default, so ent must fill them on insert like ActiveRecord did.
func (Lead) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
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
		// The Metrika client id the form carried.
		field.String("ym_client_id").Optional().Nillable(),
	}
}
