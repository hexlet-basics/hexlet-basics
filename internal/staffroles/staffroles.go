// Package staffroles holds the staff-role operations the admin handlers call.
package staffroles

import (
	"context"
	"slices"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrolepermission"
	"hexletbasics/internal/api"
)

// PermissionWriter is the permission-matrix seam the handlers depend on.
type PermissionWriter interface {
	Apply(ctx context.Context, roleID int, permissions []api.RolePermissionInput) error
}

// Permissions writes a role's permission matrix.
type Permissions struct {
	db *ent.Client
}

// NewPermissions builds the production permission writer.
func NewPermissions(db *ent.Client) *Permissions {
	return &Permissions{db: db}
}

// Apply upserts each submitted row by (role_id, resource) — legacy
// StaffRoleMutator#sync_permissions (find_or_initialize_by(resource).update!),
// so unlisted resources are left untouched. The unique index on (role_id,
// resource) decides insert or update in one statement, so the matrix lands
// whole or not at all without a transaction, and two admins saving at once
// cannot both insert a row.
//
// The input type is the contract's because the ent builders' generated SetInput
// takes it.
func (p *Permissions) Apply(ctx context.Context, roleID int, permissions []api.RolePermissionInput) error {
	// A resource listed twice would make the upsert touch one row twice, which
	// PostgreSQL rejects; the last entry wins, as it did when rows were written
	// one after another.
	perms := slices.Clone(permissions)
	slices.Reverse(perms)
	perms = lo.UniqBy(perms, func(r api.RolePermissionInput) api.PermissionResource { return r.Resource })
	slices.Reverse(perms)

	builders := lo.Map(perms, func(r api.RolePermissionInput, _ int) *ent.StaffRolePermissionCreate {
		return p.db.StaffRolePermission.Create().SetRoleID(roleID).SetInput(&r)
	})
	return p.db.StaffRolePermission.CreateBulk(builders...).
		OnConflictColumns(staffrolepermission.FieldRoleID, staffrolepermission.FieldResource).
		UpdateNewValues().
		Exec(ctx)
}
