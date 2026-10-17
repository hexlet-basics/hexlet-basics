// Package staffroles holds the staff-role operations the admin handlers call.
package staffroles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrole"
	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/permissions"
	"hexletbasics/internal/store"
)

// Permissions writes a role's grants.
type Permissions struct {
	db *ent.Client
	tx store.Transactor
}

// NewPermissions builds the production permission writer. The grants are
// replaced inside a transaction, so a save lands whole or not at all.
func NewPermissions(db *ent.Client, tx store.Transactor) *Permissions {
	return &Permissions{db: db, tx: tx}
}

// ErrUnknownPermission is a grant that names no Permission the contract
// declares: a typo would otherwise be stored and silently open nothing. The
// wrapped message lists the names.
var ErrUnknownPermission = errors.New("unknown permission")

// Apply writes the roles screen's save: the role holds exactly grants
// afterwards. A missing role is domain.ErrNotFound, before any write, and a
// grant the contract does not declare is ErrUnknownPermission.
func (p *Permissions) Apply(ctx context.Context, roleID int, grants []string) error {
	if _, err := p.db.StaffRole.Query().Where(staffrole.ID(roleID)).Only(ctx); err != nil {
		return domain.FromStore(err)
	}
	if unknown := lo.Reject(grants, func(name string, _ int) bool { return permissions.Known(name) }); len(unknown) > 0 {
		return fmt.Errorf("%w: %s", ErrUnknownPermission, strings.Join(unknown, ", "))
	}
	return p.tx.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
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
