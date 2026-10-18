package handlers_test

import (
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/jobs"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/events"
	"hexletbasics/internal/testsupport"
)

// fixtureUser is the fixture user with the label.
func fixtureUser(t *testing.T, h *testsupport.Harness, label fixture.Label) *ent.User {
	t.Helper()

	return h.DB.User.GetX(t.Context(), fixture.ID(label))
}

func issueEmailToken(t *testing.T, purpose emailtoken.Purpose, u *ent.User, opts ...emailtoken.Option) string {
	t.Helper()

	token, err := testsupport.EmailTokens(opts...).Issue(purpose, u)
	require.NoError(t, err)

	return token
}

// issuedAgo backdates a token by the given duration.
func issuedAgo(d time.Duration) emailtoken.Option {
	return emailtoken.WithClock(func() time.Time { return time.Now().Add(-d) })
}

func TestCreateMagicLinkEmailsARegisteredUser(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "es")
	u := fixtureUser(t, h, fixture.UserBlank)

	err := h.Client.CreateMagicLink(t.Context(), &api.EmailInput{Email: "blank-user@example.com"})
	require.NoError(t, err)
	assert.Equal(t, jobs.AccountEmailArgs{Purpose: emailtoken.MagicLink, UserID: u.ID, Locale: "es"}, testsupport.RequireEnqueued(t, h.Services, jobs.AccountEmailArgs{
		Purpose: "",
		UserID:  0,
		Locale:  "",
	}))
}

func TestCreateMagicLinkIgnoresCase(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	u := fixtureUser(t, h, fixture.UserBlank)

	err := h.Client.CreateMagicLink(t.Context(), &api.EmailInput{Email: "Blank-User@Example.com"})
	require.NoError(t, err)

	email := testsupport.RequireEnqueued(t, h.Services, jobs.AccountEmailArgs{
		Purpose: "",
		UserID:  0,
		Locale:  "",
	})
	assert.Equal(t, u.ID, email.UserID)
}

// Only ru and es are emailed; a request in any other language gets ru.
func TestCreateMagicLinkEmailsRussianForOtherLocales(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	testsupport.SpeakTo(h, "en")
	fixtureUser(t, h, fixture.UserBlank)

	err := h.Client.CreateMagicLink(t.Context(), &api.EmailInput{Email: "blank-user@example.com"})
	require.NoError(t, err)

	email := testsupport.RequireEnqueued(t, h.Services, jobs.AccountEmailArgs{
		Purpose: "",
		UserID:  0,
		Locale:  "",
	})
	assert.Equal(t, "ru", email.Locale)
}

func TestCreateMagicLinkAnswersTheSameForAnUnknownEmail(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	before := h.DB.User.Query().CountX(t.Context())

	err := h.Client.CreateMagicLink(t.Context(), &api.EmailInput{Email: "magic-nobody@example.com"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())
	testsupport.RequireNotEnqueued(t, h.Services, jobs.AccountEmailArgs{
		Purpose: "",
		UserID:  0,
		Locale:  "",
	})
	assert.Equal(t, before, h.DB.User.Query().CountX(t.Context()), "no account is created")
}

func TestCreateMagicLinkTwiceWithinTheCooldownEnqueuesOnce(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	fixtureUser(t, h, fixture.UserBlank)

	for range 2 {
		err := h.Client.CreateMagicLink(t.Context(), &api.EmailInput{Email: "blank-user@example.com"})
		require.NoError(t, err)
	}

	// RequireEnqueued fails on a second job: River skipped the repeat.
	testsupport.RequireEnqueued(t, h.Services, jobs.AccountEmailArgs{
		Purpose: "",
		UserID:  0,
		Locale:  "",
	})
}

// Following the link signs in exactly as a password does: the session cookie,
// the guest's progress credited, the sign-in recorded.
func TestConsumeMagicLinkSignsInAndCreditsGuestProgress(t *testing.T) {
	t.Parallel()
	h := testsupport.NewVisitorHarness(t, progressCookie(firstLessonSlug, secondLessonSlug))
	u := fixtureUser(t, h, fixture.UserBlank)

	res, err := h.Client.ConsumeMagicLink(t.Context(),
		api.ConsumeMagicLinkParams{Token: issueEmailToken(t, emailtoken.MagicLink, u)})
	require.NoError(t, err)

	session := res
	assert.Equal(t, u.ID, int(session.Response.ID))
	assert.True(t, startsSession(h.ResponseCookies()))
	assert.ElementsMatch(t, []string{firstLessonSlug, secondLessonSlug},
		finishedLessonSlugs(t, h, u.ID, jsCourseSlug))
	assert.True(t, clearsGuestCookie(h.ResponseCookies()))
	assert.Len(t, publishedOf[events.UserSignedIn](t, h), 1)
}

func TestConsumeMagicLinkWorksTwiceWhileItLives(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	u := fixtureUser(t, h, fixture.UserBlank)
	token := issueEmailToken(t, emailtoken.MagicLink, u)

	for range 2 {
		res, err := h.Client.ConsumeMagicLink(t.Context(), api.ConsumeMagicLinkParams{Token: token})
		require.NoError(t, err)
		assert.IsType(t, &api.UserHeaders{}, res)
	}
}

func TestConsumeMagicLinkRefusesDeadLinks(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)
	u := fixtureUser(t, h, fixture.UserBlank)
	changed := fixtureUser(t, h, fixture.UserWithPassword)
	sentToOldAddress := issueEmailToken(t, emailtoken.MagicLink, changed)
	h.DB.User.UpdateOne(changed).SetEmail("magic-new-address@example.com").ExecX(t.Context())

	cases := map[string]string{
		"expired":             issueEmailToken(t, emailtoken.MagicLink, u, issuedAgo(16*time.Minute)),
		"email changed":       sentToOldAddress,
		"password reset link": issueEmailToken(t, emailtoken.PasswordReset, u),
		"forged":              issueEmailToken(t, emailtoken.MagicLink, u)[:20] + "forged",
		// A link the Rails app emailed before cutover: a MessageVerifier token,
		// which never parses as ours.
		"legacy": "eyJfcmFpbHMiOnsiZGF0YSI6WzFdLCJwdXIiOiJVc2VyXG5tYWdpY19saW5rIn19--0123456789abcdef",
	}
	for name, token := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := h.Client.ConsumeMagicLink(t.Context(), api.ConsumeMagicLinkParams{Token: token})
			requireProblem(t, err, http.StatusNotFound)
			assert.Equal(t, http.StatusNotFound, h.LastStatus())
			assert.Empty(t, publishedOf[events.UserSignedIn](t, h), "nobody is signed in")
		})
	}
}
