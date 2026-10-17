package accounts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/user"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

const (
	withPasswordEmail = "with-password@example.com"
	removedEmail      = "removed-user@example.com"
	blankEmail        = "blank-user@example.com"
)

func TestAuthenticateIgnoresEmailCaseAndSpaces(t *testing.T) {
	db := testsupport.NewClient(t)

	got, err := accounts.NewCredentials(db).Authenticate(t.Context(), "  With-Password@Example.COM ", fixtures.Password)

	require.NoError(t, err)
	assert.Equal(t, fixtures.ID(fixtures.UserWithPassword), got.ID)
}

func TestAuthenticateRefusesEveryFailureTheSameWay(t *testing.T) {
	db := testsupport.NewClient(t)
	credentials := accounts.NewCredentials(db)

	_, err := credentials.Authenticate(t.Context(), withPasswordEmail, "nope")
	assert.ErrorIs(t, err, accounts.ErrInvalidCredentials, "wrong password")

	_, err = credentials.Authenticate(t.Context(), "nobody@example.com", fixtures.Password)
	assert.ErrorIs(t, err, accounts.ErrInvalidCredentials, "unknown email")

	_, err = credentials.Authenticate(t.Context(), blankEmail, fixtures.Password)
	assert.ErrorIs(t, err, accounts.ErrInvalidCredentials, "no password")

	_, err = credentials.Authenticate(t.Context(), removedEmail, fixtures.Password)
	assert.ErrorIs(t, err, accounts.ErrInvalidCredentials, "removed account")
}

func TestSetPasswordReplacesTheOldOne(t *testing.T) {
	db := testsupport.NewClient(t)
	credentials := accounts.NewCredentials(db)

	_, err := credentials.SetPassword(t.Context(), db.User.GetX(t.Context(), fixtures.ID(fixtures.UserWithPassword)), "new")
	require.NoError(t, err)

	_, err = credentials.Authenticate(t.Context(), withPasswordEmail, fixtures.Password)
	assert.ErrorIs(t, err, accounts.ErrInvalidCredentials)
	_, err = credentials.Authenticate(t.Context(), withPasswordEmail, "new")
	assert.NoError(t, err)
}

func TestEnsureAnotherAdmin(t *testing.T) {
	db := testsupport.NewClient(t)
	admin, secondAdmin, blank := fixtures.ID(fixtures.UserAdminLearner), fixtures.ID(fixtures.UserSecondAdmin), fixtures.ID(fixtures.UserBlank)

	require.NoError(t, accounts.EnsureAnotherAdmin(t.Context(), db, admin), "user_second_admin is the other administrator")
	require.NoError(t, accounts.EnsureAnotherAdmin(t.Context(), db, blank), "user_blank is not an administrator")

	// The fixtures hold a third administrator, the staff member user_admin_staff.
	db.User.Update().Where(user.IDIn(secondAdmin, fixtures.ID("user_admin_staff"))).SetAdmin(false).ExecX(t.Context())
	assert.ErrorIs(t, accounts.EnsureAnotherAdmin(t.Context(), db, admin), accounts.ErrLastAdmin)
	assert.NoError(t, accounts.EnsureAnotherAdmin(t.Context(), db, blank), "still not an administrator")
}
