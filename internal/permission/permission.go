// Package permission holds the staff-role operations the admin handlers call.
package permission

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/samber/lo"
	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrole"
	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/permissioncatalog"
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

// errs is the package template for errors (ADR-0018): static messages, the data
// in With, wrapped once at the exported function that knows the operation.
var errs = oops.In("permission")

// ErrUnknownPermission is a grant that names no Permission the contract
// declares: a typo would otherwise be stored and silently open nothing. The
// wrapped message lists the names.
var ErrUnknownPermission = errors.New("unknown permission")

// Check judges the roles screen's submission before it is written: a missing
// role is domain.ErrNotFound and a grant the contract does not declare is
// ErrUnknownPermission. It is the request's rule and not Apply's, because a
// role may still hold a grant the contract has since retired (a migration
// splits one Permission in two), and trusted callers rewriting what a role
// holds must not be refused for it.
func (p *Permissions) Check(ctx context.Context, roleID int, grants []string) error {
	if err := p.requireRole(ctx, roleID); err != nil {
		return err
	}
	if unknown := lo.Reject(grants, func(name string, _ int) bool { return permissioncatalog.Known(name) }); len(unknown) > 0 {
		// The names stay in the message on purpose: the roles screen shows
		// err.Error() as the field's detail, so they are user-facing text and
		// not log data. The HTTP error step (ADR-0018 follow-up) replaces this
		// with a localized detail and the message becomes static.
		names := strings.Join(unknown, ", ")
		return errs.With("unknown", names).Wrapf(ErrUnknownPermission, "%s", names)
	}
	return nil
}

// Apply writes the roles screen's save: the role holds exactly grants
// afterwards. A missing role is domain.ErrNotFound, before any write.
func (p *Permissions) Apply(ctx context.Context, roleID int, grants []string) error {
	if err := p.requireRole(ctx, roleID); err != nil {
		return err
	}
	return p.tx.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		return replaceGrants(ctx, db, roleID, grants)
	})
}

func (p *Permissions) requireRole(ctx context.Context, roleID int) error {
	_, err := p.db.StaffRole.Query().Where(staffrole.ID(roleID)).Only(ctx)
	return domain.FromStore(err)
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
			return errs.With("role_id", roleID).Wrapf(err, "add grants")
		}
	}
	_, err := db.StaffRoleGrant.Delete().
		Where(staffrolegrant.RoleID(roleID), staffrolegrant.PermissionNotIn(grants...)).
		Exec(ctx)
	if err != nil {
		return errs.With("role_id", roleID).Wrapf(err, "remove grants")
	}
	return nil
}
