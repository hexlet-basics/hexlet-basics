package handlers_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// A password the hash cannot take (bcrypt stops at 72 bytes) is refused on the
// password field, the old password keeps working and the link is still good for
// another try with a password that fits.
func TestUpdatePasswordRefusesAPasswordTheHashCannotTake(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	u := fixtureUser(t, h, fixture.UserWithPassword)
	token := issueEmailToken(t, emailtoken.PasswordReset, u)

	_, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: strings.Repeat("a", 100)}, api.UpdatePasswordParams{Token: token})
	problem := problemWithStatus(t, err, http.StatusUnprocessableEntity)
	assert.Contains(t, problem.Response.Errors.Value, "password")
	assert.True(t, signsInWith(t, h, "with-password@example.com", fixture.Password))

	_, err = h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "a-password-that-fits"}, api.UpdatePasswordParams{Token: token})
	require.NoError(t, err)
}
