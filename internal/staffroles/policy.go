package staffroles

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/ent/staffrolepermission"
	"hexletbasics/internal/api"
	"hexletbasics/internal/permissions"
)

// ErrForbidden means the user may not perform the operation.
var ErrForbidden = errors.New("staff permission is required")

// Policy decides which admin operations a user may perform. An administrator
// may do everything. A staff member may do what their role is granted, in the
// locales their record lists — legacy `authorize_staff_locale!`.
//
// What an operation needs comes from one of two places while resources move
// over (ADR-0017): the Permission it declares in the contract, which the role
// must hold as a grant, or, for the resources not yet migrated, the table in
// rules.go and the role's create/read/update/delete matrix. Anything neither
// names is for administrators only, so a new admin operation is closed until
// someone classifies it.
type Policy struct {
	db *ent.Client
}

// NewPolicy builds the policy over an ent client.
func NewPolicy(db *ent.Client) *Policy {
	return &Policy{db: db}
}

// Allow reports nil when u may perform operation (an ogen operation name) in
// locale, and ErrForbidden otherwise.
func (p *Policy) Allow(ctx context.Context, u *ent.User, locale, operation string) error {
	if lo.FromPtr(u.Admin) {
		return nil
	}

	staff, err := p.staffMember(ctx, u)
	if err != nil {
		return err
	}
	if staff == nil || !slices.Contains(staff.AllowedLocales, locale) {
		return ErrForbidden
	}

	if required, ok := permissions.For(operation); ok {
		held, err := p.db.StaffRoleGrant.Query().
			Where(
				staffrolegrant.RoleID(staff.RoleID),
				staffrolegrant.Permission(required.Permission),
			).
			Exist(ctx)
		if err != nil {
			return err
		}
		if !held {
			return ErrForbidden
		}
		return nil
	}

	rule, ok := rules[operation]
	if !ok {
		return ErrForbidden
	}
	if rule.anyStaff {
		return nil
	}

	permission, err := p.db.StaffRolePermission.Query().
		Where(
			staffrolepermission.RoleID(staff.RoleID),
			staffrolepermission.Resource(string(rule.resource)),
		).
		Only(ctx)
	if ent.IsNotFound(err) {
		return ErrForbidden
	}
	if err != nil {
		return err
	}
	if !granted(permission, rule.action) {
		return ErrForbidden
	}
	return nil
}

// staffMember is the user's staff record, nil for someone who is not an
// employee. A record already loaded with the user is used as is.
func (p *Policy) staffMember(ctx context.Context, u *ent.User) (*ent.StaffMember, error) {
	if staff := u.Edges.StaffMember; staff != nil {
		return staff, nil
	}
	staff, err := u.QueryStaffMember().Only(ctx)
	if ent.IsNotFound(err) {
		return nil, nil
	}
	return staff, err
}

// granted reads the matrix bit for action.
func granted(p *ent.StaffRolePermission, a action) bool {
	switch a {
	case index:
		return p.CanIndex
	case create:
		return p.CanCreate
	case update:
		return p.CanUpdate
	case destroy:
		return p.CanDestroy
	}
	return false
}

// IsClassified reports whether operation is classified: it declares a
// Permission, or the table names it. Administrators-only operations are listed
// too, so a test can tell "closed on purpose" from "forgotten".
func IsClassified(operation string) bool {
	_, declared := permissions.For(operation)
	_, ok := rules[operation]
	return declared || ok || adminOnly[operation]
}

// OperationName turns an OpenAPI operationId (adminListBanners) into the ogen
// name (AdminListBanners) the security handler receives.
func OperationName(operationID string) api.OperationName {
	return strings.ToUpper(operationID[:1]) + operationID[1:]
}
