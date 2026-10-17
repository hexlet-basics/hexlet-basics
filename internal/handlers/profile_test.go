package handlers_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/user"
	"hexletbasics/ent/useraccount"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/emailtokens"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/sessions"
	"hexletbasics/internal/testsupport"
)

func TestGetProfile(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.GetProfile(t.Context())
	require.NoError(t, err)

	profile, ok := res.(*api.User)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, api.NewNilString("Admin"), profile.FirstName)
	assert.Equal(t, api.NewNilString("Learner"), profile.LastName)
}

func TestGetProfileRequiresSignIn(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	_, err := h.Client.GetProfile(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus())
}

func TestUpdateProfile(t *testing.T) {
	h := testsupport.NewHarness(t)

	var noLastName api.NilProfileName
	noLastName.SetToNull()
	res, err := h.Client.UpdateProfile(t.Context(), &api.ProfileInput{
		FirstName: api.NewNilProfileName("Alicia"),
		LastName:  noLastName,
	})
	require.NoError(t, err)

	profile, ok := res.(*api.User)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, api.NewNilString("Alicia"), profile.FirstName)

	u := h.DB.User.GetX(t.Context(), h.UserID)
	assert.Equal(t, "Alicia", lo.FromPtr(u.FirstName))
	assert.Nil(t, u.LastName)
}

func TestDeleteAccount(t *testing.T) {
	h := testsupport.NewHarness(t)
	u := h.DB.User.GetX(t.Context(), h.UserID)
	require.NotNil(t, u.Nickname, "the fixture user has a nickname and a phone to erase")
	other := h.DB.UserAccount.GetX(t.Context(), fixtures.ID(fixtures.UserAccountStaffEditorGithub))

	res, err := h.Client.DeleteAccount(t.Context())
	require.NoError(t, err)
	assert.IsType(t, &api.DeleteAccountNoContent{}, res)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	// The session cookie comes back expired, as on sign-out.
	expired := map[string]bool{}
	for _, raw := range h.ResponseCookies() {
		cookie, err := http.ParseSetCookie(raw)
		require.NoError(t, err)
		expired[cookie.Name] = cookie.MaxAge < 0 || cookie.Value == ""
	}
	assert.Equal(t, map[string]bool{sessions.CookieName: true}, expired)

	// The row stays, removed and erased; what is not identifying is kept.
	removed := h.DB.User.GetX(t.Context(), u.ID)
	assert.Equal(t, accounts.StateRemoved, lo.FromPtr(removed.State))
	assert.Nil(t, removed.FirstName)
	assert.Nil(t, removed.LastName)
	assert.Nil(t, removed.Nickname)
	assert.Nil(t, removed.Email)
	assert.Nil(t, removed.Phone)
	assert.Nil(t, removed.PasswordDigest)
	assert.Equal(t, u.Admin, removed.Admin)

	assert.False(t, h.DB.UserAccount.Query().Where(useraccount.UserID(u.ID)).ExistX(t.Context()))
	assert.True(t, h.DB.UserAccount.Query().Where(useraccount.UID(other.UID)).ExistX(t.Context()))
}

// A session is only as good as the account it names, so the cookie a browser
// still holds — on this device until it is dropped, on any other one for its
// whole life — must stop naming anyone once the account is gone.
func TestRemovedAccountIsSignedOutEverywhere(t *testing.T) {
	h := testsupport.NewHarness(t)
	_, err := h.Client.DeleteAccount(t.Context())
	require.NoError(t, err)

	_, err = h.Client.GetProfile(t.Context())
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus())

	current, err := h.Client.GetCurrentUser(t.Context(), api.GetCurrentUserParams{})
	require.NoError(t, err)
	assert.True(t, current.User.Null)
}

// Removal frees the email: the same address registers again as a new account.
func TestRemovedAccountEmailCanRegisterAgain(t *testing.T) {
	h := testsupport.NewHarness(t)
	email := lo.FromPtr(h.DB.User.GetX(t.Context(), h.UserID).Email)
	_, err := h.Client.DeleteAccount(t.Context())
	require.NoError(t, err)

	res, err := h.Client.CreateUser(t.Context(), &api.SignUpInput{
		FirstName: api.NewNilString("Admin"),
		Email:     email,
		Password:  "s3cret-pass",
	})
	require.NoError(t, err)
	assert.IsType(t, &api.UserHeaders{}, res)
	assert.Equal(t, 2, h.DB.User.Query().Where(user.Or(user.Email(email), user.ID(h.UserID))).CountX(t.Context()))
}

