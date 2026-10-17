// Package staffroles holds the staff-role operations the admin handlers call.
package staffroles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/ent/staffrolepermission"
	"hexletbasics/internal/api"
	"hexletbasics/internal/store"
)

// Permissions writes a role's permission matrix and grants.
type Permissions struct {
	db *ent.Client
	tx store.Transactor
}

// NewPermissions builds the production permission writer. The transactor makes
// the roles screen's single save land whole: matrix and grants together or not
// at all.
func NewPermissions(db *ent.Client, tx store.Transactor) *Permissions {
	return &Permissions{db: db, tx: tx}
}

// Apply writes the roles screen's save: the matrix rows, then the grants.
func (p *Permissions) Apply(ctx context.Context, roleID int, permissions []api.RolePermissionInput, grants []string) error {
	return p.tx.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		if err := applyMatrix(ctx, db, roleID, permissions); err != nil {
			return err
		}
		return replaceGrants(ctx, db, roleID, grants)
	})
}

// replaceGrants makes the role hold exactly grants. Rows already held are left
// alone (the upsert does nothing on conflict), so a save that changes nothing
// writes nothing, and the rest are removed; the unique index on (role_id,
// permission) decides insert or skip in one statement, so two admins saving at
// once cannot both insert the same grant.
func replaceGrants(ctx context.Context, db *ent.Client, roleID int, grants []string) error {
	grants = lo.Uniq(grants)
	if len(grants) > 0 {
		builders := lo.Map(grants, func(name string, _ int) *ent.StaffRoleGrantCreate {
			return db.StaffRoleGrant.Create().SetRoleID(roleID).SetPermission(name)
		})
		err := db.StaffRoleGrant.CreateBulk(builders...).
			OnConflictColumns(staffrolegrant.FieldRoleID, staffrolegrant.FieldPermission).
			DoNothing().
			Exec(ctx)
		// DoNothing reports no row back for a skipped conflict, which ent
		// surfaces as sql.ErrNoRows; it is the expected outcome, not a failure.
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("add grants: %w", err)
		}
	}
	_, err := db.StaffRoleGrant.Delete().
		Where(staffrolegrant.RoleID(roleID), staffrolegrant.PermissionNotIn(grants...)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("remove grants: %w", err)
	}
	return nil
}

// applyMatrix upserts each submitted row by (role_id, resource) — legacy
// StaffRoleMutator#sync_permissions (find_or_initialize_by(resource).update!),
// so unlisted resources are left untouched. The unique index on (role_id,
// resource) decides insert or update in one statement, so the matrix lands
// whole or not at all without a transaction, and two admins saving at once
// cannot both insert a row.
//
// The input type is the contract's because the ent builders' generated SetInput
// takes it.
func applyMatrix(ctx context.Context, db *ent.Client, roleID int, permissions []api.RolePermissionInput) error {
	// A resource listed twice would make the upsert touch one row twice, which
	// PostgreSQL rejects; the last entry wins, as it did when rows were written
	// one after another.
	perms := slices.Clone(permissions)
	slices.Reverse(perms)
	perms = lo.UniqBy(perms, func(r api.RolePermissionInput) api.PermissionResource { return r.Resource })
	slices.Reverse(perms)

	builders := lo.Map(perms, func(r api.RolePermissionInput, _ int) *ent.StaffRolePermissionCreate {
		return db.StaffRolePermission.Create().SetRoleID(roleID).SetInput(&r)
	})
	return db.StaffRolePermission.CreateBulk(builders...).
		OnConflictColumns(staffrolepermission.FieldRoleID, staffrolepermission.FieldResource).
		UpdateNewValues().
		Exec(ctx)
}
