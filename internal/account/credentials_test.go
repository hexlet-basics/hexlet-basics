package account_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/user"
	"hexletbasics/internal/account"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

const (
	withPasswordEmail = "with-password@example.com"
	removedEmail      = "removed-user@example.com"
	blankEmail        = "blank-user@example.com"
)

func TestAuthenticateIgnoresEmailCaseAndSpaces(t *testing.T) {
	db := testsupport.NewClient(t)

	got, err := account.NewCredentials(db).Authenticate(t.Context(), "  With-Password@Example.COM ", fixture.Password)

	require.NoError(t, err)
	assert.Equal(t, fixture.ID(fixture.UserWithPassword), got.ID)
}

func TestAuthenticateRefusesEveryFailureTheSameWay(t *testing.T) {
	db := testsupport.NewClient(t)
	credentials := account.NewCredentials(db)

	_, err := credentials.Authenticate(t.Context(), withPasswordEmail, "nope")
	assert.ErrorIs(t, err, account.ErrInvalidCredentials, "wrong password")

	_, err = credentials.Authenticate(t.Context(), "nobody@example.com", fixture.Password)
	assert.ErrorIs(t, err, account.ErrInvalidCredentials, "unknown email")

	_, err = credentials.Authenticate(t.Context(), blankEmail, fixture.Password)
	assert.ErrorIs(t, err, account.ErrInvalidCredentials, "no password")

	_, err = credentials.Authenticate(t.Context(), removedEmail, fixture.Password)
	assert.ErrorIs(t, err, account.ErrInvalidCredentials, "removed account")
}

func TestSetPasswordReplacesTheOldOne(t *testing.T) {
	db := testsupport.NewClient(t)
	credentials := account.NewCredentials(db)

	_, err := credentials.SetPassword(t.Context(), db.User.GetX(t.Context(), fixture.ID(fixture.UserWithPassword)), "new")
	require.NoError(t, err)

	_, err = credentials.Authenticate(t.Context(), withPasswordEmail, fixture.Password)
	assert.ErrorIs(t, err, account.ErrInvalidCredentials)
	_, err = credentials.Authenticate(t.Context(), withPasswordEmail, "new")
	assert.NoError(t, err)
}

func TestEnsureAnotherAdmin(t *testing.T) {
	db := testsupport.NewClient(t)
	admin, secondAdmin, blank := fixture.ID(fixture.UserAdminLearner), fixture.ID(fixture.UserSecondAdmin), fixture.ID(fixture.UserBlank)

	require.NoError(t, account.EnsureAnotherAdmin(t.Context(), db, admin), "user_second_admin is the other administrator")
	require.NoError(t, account.EnsureAnotherAdmin(t.Context(), db, blank), "user_blank is not an administrator")

	// The fixtures hold a third administrator, the staff member user_admin_staff.
	db.User.Update().Where(user.IDIn(secondAdmin, fixture.ID(fixture.UserAdminStaff))).SetAdmin(false).ExecX(t.Context())
	assert.ErrorIs(t, account.EnsureAnotherAdmin(t.Context(), db, admin), account.ErrLastAdmin)
	assert.NoError(t, account.EnsureAnotherAdmin(t.Context(), db, blank), "still not an administrator")
}
