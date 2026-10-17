package schema

import (
	"encoding/json"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// AuditRecord is one line of the append-only audit trail (ADR-0017): who did
// what, through which channel, to which entity, and how it ended. The table is
// written only by internal/audit, in the same transaction as the change it
// describes, and a database trigger (see the migration) refuses UPDATE and
// DELETE, so nothing here is ever edited.
//
// Deliberately no foreign keys and no created_at/updated_at mixin: a record
// must outlive the user, the Connected App and the entity it names, and it has
// exactly one timestamp, occurred_at, which is also part of the primary key
// (id, occurred_at) so the table can be partitioned by month (ticket #846)
// without a rewrite. ent never creates the table; atlas owns it.
//
// Fields marked Sensitive() in any other schema (password material, tokens) are
// recorded in `changes` by name only; see internal/audit.
type AuditRecord struct {
	ent.Schema
}

func (AuditRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "audit_records"},
	}
}

func (AuditRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Time("occurred_at").Default(time.Now).Immutable(),
		// Null for a refusal of an unauthenticated caller (ticket #845).
		field.Int("actor_user_id").Optional().Nillable().Immutable(),
		// web or mcp. The Connected App and token only exist for mcp.
		field.String("channel").Immutable(),
		field.Int("connected_app_id").Optional().Nillable().Immutable(),
		field.Int("token_id").Optional().Nillable().Immutable(),
		field.String("request_id").Optional().Nillable().Immutable(),
		field.String("address").Optional().Nillable().Immutable(),
		field.String("user_agent").Optional().Nillable().Immutable(),
		// The ogen operation name, e.g. AdminUpdateBlogPost.
		field.String("operation").Immutable(),
		// The ent type and id of the first single-row change, empty when the
		// operation changed no row.
		field.String("entity_type").Optional().Nillable().Immutable(),
		field.String("entity_id").Optional().Nillable().Immutable(),
		field.String("locale").Optional().Nillable().Immutable(),
		// succeeded, forbidden or failed.
		field.String("outcome").Immutable(),
		// The row changes of the operation as JSON (internal/audit.Change).
		field.JSON("changes", json.RawMessage{}).Optional().Immutable(),
	}
}
