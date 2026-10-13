package staffroles

import (
	"context"
	"errors"
	"slices"
	"strings"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrolepermission"
	"hexletbasics/internal/api"
)

// ErrForbidden means the user may not perform the operation.
var ErrForbidden = errors.New("staff permission is required")

// Policy decides which admin operations a user may perform. An administrator
// may do everything. A staff member may do what their role's permission matrix
// grants on the operation's resource, in the locales their record lists —
// legacy `authorize_staff_resource!` and `authorize_staff_locale!`. Anything the
// table below does not name is for administrators only, so a new admin
// operation is closed until someone classifies it.
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
	if u.Admin != nil && *u.Admin {
		return nil
	}

	staff, err := p.staffMember(ctx, u)
	if err != nil {
		return err
	}
	if staff == nil || !slices.Contains(staff.AllowedLocales, locale) {
		return ErrForbidden
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

// IsClassified reports whether operation is classified: administrators-only
// operations are listed too, so a test can tell "closed on purpose" from
// "forgotten".
func IsClassified(operation string) bool {
	_, ok := rules[operation]
	return ok || adminOnly[operation]
}

// OperationName turns an OpenAPI operationId (adminListBanners) into the ogen
// name (AdminListBanners) the security handler receives.
func OperationName(operationID string) api.OperationName {
	return strings.ToUpper(operationID[:1]) + operationID[1:]
}
