package accounts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/testsupport"
)

// The fixtures carry no password digests: bob (1002) gets one here, carol (1003)
// stays without.
const (
	bobID   = 1002
	carolID = 1003
)

func bobWithPassword(t *testing.T, db *ent.Client, password string) *accounts.Credentials {
	t.Helper()
	credentials := accounts.NewCredentials(db)
	_, err := credentials.SetPassword(t.Context(), db.User.GetX(t.Context(), bobID), password)
	require.NoError(t, err)
	return credentials
}

func TestAuthenticateIgnoresEmailCaseAndSpaces(t *testing.T) {
	db := testsupport.NewClient(t)
	credentials := bobWithPassword(t, db, "secret")

	got, err := credentials.Authenticate(t.Context(), "  Bob@Example.COM ", "secret")

	require.NoError(t, err)
	assert.Equal(t, bobID, got.ID)
}

func TestAuthenticateRefusesEveryFailureTheSameWay(t *testing.T) {
	db := testsupport.NewClient(t)
	credentials := bobWithPassword(t, db, "secret")

	_, err := credentials.Authenticate(t.Context(), "bob@example.com", "nope")
	assert.ErrorIs(t, err, accounts.ErrInvalidCredentials, "wrong password")

	_, err = credentials.Authenticate(t.Context(), "nobody@example.com", "secret")
	assert.ErrorIs(t, err, accounts.ErrInvalidCredentials, "unknown email")

	_, err = credentials.Authenticate(t.Context(), "carol@example.com", "secret")
	assert.ErrorIs(t, err, accounts.ErrInvalidCredentials, "no password")

	db.User.UpdateOneID(bobID).SetState(accounts.StateRemoved).ClearEmail().ExecX(t.Context())
	_, err = credentials.Authenticate(t.Context(), "bob@example.com", "secret")
	assert.ErrorIs(t, err, accounts.ErrInvalidCredentials, "removed account")
}

func TestSetPasswordReplacesTheOldOne(t *testing.T) {
	db := testsupport.NewClient(t)
	credentials := bobWithPassword(t, db, "old")

	_, err := credentials.SetPassword(t.Context(), db.User.GetX(t.Context(), bobID), "new")
	require.NoError(t, err)

	_, err = credentials.Authenticate(t.Context(), "bob@example.com", "old")
	assert.ErrorIs(t, err, accounts.ErrInvalidCredentials)
	_, err = credentials.Authenticate(t.Context(), "bob@example.com", "new")
	assert.NoError(t, err)
}

func TestEnsureAnotherAdmin(t *testing.T) {
	db := testsupport.NewClient(t)
	const aliceID, danaID = 1001, 1004

	require.NoError(t, accounts.EnsureAnotherAdmin(t.Context(), db, aliceID), "dana is the other administrator")
	require.NoError(t, accounts.EnsureAnotherAdmin(t.Context(), db, carolID), "carol is not an administrator")

	db.User.UpdateOneID(danaID).SetAdmin(false).ExecX(t.Context())
	assert.ErrorIs(t, accounts.EnsureAnotherAdmin(t.Context(), db, aliceID), accounts.ErrLastAdmin)
	assert.NoError(t, accounts.EnsureAnotherAdmin(t.Context(), db, carolID), "still not an administrator")
}