func TestRemovedUserCannotSignInWithPassword(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	assert.False(t, signsInWith(t, h, "removed-user@example.com", fixtures.Password))
	assert.Equal(t, http.StatusUnprocessableEntity, h.LastStatus())
}

func TestRemovedUserCannotFollowAMagicLink(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := fixtureUser(t, h, fixtures.UserRemoved)
	token := issueEmailToken(t, emailtokens.MagicLink, u)

	res, err := h.Client.ConsumeMagicLink(t.Context(), api.ConsumeMagicLinkParams{Token: token})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)
}

// A Password Reset issued to a removed user without a password is bound to an
// empty digest and to an email that is still there; only the state can refuse it.
func TestRemovedUserCannotResetAPassword(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := fixtureUser(t, h, fixtures.UserRemovedWithoutPassword)
	token := issueEmailToken(t, emailtokens.PasswordReset, u)

	res, err := h.Client.UpdatePassword(t.Context(),
		&api.ResetPasswordInput{Password: "new-password"},
		api.UpdatePasswordParams{Token: token})
	require.NoError(t, err)
	assert.IsType(t, &api.NotFoundError{}, res)
}

func localeCookieFrom(t *testing.T, h *testsupport.Harness) *http.Cookie {
	t.Helper()
	for _, raw := range h.ResponseCookies() {
		cookie, err := http.ParseSetCookie(raw)
		require.NoError(t, err)
		if cookie.Name == "locale" {
			return cookie
		}
	}
	t.Fatalf("no locale cookie in %v", h.ResponseCookies())
	return nil
}

func TestSwitchLocaleStoresTheChoiceOnTheUser(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.SwitchLocale(t.Context(), api.SwitchLocaleParams{Locale: api.LocaleEs})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	assert.Equal(t, "es", lo.FromPtr(h.DB.User.GetX(t.Context(), h.UserID).Locale))
	cookie := localeCookieFrom(t, h)
	assert.Equal(t, "es", cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.Greater(t, time.Duration(cookie.MaxAge)*time.Second, 300*24*time.Hour)
}

func TestSwitchLocaleRemembersAVisitorsChoiceInACookie(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	_, err := h.Client.SwitchLocale(t.Context(), api.SwitchLocaleParams{Locale: api.LocaleEn})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())
	assert.Equal(t, "en", localeCookieFrom(t, h).Value)
	assert.Nil(t, h.DB.User.GetX(t.Context(), h.UserID).Locale)
}

var secondAdminID = fixtures.ID(fixtures.UserSecondAdmin)

func TestLastAdministratorCannotRemoveTheirOwnAccount(t *testing.T) {
	h := testsupport.NewHarness(t)
	// The fixtures have three administrators: user_admin_learner (the caller), user_second_admin
	// and the staff member user_admin_staff; the caller is left the only one.
	leaveOnlyAlice(t, h)

	_, err := h.Client.DeleteAccount(t.Context())

	require.Error(t, err)
	assert.Equal(t, http.StatusConflict, h.LastStatus())
	assert.True(t, h.DB.User.Query().Where(user.ID(h.UserID), user.StateIsNil()).ExistX(t.Context()))
}

// A null first name clears the column, as a null last name does: the profile
// form sends null for exactly the names the learner emptied.
func TestUpdateProfileClearsAFirstNameSentAsNull(t *testing.T) {
	h := testsupport.NewHarness(t)

	var noFirstName api.NilProfileName
	noFirstName.SetToNull()
	res, err := h.Client.UpdateProfile(t.Context(), &api.ProfileInput{
		FirstName: noFirstName,
		LastName:  api.NewNilProfileName("Learner"),
	})
	require.NoError(t, err)

	profile, ok := res.(*api.User)
	require.True(t, ok, "got %T", res)
	assert.True(t, profile.FirstName.Null)
	u := h.DB.User.GetX(t.Context(), h.UserID)
	assert.Nil(t, u.FirstName)
	assert.Equal(t, "Learner", lo.FromPtr(u.LastName))
}
