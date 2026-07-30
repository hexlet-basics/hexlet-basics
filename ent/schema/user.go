package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User maps the legacy `users` table. The admin surface (UserCrud) writes this
// table, so the Rails-owned `created_at`/`updated_at` (NOT NULL, no DB default)
// are supplied by ent on insert/update. The table name already matches ent's
// default plural of `User`, so no @Table annotation is needed.
//
// Only the columns the admin CRUD surface touches are mapped; the many other
// baseline columns (password_digest, oauth uids, webauthn, …) are intentionally
// omitted — atlas owns the schema, so unmapped columns are simply not selected
// and never dropped. Email uniqueness is enforced by the baseline unique index
// (`index_users_on_email`), surfaced as 409 by the central ErrorHandler, so
// there is no code-side pre-check.
type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").Optional().Nillable(),
		// bcrypt hash from the legacy `has_secure_password` column. Mapped so the
		// auth handlers can verify (login) and write (registration) it; atlas owns
		// the schema, so mapping this existing baseline column needs no migration.
		// Sensitive() keeps it out of ent's Stringer/log output.
		field.String("password_digest").Optional().Nillable().Sensitive(),
		field.String("first_name").Optional().Nillable(),
		field.String("last_name").Optional().Nillable(),
		field.Bool("admin").Optional().Nillable(),
		field.Int("assistant_messages_count").Optional().Nillable(),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
