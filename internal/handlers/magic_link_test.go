package handlers_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/emailtokens"
	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

func createEmailUser(t *testing.T, h *testsupport.Harness, email string) *ent.User {
	t.Helper()
	return h.DB.User.Create().SetEmail(email).SaveX(t.Context())
}

func issueEmailToken(t *testing.T, purpose emailtokens.Purpose, u *ent.User, opts ...emailtokens.Option) string {
	t.Helper()
	token, err := testsupport.EmailTokens(opts...).Issue(purpose, u)
	require.NoError(t, err)
	return token
}

// issuedAgo backdates a token by the given duration.
func issuedAgo(d time.Duration) emailtokens.Option {
	return emailtokens.WithClock(func() time.Time { return time.Now().Add(-d) })
}

func TestCreateMagicLinkEmailsARegisteredUser(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "es")
	u := createEmailUser(t, h, "magic-registered@example.com")

	res, err := h.Client.CreateMagicLink(t.Context(), &api.EmailInput{Email: "magic-registered@example.com"})
	require.NoError(t, err)

	assert.IsType(t, &api.CreateMagicLinkNoContent{}, res)
	assert.Equal(t, []jobs.AccountEmailArgs{
		{Purpose: emailtokens.MagicLink, UserID: u.ID, Locale: "es"},
	}, h.Enqueuer.AccountEmails())
}

func TestCreateMagicLinkIgnoresCase(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := createEmailUser(t, h, "magic-case@example.com")

	_, err := h.Client.CreateMagicLink(t.Context(), &api.EmailInput{Email: "Magic-Case@Example.com"})
	require.NoError(t, err)

	emails := h.Enqueuer.AccountEmails()
	require.Len(t, emails, 1)
	assert.Equal(t, u.ID, emails[0].UserID)
}

// Only ru and es are emailed; a request in any other language gets ru.
func TestCreateMagicLinkEmailsRussianForOtherLocales(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "en")
	createEmailUser(t, h, "magic-en@example.com")

	_, err := h.Client.CreateMagicLink(t.Context(), &api.EmailInput{Email: "magic-en@example.com"})
	require.NoError(t, err)

	emails := h.Enqueuer.AccountEmails()
	require.Len(t, emails, 1)
	assert.Equal(t, "ru", emails[0].Locale)
}

func TestCreateMagicLinkAnswersTheSameForAnUnknownEmail(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	before := h.DB.User.Query().CountX(t.Context())

	res, err := h.Client.CreateMagicLink(t.Context(), &api.EmailInput{Email: "magic-nobody@example.com"})
	require.NoError(t, err)

	assert.IsType(t, &api.CreateMagicLinkNoContent{}, res)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())
	assert.Empty(t, h.Enqueuer.AccountEmails())
	assert.Equal(t, before, h.DB.User.Query().CountX(t.Context()), "no account is created")
}

func TestCreateMagicLinkTwiceWithinTheCooldownEnqueuesOnce(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	createEmailUser(t, h, "magic-twice@example.com")

	for range 2 {
		res, err := h.Client.CreateMagicLink(t.Context(), &api.EmailInput{Email: "magic-twice@example.com"})
		require.NoError(t, err)
		assert.IsType(t, &api.CreateMagicLinkNoContent{}, res)
	}

	assert.Len(t, h.Enqueuer.AccountEmails(), 1)
}

// Following the link signs in exactly as a password does: the session cookie,
// the guest's progress credited, the sign-in recorded.
func TestConsumeMagicLinkSignsInAndCreditsGuestProgress(t *testing.T) {
	h := testsupport.NewVisitorHarness(t, progress.GuestProgress{}.Record(jsCourseSlug, secondLessonSlug))
	u := createEmailUser(t, h, "magic-guest@example.com")

	res, err := h.Client.ConsumeMagicLink(t.Context(),
		api.ConsumeMagicLinkParams{Token: issueEmailToken(t, emailtokens.MagicLink, u)})
	require.NoError(t, err)

	session, ok := res.(*api.UserHeaders)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, u.ID, int(session.Response.ID))
	assert.Contains(t, h.ResponseCookies()[0], "JWT=")
	assert.ElementsMatch(t, []string{firstLessonSlug, secondLessonSlug},
		finishedLessonSlugs(t, h, u.ID, jsCourseSlug))
	assert.True(t, clearsGuestCookie(h.ResponseCookies()))
	assert.Len(t, publishedOf[events.UserSignedIn](h), 1)
}

func TestConsumeMagicLinkWorksTwiceWhileItLives(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := createEmailUser(t, h, "magic-repeat@example.com")
	token := issueEmailToken(t, emailtokens.MagicLink, u)

	for range 2 {
		res, err := h.Client.ConsumeMagicLink(t.Context(), api.ConsumeMagicLinkParams{Token: token})
		require.NoError(t, err)
		assert.IsType(t, &api.UserHeaders{}, res)
	}
}

func TestConsumeMagicLinkRefusesDeadLinks(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	u := createEmailUser(t, h, "magic-dead@example.com")
	changed := createEmailUser(t, h, "magic-old-address@example.com")
	sentToOldAddress := issueEmailToken(t, emailtokens.MagicLink, changed)
	h.DB.User.UpdateOne(changed).SetEmail("magic-new-address@example.com").ExecX(t.Context())

	cases := map[string]string{
		"expired":             issueEmailToken(t, emailtokens.MagicLink, u, issuedAgo(16*time.Minute)),
		"email changed":       sentToOldAddress,
		"password reset link": issueEmailToken(t, emailtokens.PasswordReset, u),
		"forged":              issueEmailToken(t, emailtokens.MagicLink, u)[:20] + "forged",
		// A link the Rails app emailed before cutover: a MessageVerifier token,
		// which never parses as ours.
		"legacy": "eyJfcmFpbHMiOnsiZGF0YSI6WzFdLCJwdXIiOiJVc2VyXG5tYWdpY19saW5rIn19--0123456789abcdef",
	}
	for name, token := range cases {
		t.Run(name, func(t *testing.T) {
			res, err := h.Client.ConsumeMagicLink(t.Context(), api.ConsumeMagicLinkParams{Token: token})
			require.NoError(t, err)

			assert.IsType(t, &api.NotFoundError{}, res)
			assert.Equal(t, http.StatusNotFound, h.LastStatus())
			assert.Empty(t, publishedOf[events.UserSignedIn](h), "nobody is signed in")
		})
	}
}
