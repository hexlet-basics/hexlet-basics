package staffroles

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/staffrolegrant"
	"hexletbasics/internal/api"
	"hexletbasics/internal/permissions"
)

// ErrForbidden means the user may not perform the operation.
var ErrForbidden = errors.New("staff permission is required")

// Policy decides which admin operations a user may perform. An administrator
// may do every operation. A staff member may do what their role is granted, in
// the locales their record lists — legacy `authorize_staff_locale!`.
//
// What an operation needs is the Permission it declares in the contract, which
// the role must hold as a grant (ADR-0020). An administrator passes every
// operation Permission without a grant, except the capability Permissions
// (mcp.read, audit.read, ...): those must be granted explicitly to the role of
// whoever uses them, an administrator included, because they open channels
// (MCP, the audit trail) the administrator flag was never meant to cover.
// Anything that declares no Permission, apart from the operations declared open
// to any staff member, is for administrators only, so a new admin operation is
// closed until someone classifies it.
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
	required, declares := permissions.For(operation)
	admin := lo.FromPtr(u.Admin)
	if admin && (!declares || !permissions.IsCapability(required.Permission)) {
		return nil
	}

	staff, err := p.staffMember(ctx, u)
	if err != nil {
		return err
	}
	// An administrator is not held to the locales of a staff record.
	if staff == nil || (!admin && !slices.Contains(staff.AllowedLocales, locale)) {
		return ErrForbidden
	}

	if declares {
		return p.requireGrant(ctx, staff, required.Permission)
	}

	if permissions.AccessOf(operation) == permissions.AnyStaff {
		return nil
	}
	// Declaring nothing: closed.
	return ErrForbidden
}

// HasCapability reports nil when u holds the capability Permission (mcp.read,
// mcp.write, mcp.manage, audit.read) through their staff role, and
// ErrForbidden otherwise. The administrator flag does not count: the grant is
// the only way in. It is the check the MCP layer makes before it serves a call.
// A name that is not a capability is an error of the caller, not a refusal.
func (p *Policy) HasCapability(ctx context.Context, u *ent.User, capability string) error {
	if !permissions.IsCapability(capability) {
		return fmt.Errorf("%q is not a capability Permission", capability)
	}
	staff, err := p.staffMember(ctx, u)
	if err != nil {
		return err
	}
	if staff == nil {
		return ErrForbidden
	}
	return p.requireGrant(ctx, staff, capability)
}

// requireGrant is nil when the staff member's role holds the Permission.
func (p *Policy) requireGrant(ctx context.Context, staff *ent.StaffMember, permission string) error {
	held, err := p.db.StaffRoleGrant.Query().
		Where(
			staffrolegrant.RoleID(staff.RoleID),
			staffrolegrant.Permission(permission),
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

// IsClassified reports whether operation says who may use it, by declaring a
// Permission or an Access in the contract.
func IsClassified(operation string) bool {
	return permissions.IsClassified(operation)
}

// OperationName turns an OpenAPI operationId (adminListBanners) into the ogen
// name (AdminListBanners) the security handler receives.
func OperationName(operationID string) api.OperationName {
	return strings.ToUpper(operationID[:1]) + operationID[1:]
}
