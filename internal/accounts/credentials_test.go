package accounts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/testsupport"
)

func newUser(t *testing.T, db *ent.Client, email, password string) *ent.User {
	t.Helper()
	credentials := accounts.NewCredentials(db)
	u := db.User.Create().SetEmail(email).SaveX(t.Context())
	u, err := credentials.SetPassword(t.Context(), u, password)
	require.NoError(t, err)
	return u
}

func TestAuthenticateIgnoresEmailCaseAndSpaces(t *testing.T) {
	db := testsupport.NewClient(t)
	want := newUser(t, db, "ann@example.test", "secret")

	got, err := accounts.NewCredentials(db).Authenticate(t.Context(), "  Ann@Example.TEST ", "secret")

	require.NoError(t, err)
	assert.Equal(t, want.ID, got.ID)
}

func TestAuthenticateRefusesEveryFailureTheSameWay(t *testing.T) {
	db := testsupport.NewClient(t)
	credentials := accounts.NewCredentials(db)
	newUser(t, db, "ann@example.test", "secret")
	db.User.Create().SetEmail("nopass@example.test").SaveX(t.Context())
	removed := newUser(t, db, "gone@example.test", "secret")
	db.User.UpdateOneID(removed.ID).SetState(accounts.StateRemoved).ClearEmail().SaveX(t.Context())

	for name, c := range map[string][2]string{
		"wrong password":  {"ann@example.test", "nope"},
		"unknown email":   {"nobody@example.test", "secret"},
		"no password":     {"nopass@example.test", "secret"},
		"removed account": {"gone@example.test", "secret"},
	} {
		_, err := credentials.Authenticate(t.Context(), c[0], c[1])
		assert.ErrorIs(t, err, accounts.ErrInvalidCredentials, name)
	}
}

func TestSetPasswordReplacesTheOldOne(t *testing.T) {
	db := testsupport.NewClient(t)
	credentials := accounts.NewCredentials(db)
	u := newUser(t, db, "ann@example.test", "old")

	_, err := credentials.SetPassword(t.Context(), u, "new")
	require.NoError(t, err)

	_, err = credentials.Authenticate(t.Context(), "ann@example.test", "old")
	assert.ErrorIs(t, err, accounts.ErrInvalidCredentials)
	_, err = credentials.Authenticate(t.Context(), "ann@example.test", "new")
	assert.NoError(t, err)
}
