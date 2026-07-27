package schema

import (
	"time"

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
		field.String("first_name").Optional().Nillable(),
		field.String("last_name").Optional().Nillable(),
		field.Bool("admin").Optional().Nillable(),
		field.Int("assistant_messages_count").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
