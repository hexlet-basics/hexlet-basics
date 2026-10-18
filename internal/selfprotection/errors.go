package selfprotection

import "hexletbasics/internal/domain"

// refused is a refusal of the self-protection rules. It is a
// domain.ErrForbidden, so every adapter answers it as a forbidden operation,
// while each value stays distinct for the adapter that words it.
type refused string

func (r refused) Error() string { return string(r) }

// Is makes every refusal match domain.ErrForbidden.
func (refused) Is(target error) bool { return target == domain.ErrForbidden }

// The refusals of the self-protection rules. The HTTP adapter words each with a
// localized Problem Details detail.
var (
	// ErrOwnStaffRecord is a change to the actor's own Staff Member record, or
	// one that would move a record onto the actor.
	ErrOwnStaffRecord error = refused("own staff record")
	// ErrOwnRoleGrants is a non-administrator changing the grants of the Staff
	// Role they hold.
	ErrOwnRoleGrants error = refused("own role grants")
	// ErrAdministratorFlagRestricted is a non-administrator changing the
	// administrator flag.
	ErrAdministratorFlagRestricted error = refused("administrator flag restricted")
	// ErrOwnAccountDeletion is the deletion of the actor's own account through
	// the admin operations.
	ErrOwnAccountDeletion error = refused("own account deletion")
)
