package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// UserAccount maps the legacy `user_accounts` table: the social sign-in
// identities (provider + uid) linked to a user. Social sign-in is not ported
// (ADR-0015), so the only use is removing an account, which deletes these rows
// as legacy `accounts.clear` did. atlas owns the table, so mapping it needs no
// migration.
type UserAccount struct {
	ent.Schema
}

func (UserAccount) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "user_accounts"}}
}

func (UserAccount) Fields() []ent.Field {
	return []ent.Field{
		// NOT NULL in the baseline (FK to users), hence a value field.
		field.Int("user_id"),
		field.String("provider"),
		field.String("uid"),
	}
}

func (UserAccount) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
