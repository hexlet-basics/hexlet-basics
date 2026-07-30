package schema

import (
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
		// One matrix row of RolePermissionsInput; role_id is skipped because
		// the parent role comes from the URL, not the input.
		AdminInput{Type: "RolePermissionInput"},
	}
}

func (StaffRolePermission) Fields() []ent.Field {
	return []ent.Field{
		field.Int("role_id").
			Annotations(AdminInputField{Skip: true}),
		field.String("resource"),
		field.Bool("can_index").Default(false),
		field.Bool("can_create").Default(false),
		field.Bool("can_update").Default(false),
		field.Bool("can_destroy").Default(false),
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

func (StaffRolePermission) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
