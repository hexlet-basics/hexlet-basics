package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lib/pq"
)

// StaffMember maps the legacy `staff_members` table — the join granting a user
// an admin role plus the locales that role may act on. The admin surface writes
// it, so the Rails-owned timestamps are supplied by ent. `user_id` is unique
// (one staff record per user) via the baseline index → 409 centrally. The
// default plural already matches the table name.
//
// `allowed_locales` is a native Postgres text array (`varchar[]`), which ent's
// field.Strings (JSON) cannot map; pq.StringArray implements the array wire
// scan/value over database/sql, so it is used as the field's Go type.
type StaffMember struct {
	ent.Schema
}

func (StaffMember) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("role_id"),
		field.Other("allowed_locales", pq.StringArray{}).
			SchemaType(map[string]string{dialect.Postgres: "varchar[]"}),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (StaffMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Field("user_id").Unique().Required(),
		edge.To("role", StaffRole.Type).Field("role_id").Unique().Required(),
	}
}
