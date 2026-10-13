package handlers_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/emailtokens"
	"hexletbasics/internal/testsupport"
)

// A password the hash cannot take (bcrypt stops at 72 bytes) is refused on the
// password field, the old password keeps working and the link is still good for
// another try with a password that fits.
func TestUpdatePasswordRefusesAPasswordTheHashCannotTake(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := createPasswordUser(t, h, "reset-too-long@example.com")
	token := issueEmailToken(t, emailtokens.PasswordReset, u)

	res, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: strings.Repeat("a", 100)}, api.UpdatePasswordParams{Token: token})
	require.NoError(t, err)

	refusal, ok := res.(*api.ValidationError)
	require.True(t, ok, "got %T", res)
	assert.Contains(t, refusal.Errors, "password")
	assert.True(t, signsInWith(t, h, "reset-too-long@example.com", oldPassword))

	res, err = h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "a-password-that-fits"}, api.UpdatePasswordParams{Token: token})
	require.NoError(t, err)
	assert.IsType(t, &api.UserHeaders{}, res)
}
