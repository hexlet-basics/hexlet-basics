package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// LanguageLesson maps the legacy `language_lessons` table (a single lesson in a
// course). The admin surface reads it only to enrich other lists (a review's or
// member's lesson slug and natural order), so this schema declares just those
// columns — atlas owns the schema, so the many unmapped columns (state,
// original_code, prepared_code, …) are simply never selected. The struct name's
// snake-plural already matches the table, so no @Table annotation is needed.
//
// slug/natural_order are nullable in the baseline even though the app validates
// slug presence, hence Optional().Nillable().
type LanguageLesson struct {
	ent.Schema
}

func (LanguageLesson) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").Optional().Nillable(),
		field.Int("natural_order").Optional().Nillable(),
	}
}
