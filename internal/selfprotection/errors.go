package selfprotection

import "errors"

// The refusals of the self-protection rules. Each adapter maps them to its own
// wire format: the HTTP handlers answer 403 with a localized Problem Details
// detail.
var (
	// ErrOwnStaffRecord is a change to the actor's own Staff Member record, or
	// one that would move a record onto the actor.
	ErrOwnStaffRecord = errors.New("own staff record")
	// ErrOwnRoleGrants is a non-administrator changing the grants of the Staff
	// Role they hold.
	ErrOwnRoleGrants = errors.New("own role grants")
	// ErrAdministratorFlagRestricted is a non-administrator changing the
	// administrator flag.
	ErrAdministratorFlagRestricted = errors.New("administrator flag restricted")
	// ErrOwnAccountDeletion is the deletion of the actor's own account through
	// the admin operations.
	ErrOwnAccountDeletion = errors.New("own account deletion")
)
