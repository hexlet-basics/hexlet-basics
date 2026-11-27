package signin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/account"
	"hexletbasics/internal/auth"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/events"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/signin"
	"hexletbasics/internal/testsupport"
)

// rig is a gate over the test's transaction with real collaborators: the
// sign-in fact lands in the outbox and the email job in River's table, both
// rolled back with the test.
type rig struct {
	db       *ent.Client
	services *testsupport.Services
	sessions *auth.Authenticator
	gate     *signin.Gate
}

func newRig(t *testing.T) *rig {
	t.Helper()
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	services := testsupport.NewServices(t, tx, transactor)
	sessions := auth.New(db)
	tracker := progress.New(db, transactor, services.Events, testsupport.NewMocks().Runner)

	return &rig{
		db:       db,
		services: services,
		sessions: sessions,
		gate: signin.New(db, account.NewCredentials(db), account.NewRegistrar(transactor, services.Events),
			testsupport.EmailTokens(), sessions, tracker, services.Events, services.Jobs),
	}
}

// within runs a sign-in inside a request the session middleware has loaded, as
// every adapter calls the gate: beginning a Session needs one.
func (r *rig) within(t *testing.T, enter func(ctx context.Context) (signin.Entered, error)) (signin.Entered, error) {
	t.Helper()

	var (
		entered signin.Entered
		err     error
	)

	handler := r.sessions.Middleware(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		entered, err = enter(req.Context())
	}))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil))

	return entered, err
}

func (r *rig) magicLink(t *testing.T, label fixture.Label) string {
	t.Helper()

	token, err := testsupport.EmailTokens().Issue(emailtoken.MagicLink, r.db.User.GetX(t.Context(), fixture.ID(label)))
	require.NoError(t, err)

	return token
}

// Every way in for an existing account records the sign-in and credits the
// visitor's Guest progress; only then may the adapter clear the cookie.
func TestASignInRecordsItAndCreditsTheGuestProgress(t *testing.T) {
	t.Parallel()
	r := newRig(t)
	userID := fixture.ID(fixture.UserBlank)
	guest := progress.GuestProgress{Entries: []progress.GuestEntry{
		{CourseSlug: "javascript", Lessons: []string{"hello-world", "variables"}},
	}}

	link := r.magicLink(t, fixture.UserBlank)

	entered, err := r.within(t, func(ctx context.Context) (signin.Entered, error) {
		return r.gate.WithMagicLink(ctx, link, signin.Visitor{Locale: "en", Guest: guest})
	})

	require.NoError(t, err)
	assert.Equal(t, userID, entered.User.ID)
	assert.True(t, entered.GuestMerged)
	assert.Len(t, testsupport.PublishedEvents[events.UserSignedIn](t, r.services), 1)

	finished := r.db.LessonProgress.Query().
		Where(lessonprogress.UserID(userID), lessonprogress.StateEQ(progress.StateFinished),
			lessonprogress.HasLessonWith(courselesson.SlugIn("hello-world", "variables"))).
		CountX(t.Context())
	assert.Equal(t, 2, finished)
}

// With no Guest progress there is nothing to credit, and the cookie is left as
// it is.
func TestASignInWithoutGuestProgressMergesNothing(t *testing.T) {
	t.Parallel()
	r := newRig(t)

	link := r.magicLink(t, fixture.UserBlank)

	entered, err := r.within(t, func(ctx context.Context) (signin.Entered, error) {
		return r.gate.WithMagicLink(ctx, link, signin.Visitor{Locale: "en", Guest: progress.GuestProgress{Entries: nil}})
	})

	require.NoError(t, err)
	assert.False(t, entered.GuestMerged)
}

// A sign-up is recorded by the registrar as a sign-up, not as a sign-in too.
func TestASignUpIsNotRecordedAsASignIn(t *testing.T) {
	t.Parallel()
	r := newRig(t)

	entered, err := r.within(t, func(ctx context.Context) (signin.Entered, error) {
		return r.gate.SignUp(ctx, "new-learner@example.com", "a long enough password", nil,
			signin.Visitor{Locale: "ru", Guest: progress.GuestProgress{Entries: nil}})
	})

	require.NoError(t, err)
	assert.Equal(t, "new-learner@example.com", *entered.User.Email)
	assert.Len(t, testsupport.PublishedEvents[events.UserSignedUp](t, r.services), 1)
	assert.Empty(t, testsupport.PublishedEvents[events.UserSignedIn](t, r.services))
}

func TestAWrongPasswordIsRefusedAndRecordsNothing(t *testing.T) {
	t.Parallel()
	r := newRig(t)

	_, err := r.within(t, func(ctx context.Context) (signin.Entered, error) {
		return r.gate.WithPassword(ctx, "blank-user@example.com", "not the password", signin.Visitor{Locale: "en", Guest: progress.GuestProgress{Entries: nil}})
	})

	require.ErrorIs(t, err, account.ErrInvalidCredentials)
	assert.Empty(t, testsupport.PublishedEvents[events.UserSignedIn](t, r.services))
}

// Asking for either emailed link answers the same for an unknown email, and
// sends nothing: neither form tells a visitor which emails are registered.
func TestAskingForALinkForAnUnknownEmailSendsNothing(t *testing.T) {
	t.Parallel()
	r := newRig(t)

	require.NoError(t, r.gate.RequestMagicLink(t.Context(), "nobody@example.com", "en"))
	require.NoError(t, r.gate.RequestPasswordReset(t.Context(), "nobody@example.com", "en"))

	testsupport.RequireNotEnqueued(t, r.services, jobs.AccountEmailArgs{Purpose: "", UserID: 0, Locale: ""})
}

func TestAskingForALinkForARegisteredEmailSendsIt(t *testing.T) {
	t.Parallel()
	r := newRig(t)

	require.NoError(t, r.gate.RequestPasswordReset(t.Context(), "  Blank-User@Example.com ", "es"))

	sent := testsupport.RequireEnqueued(t, r.services, jobs.AccountEmailArgs{Purpose: "", UserID: 0, Locale: ""})
	assert.Equal(t, jobs.AccountEmailArgs{Purpose: emailtoken.PasswordReset, UserID: fixture.ID(fixture.UserBlank), Locale: "es"}, sent)
}

// A reset changes the password, which kills the link: it cannot be used twice.
func TestAResetLinkWorksOnce(t *testing.T) {
	t.Parallel()
	r := newRig(t)
	token, err := testsupport.EmailTokens().Issue(emailtoken.PasswordReset, r.db.User.GetX(t.Context(), fixture.ID(fixture.UserBlank)))
	require.NoError(t, err)

	visitor := signin.Visitor{Locale: "en", Guest: progress.GuestProgress{Entries: nil}}

	require.NoError(t, r.gate.CheckPasswordReset(t.Context(), token))

	_, err = r.within(t, func(ctx context.Context) (signin.Entered, error) {
		return r.gate.ResetPassword(ctx, token, "a brand new password", visitor)
	})
	require.NoError(t, err)

	require.ErrorIs(t, r.gate.CheckPasswordReset(t.Context(), token), emailtoken.ErrInvalid)
}
