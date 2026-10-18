package handlers_test

import (
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/jobs"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/events"
	"hexletbasics/internal/testsupport"
)

func signsInWith(t *testing.T, h *testsupport.Harness, email, password string) bool {
	t.Helper()
	_, err := h.Client.CreateSession(t.Context(), &api.SessionInput{Email: email, Password: password})
	if err == nil {
		return true
	}
	requireProblem(t, err, http.StatusUnprocessableEntity)
	return false
}

func TestCreatePasswordReminderEmailsARegisteredUser(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "es")
	u := fixtureUser(t, h, fixture.UserWithPassword)

	err := h.Client.CreatePasswordReminder(t.Context(), &api.EmailInput{Email: "with-password@example.com"})
	require.NoError(t, err)
	assert.Equal(t, jobs.AccountEmailArgs{Purpose: emailtoken.PasswordReset, UserID: u.ID, Locale: "es"}, testsupport.RequireEnqueued(t, h.Services, jobs.AccountEmailArgs{}))
}

// Unlike legacy, an unknown email is not reported.
func TestCreatePasswordReminderAnswersTheSameForAnUnknownEmail(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	err := h.Client.CreatePasswordReminder(t.Context(), &api.EmailInput{Email: "reset-nobody@example.com"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())
	testsupport.RequireNotEnqueued(t, h.Services, jobs.AccountEmailArgs{})
}

func TestCheckPasswordResetToken(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := fixtureUser(t, h, fixture.UserWithPassword)
	token := issueEmailToken(t, emailtoken.PasswordReset, u)

	err := h.Client.CheckPasswordResetToken(t.Context(), api.CheckPasswordResetTokenParams{Token: token})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	err = h.Client.CheckPasswordResetToken(t.Context(),
		api.CheckPasswordResetTokenParams{Token: token[:20] + "forged"})
	requireProblem(t, err, http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestUpdatePasswordSignsInWithTheNewPassword(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := fixtureUser(t, h, fixture.UserWithPassword)

	res, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "new-password"},
		api.UpdatePasswordParams{Token: issueEmailToken(t, emailtoken.PasswordReset, u)})
	require.NoError(t, err)

	session := res
	assert.Equal(t, u.ID, int(session.Response.ID))
	assert.True(t, startsSession(h.ResponseCookies()))
	assert.Len(t, publishedOf[events.UserSignedIn](t, h), 1)
	assert.True(t, signsInWith(t, h, "with-password@example.com", "new-password"))
	assert.False(t, signsInWith(t, h, "with-password@example.com", fixture.Password))
}

// Changing the digest is what kills the link, so it cannot be used twice.
func TestUpdatePasswordRefusesTheSameLinkTwice(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := fixtureUser(t, h, fixture.UserWithPassword)
	token := issueEmailToken(t, emailtoken.PasswordReset, u)

	_, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "new-password"}, api.UpdatePasswordParams{Token: token})
	require.NoError(t, err)

	_, err = h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "another-password"}, api.UpdatePasswordParams{Token: token})
	requireProblem(t, err, http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
	assert.True(t, signsInWith(t, h, "with-password@example.com", "new-password"))
}

func TestUpdatePasswordRefusesAnExpiredLink(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := fixtureUser(t, h, fixture.UserWithPassword)

	_, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "new-password"},
		api.UpdatePasswordParams{Token: issueEmailToken(t, emailtoken.PasswordReset, u, issuedAgo(emailtoken.Lifetime+time.Minute))})
	requireProblem(t, err, http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
	assert.Empty(t, publishedOf[events.UserSignedIn](t, h), "nobody is signed in")
	assert.True(t, signsInWith(t, h, "with-password@example.com", fixture.Password))
}

// The contract's minLength is enforced by the generated server; the rejected
// submission leaves the digest, and so the link, untouched.
func TestUpdatePasswordRejectsAShortPasswordAndKeepsTheLink(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := fixtureUser(t, h, fixture.UserWithPassword)
	token := issueEmailToken(t, emailtoken.PasswordReset, u)

	_, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "12345"}, api.UpdatePasswordParams{Token: token})
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, h.LastStatus())

	err = h.Client.CheckPasswordResetToken(t.Context(), api.CheckPasswordResetTokenParams{Token: token})
	require.NoError(t, err)
}

// Resetting a password signs the account out wherever it was signed in: the
// session a thief may hold stops naming anyone, and only the reset's own new
// session is left.
func TestUpdatePasswordEndsTheSessionsOfTheOldPassword(t *testing.T) {
	h := testsupport.NewHarnessAs(t, fixture.UserWithPassword)
	u := fixtureUser(t, h, fixture.UserWithPassword)

	_, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "new-password"},
		api.UpdatePasswordParams{Token: issueEmailToken(t, emailtoken.PasswordReset, u)})
	require.NoError(t, err)

	_, err = h.Client.GetProfile(t.Context())
	requireProblem(t, err, http.StatusUnauthorized)
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus(), "the session from before the reset is dead")
}
