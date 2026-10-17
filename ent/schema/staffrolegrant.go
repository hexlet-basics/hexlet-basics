package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StaffRoleGrant is one Permission a Staff Role holds (ADR-0018): a row exists
// exactly when the role has the Permission, so a grant is checked with one
// indexed lookup and a new Permission is a catalog entry, not a column. The
// permission name is a plain string because the catalog is generated from the
// contract, not from the database; the unique index on (role_id, permission)
// (see the migration) is what the replace-all write relies on.
type StaffRoleGrant struct {
	ent.Schema
}

func (StaffRoleGrant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "staff_member_role_grants"},
	}
}

func (StaffRoleGrant) Fields() []ent.Field {
	return []ent.Field{
		field.Int("role_id"),
		field.String("permission"),
	}
}

func (StaffRoleGrant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", StaffRole.Type).
			Ref("grants").
			Field("role_id").
			Unique().
			Required(),
	}
}

func (StaffRoleGrant) Mixin() []ent.Mixin {
	return []ent.Mixin{TimestampsMixin{}}
}
