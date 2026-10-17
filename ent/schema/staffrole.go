package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StaffRole maps the legacy `staff_member_roles` table — a named admin role with
// a set of per-resource permissions. The admin surface writes it, so the
// Rails-owned timestamps are supplied by ent. Name uniqueness is enforced by the
// `staff_member_roles_name_uniq` migration (409 via the central handler). The
// default plural doesn't match the legacy table name, so it is pinned.
type StaffRole struct {
	ent.Schema
}

func (StaffRole) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "staff_member_roles"},
		// The contract input is RoleInput, not StaffRoleInput; the default
		// mapping (name required, description nullable) fits as-is.
		AdminInput{Type: "RoleInput"},
		// Legacy `order(:name)`; the list rows carry the permission count, and
		// the single view is the StaffRoleDetail permission matrix.
		AdminCRUD{
			Detail: "StaffRoleDetail",
			Order:  []AdminOrder{{Field: "Name"}},
			Eager:  []string{"WithPermissions()"},
		},
	}
}

func (StaffRole) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description").Optional().Nillable(),
	}
}

func (StaffRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("permissions", StaffRolePermission.Type),
		// The declared Permissions the role holds (ADR-0017); they replace the
		// matrix resource by resource.
		edge.To("grants", StaffRoleGrant.Type),
	}
}

func (StaffRole) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
