package staffroles

import (
	"context"
	"errors"
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
// may do everything. A staff member may do what their role is granted, in the
// locales their record lists — legacy `authorize_staff_locale!`.
//
// What an operation needs is the Permission it declares in the contract, which
// the role must hold as a grant (ADR-0017). Anything that declares none, apart
// from the any-staff upload, is for administrators only, so a new admin
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

	if slices.Contains(anyStaffOperations, operation) {
		return nil
	}
	return ErrForbidden
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

// IsClassified reports whether operation is classified: it declares a
// Permission, or it is one of the any-staff operations. Administrators-only operations are listed
// too, so a test can tell "closed on purpose" from "forgotten".
func IsClassified(operation string) bool {
	_, declared := permissions.For(operation)
	return declared || slices.Contains(anyStaffOperations, operation) || adminOnly[operation]
}

// OperationName turns an OpenAPI operationId (adminListBanners) into the ogen
// name (AdminListBanners) the security handler receives.
func OperationName(operationID string) api.OperationName {
	return strings.ToUpper(operationID[:1]) + operationID[1:]
}
