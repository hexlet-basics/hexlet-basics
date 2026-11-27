// Package permission holds staff access: what a Staff Role may do (Policy),
// which Permissions a role holds (Permissions) and who holds a role
// (StaffEditor, the Staff Member records). The writes carry the
// self-protection rules of ADR-0020, so no adapter can write past them.
package permission

import (
	"context"
	"database/sql"
	"errors"

	"github.com/samber/lo"
	"github.com/samber/oops"

	"hexletbasics/ent"
	"hexletbasics/ent/staffmember"
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
// message is static; UnknownNames reads the offending names back.
var ErrUnknownPermission = errors.New("unknown permission")

// The self-protection refusals of staff management (ADR-0020): whoever may
// manage staff must not raise themselves or remove their own footing. Each is
// a domain.ErrForbidden, distinct for the adapter that words it.
var (
	// ErrOwnStaffRecord is a change to the actor's own Staff Member record, or
	// one that would move a record onto the actor.
	ErrOwnStaffRecord = errors.Join(domain.ErrForbidden, errOwnStaffRecord)
	// ErrOwnRoleGrants is a non-administrator changing the grants of the Staff
	// Role they hold: delegating staff_roles.update would otherwise let its
	// holder grant their own role anything.
	ErrOwnRoleGrants = errors.Join(domain.ErrForbidden, errOwnRoleGrants)
)

// The messages of the refusals above.
var (
	errOwnStaffRecord = errors.New("own staff record")
	errOwnRoleGrants  = errors.New("own role grants")
)

// Set writes the roles screen's save on actor's behalf: the role holds exactly
// grants afterwards. Before any write, a missing role is domain.ErrNotFound, a
// non-administrator changing the role they hold is ErrOwnRoleGrants, and a
// grant the contract does not declare is ErrUnknownPermission. The last check
// is the request's and not Apply's, because a role may still hold a grant the
// contract has since retired (a migration splits one Permission in two), and a
// trusted caller rewriting what a role holds must not be refused for it.
func (p *Permissions) Set(ctx context.Context, actor domain.Actor, roleID int, grants []string) error {
	err := p.requireRole(ctx, roleID)
	if err != nil {
		return err
	}

	if !actor.Administrator {
		holds, err := p.db.StaffMember.Query().
			Where(staffmember.UserID(actor.UserID), staffmember.RoleID(roleID)).
			Exist(ctx)
		if err != nil {
			return store.ToDomain(err)
		}

		if holds {
			return ErrOwnRoleGrants
		}
	}

	if unknown := lo.Reject(grants, func(name string, _ int) bool { return permissioncatalog.Known(name) }); len(unknown) > 0 {
		return errs.With("unknown", unknown).Wrapf(ErrUnknownPermission, "check grants")
	}

	return p.Apply(ctx, roleID, grants)
}

// UnknownNames lists the grants an ErrUnknownPermission refused, for an adapter
// that shows them to the person who typed them (the message stays static).
func UnknownNames(err error) []string {
	oopsErr, ok := oops.AsOops(err)
	if !ok {
		return nil
	}

	names, _ := oopsErr.Context()["unknown"].([]string)

	return names
}

// Apply makes the role hold exactly grants, for a trusted caller (a test's
// arrange, a data fix): no actor, and grants the contract no longer declares
// are kept. A missing role is domain.ErrNotFound, before any write.
func (p *Permissions) Apply(ctx context.Context, roleID int, grants []string) error {
	err := p.requireRole(ctx, roleID)
	if err != nil {
		return err
	}

	return p.tx.WithinTx(ctx, func(_ *sql.Tx, db *ent.Client) error {
		return replaceGrants(ctx, db, roleID, grants)
	})
}

func (p *Permissions) requireRole(ctx context.Context, roleID int) error {
	_, err := p.db.StaffRole.Query().Where(staffrole.ID(roleID)).Only(ctx)

	return store.ToDomain(err)
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
