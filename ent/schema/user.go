package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
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

// Annotations opts the schema into the generated SetInput builders. Explicit
// because the schema is wider than the admin surface: only the fields carrying
// an AdminInputField marker map onto UserInput — password_digest and the
// counters must never be reachable from an admin payload.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{AdminInput{Type: "UserInput", Explicit: true}}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		// Nullable in the legacy column, required in UserInput.
		field.String("email").Optional().Nillable().
			Annotations(AdminInputField{Required: true}),
		// bcrypt hash from the legacy `has_secure_password` column. Mapped so the
		// auth handlers can verify (login) and write (registration) it; atlas owns
		// the schema, so mapping this existing baseline column needs no migration.
		// Sensitive() keeps it out of ent's Stringer/log output.
		field.String("password_digest").Optional().Nillable().Sensitive(),
		field.String("first_name").Optional().Nillable().
			Annotations(AdminInputField{}),
		field.String("last_name").Optional().Nillable().
			Annotations(AdminInputField{}),
		field.Bool("admin").Optional().Nillable().
			Annotations(AdminInputField{}),
		field.Int("assistant_messages_count").Optional().Nillable(),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
