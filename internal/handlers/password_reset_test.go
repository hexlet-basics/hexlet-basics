package handlers_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/emailtokens"
	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport"
)

const oldPassword = "old-password"

// createPasswordUser is a User who signs in with a password: the reset link's
// fingerprint covers the digest, so a passwordless User would not show that
// changing the password kills the link.
func createPasswordUser(t *testing.T, h *testsupport.Harness, email string) *ent.User {
	t.Helper()
	digest, err := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.MinCost)
	require.NoError(t, err)
	return h.DB.User.Create().SetEmail(email).SetPasswordDigest(string(digest)).SaveX(t.Context())
}

func signsInWith(t *testing.T, h *testsupport.Harness, email, password string) bool {
	t.Helper()
	res, err := h.Client.CreateSession(t.Context(), &api.SessionInput{Email: email, Password: password})
	require.NoError(t, err)
	_, ok := res.(*api.UserHeaders)
	return ok
}

func TestCreatePasswordReminderEmailsARegisteredUser(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "es")
	u := createPasswordUser(t, h, "reset-registered@example.com")

	res, err := h.Client.CreatePasswordReminder(t.Context(), &api.EmailInput{Email: "reset-registered@example.com"})
	require.NoError(t, err)

	assert.IsType(t, &api.CreatePasswordReminderNoContent{}, res)
	assert.Equal(t, []jobs.AccountEmailArgs{
		{Purpose: emailtokens.PasswordReset, UserID: u.ID, Locale: "es"},
	}, h.Enqueuer.AccountEmails())
}

// Unlike legacy, an unknown email is not reported.
func TestCreatePasswordReminderAnswersTheSameForAnUnknownEmail(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	res, err := h.Client.CreatePasswordReminder(t.Context(), &api.EmailInput{Email: "reset-nobody@example.com"})
	require.NoError(t, err)

	assert.IsType(t, &api.CreatePasswordReminderNoContent{}, res)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())
	assert.Empty(t, h.Enqueuer.AccountEmails())
}

func TestCheckPasswordResetToken(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := createPasswordUser(t, h, "reset-check@example.com")
	token := issueEmailToken(t, emailtokens.PasswordReset, u)

	res, err := h.Client.CheckPasswordResetToken(t.Context(), api.CheckPasswordResetTokenParams{Token: token})
	require.NoError(t, err)
	assert.IsType(t, &api.CheckPasswordResetTokenNoContent{}, res)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	res, err = h.Client.CheckPasswordResetToken(t.Context(),
		api.CheckPasswordResetTokenParams{Token: token[:20] + "forged"})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestUpdatePasswordSignsInWithTheNewPassword(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := createPasswordUser(t, h, "reset-success@example.com")

	res, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "new-password"},
		api.UpdatePasswordParams{Token: issueEmailToken(t, emailtokens.PasswordReset, u)})
	require.NoError(t, err)

	session, ok := res.(*api.UserHeaders)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, u.ID, int(session.Response.ID))
	assert.Contains(t, h.ResponseCookies()[0], "JWT=")
	assert.Len(t, publishedOf[events.UserSignedIn](h), 1)
	assert.True(t, signsInWith(t, h, "reset-success@example.com", "new-password"))
	assert.False(t, signsInWith(t, h, "reset-success@example.com", oldPassword))
}

// Changing the digest is what kills the link, so it cannot be used twice.
func TestUpdatePasswordRefusesTheSameLinkTwice(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := createPasswordUser(t, h, "reset-twice@example.com")
	token := issueEmailToken(t, emailtokens.PasswordReset, u)

	_, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "new-password"}, api.UpdatePasswordParams{Token: token})
	require.NoError(t, err)

	res, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "another-password"}, api.UpdatePasswordParams{Token: token})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
	assert.True(t, signsInWith(t, h, "reset-twice@example.com", "new-password"))
}

func TestUpdatePasswordRefusesAnExpiredLink(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := createPasswordUser(t, h, "reset-expired@example.com")

	res, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "new-password"},
		api.UpdatePasswordParams{Token: issueEmailToken(t, emailtokens.PasswordReset, u, issuedAgo(16*time.Minute))})
	require.NoError(t, err)

	assert.IsType(t, &api.NotFoundError{}, res)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
	assert.Empty(t, publishedOf[events.UserSignedIn](h), "nobody is signed in")
	assert.True(t, signsInWith(t, h, "reset-expired@example.com", oldPassword))
}

// The contract's minLength is enforced by the generated server; the rejected
// submission leaves the digest, and so the link, untouched.
func TestUpdatePasswordRejectsAShortPasswordAndKeepsTheLink(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := createPasswordUser(t, h, "reset-short@example.com")
	token := issueEmailToken(t, emailtokens.PasswordReset, u)

	_, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "12345"}, api.UpdatePasswordParams{Token: token})
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, h.LastStatus())

	res, err := h.Client.CheckPasswordResetToken(t.Context(), api.CheckPasswordResetTokenParams{Token: token})
	require.NoError(t, err)
	assert.IsType(t, &api.CheckPasswordResetTokenNoContent{}, res)
}
