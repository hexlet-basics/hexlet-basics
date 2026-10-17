package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Session is one signed-in visitor's server-side session (ADR-0017): the table
// the session store reads and writes. It is a new table, so it carries no Rails
// timestamps; the lifetime is the expiry column alone.
//
// The table is `user_sessions`, not `sessions`: the Rails baseline already has a
// `sessions` table (Rails 8's own authentication, keyed by user and never by a
// token this store could use). It is left as it is; a name of our own keeps a
// migration from touching rows that are not ours.
//
// The primary key is the token's digest, not the token. scs hashes the token
// before it reaches the store, so a leaked table cannot be replayed as cookies.
// `data` is the codec's opaque encoding of the session values (the user and the
// password epoch it started with); nothing queries inside it, which is why who a
// session belongs to is not a column.
type Session struct {
	ent.Schema
}

func (Session) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "user_sessions"}}
}

func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("token").Immutable().NotEmpty(),
		field.Bytes("data"),
		field.Time("expiry"),
	}
}

// Indexes serve the periodic removal of expired sessions, which is the only
// query that is not by primary key.
func (Session) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("expiry"),
	}
}
