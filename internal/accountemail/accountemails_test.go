package accountemail_test

import (
	"log/slog"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/accountemail"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/mailer"
	"hexletbasics/internal/testsupport"
)

const siteURL = "https://code-basics.test"

var magicLinkURL = regexp.MustCompile(`https://code-basics\.test/ru/magic_links/([A-Za-z0-9_.-]+)`)

func TestSenderMailsRussianMagicLinkThatVerifies(t *testing.T) {
	db, _ := testsupport.NewClientWithTx(t)
	u := db.User.GetX(t.Context(), fixture.ID(fixture.UserBlank))
	tokens := emailtoken.New("test-email-secret")
	recorder := &mailer.Recorder{}
	sender := accountemail.NewSender(db, tokens, recorder, testsupport.NewTranslator(t), siteURL)

	require.NoError(t, sender.Send(t.Context(), emailtoken.MagicLink, u.ID, "ru"))

	sent := recorder.Sent()
	require.Len(t, sent, 1)
	assert.Equal(t, "blank-user@example.com", sent[0].To)
	assert.Equal(t, "Ссылка для входа", sent[0].Subject)
	match := magicLinkURL.FindStringSubmatch(sent[0].HTML)
	require.NotNil(t, match, "the email must carry a ru Magic Link: %s", sent[0].HTML)
	verified, err := tokens.Verify(t.Context(), db, emailtoken.MagicLink, match[1])
	require.NoError(t, err)
	assert.Equal(t, u.ID, verified.ID)
}

func TestSenderMailsSpanishPasswordResetLink(t *testing.T) {
	db, _ := testsupport.NewClientWithTx(t)
	u := db.User.GetX(t.Context(), fixture.ID(fixture.UserBlank))
	recorder := &mailer.Recorder{}
	sender := accountemail.NewSender(db, emailtoken.New("test-email-secret"), recorder,
		testsupport.NewTranslator(t), siteURL)

	require.NoError(t, sender.Send(t.Context(), emailtoken.PasswordReset, u.ID, "es"))

	sent := recorder.Sent()
	require.Len(t, sent, 1)
	assert.Equal(t, "Restablecer contraseña", sent[0].Subject)
	assert.Regexp(t, `https://code-basics\.test/es/password/[A-Za-z0-9_.-]+/edit`, sent[0].HTML)
}

func TestSenderSkipsDeletedUser(t *testing.T) {
	db, _ := testsupport.NewClientWithTx(t)
	recorder := &mailer.Recorder{}
	sender := accountemail.NewSender(db, emailtoken.New("test-email-secret"), recorder,
		testsupport.NewTranslator(t), siteURL)

	require.NoError(t, sender.Send(t.Context(), emailtoken.MagicLink, -1, "ru"))

	assert.Empty(t, recorder.Sent())
}

func TestEnqueueIsUniquePerUserAndPurposeForAMinute(t *testing.T) {
	db, tx := testsupport.NewClientWithTx(t)
	u := db.User.GetX(t.Context(), fixture.ID(fixture.UserBlank))
	client, err := jobs.NewInsertOnlyClient(nil, slog.New(slog.DiscardHandler), nil, nil)
	require.NoError(t, err)

	first, err := client.InsertTx(t.Context(), tx,
		jobs.AccountEmailArgs{Purpose: emailtoken.MagicLink, UserID: u.ID, Locale: "ru"}, nil)
	require.NoError(t, err)
	// Another language within the cooldown is still the same request.
	again, err := client.InsertTx(t.Context(), tx,
		jobs.AccountEmailArgs{Purpose: emailtoken.MagicLink, UserID: u.ID, Locale: "es"}, nil)
	require.NoError(t, err)
	reset, err := client.InsertTx(t.Context(), tx,
		jobs.AccountEmailArgs{Purpose: emailtoken.PasswordReset, UserID: u.ID, Locale: "ru"}, nil)
	require.NoError(t, err)

	assert.False(t, first.UniqueSkippedAsDuplicate)
	assert.True(t, again.UniqueSkippedAsDuplicate)
	assert.False(t, reset.UniqueSkippedAsDuplicate, "the other email kind has its own cooldown")
}

func TestSenderRefusesAnUnknownPurpose(t *testing.T) {
	db, _ := testsupport.NewClientWithTx(t)
	u := db.User.GetX(t.Context(), fixture.ID(fixture.UserBlank))
	recorder := &mailer.Recorder{}
	sender := accountemail.NewSender(db, emailtoken.New("test-email-secret"), recorder,
		testsupport.NewTranslator(t), siteURL)

	err := sender.Send(t.Context(), emailtoken.Purpose("newsletter"), u.ID, "en")

	require.ErrorContains(t, err, "unknown account email purpose")
	assert.Empty(t, recorder.Sent())
}

func TestSenderSkipsAUserWithoutAnEmail(t *testing.T) {
	db, _ := testsupport.NewClientWithTx(t)
	u := db.User.GetX(t.Context(), fixture.ID(fixture.UserWithoutEmail))
	recorder := &mailer.Recorder{}
	sender := accountemail.NewSender(db, emailtoken.New("test-email-secret"), recorder,
		testsupport.NewTranslator(t), siteURL)

	require.NoError(t, sender.Send(t.Context(), emailtoken.MagicLink, u.ID, "en"))

	assert.Empty(t, recorder.Sent(), "there is nobody to mail, and no retry is wanted")
}
