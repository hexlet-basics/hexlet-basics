package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StaffRolePermission maps the legacy `staff_member_role_permissions` table —
// one resource's CRUD bits for a role. The baseline has a unique index on
// (role_id, resource), which the permission-matrix upsert targets. The default
// plural doesn't match the legacy table name, so it is pinned.
type StaffRolePermission struct {
	ent.Schema
}

func (StaffRolePermission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "staff_member_role_permissions"},
	}
}

func (StaffRolePermission) Fields() []ent.Field {
	return []ent.Field{
		field.Int("role_id"),
		field.String("resource"),
		field.Bool("can_index").Default(false),
		field.Bool("can_create").Default(false),
		field.Bool("can_update").Default(false),
		field.Bool("can_destroy").Default(false),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (StaffRolePermission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", StaffRole.Type).
			Ref("permissions").
			Field("role_id").
			Unique().
			Required(),
	}
}
