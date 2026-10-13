package handlers_test

import (
	"context"
	"encoding/json"
	wmsql "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/lessonprogress"
	"hexletbasics/ent/user"
	"hexletbasics/internal/api"
	"hexletbasics/internal/events"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

// An email already registered is refused as taken, on the email field, and no
// second account is made.
func TestSignUpRefusesATakenEmail(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	taken := h.DB.User.Query().FirstX(t.Context())
	before := h.DB.User.Query().CountX(t.Context())

	res, err := h.Client.CreateUser(t.Context(), &api.SignUpInput{
		Email:    *taken.Email,
		Password: "correct horse battery staple",
	})
	require.NoError(t, err)

	refusal, ok := res.(*api.ValidationError)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, http.StatusUnprocessableEntity, h.LastStatus())
	assert.Contains(t, refusal.Errors, "email")
	assert.Equal(t, before, h.DB.User.Query().CountX(t.Context()))
}

// A password the hash cannot take (bcrypt stops at 72 bytes) is refused on the
// password field rather than failing the request.
func TestSignUpRefusesAPasswordTheHashCannotTake(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	res, err := h.Client.CreateUser(t.Context(), &api.SignUpInput{
		Email:    "long-password@example.com",
		Password: strings.Repeat("a", 100),
	})
	require.NoError(t, err)

	refusal, ok := res.(*api.ValidationError)
	require.True(t, ok, "got %T", res)
	assert.Contains(t, refusal.Errors, "password")
	assert.False(t, h.DB.User.Query().Where(user.Email("long-password@example.com")).ExistX(t.Context()))
}

// An operation that needs a signed-in user and is reached without one is a
// programming error in the contract's security wiring, not a user mistake: it
// panics loudly instead of continuing as nobody.
func TestMustAuthenticatedUserPanicsWithoutOne(t *testing.T) {
	assert.Panics(t, func() { handlers.MustAuthenticatedUser(context.Background()) })
}

// What the visitor is told when storage fails at sign-in, at each step: the
// session is not issued, the answer is a server error.
func TestSignInAnswersAStorageFailureAsAServerError(t *testing.T) {
	outbox := outboxTable()
	cases := []struct {
		name    string
		storage func(*testing.T, *testsupport.Harness)
	}{
		{"the user cannot be looked up", reads(user.Table)},
		{"the sign-in cannot be recorded", writes(outbox, "INSERT")},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := testsupport.NewAnonymousHarness(t)
			// The credentials are read before the storage is broken.
			input := signInAsActingUser(t, h)
			c.storage(t, h)

			_, err := h.Client.CreateSession(t.Context(), input)

			require.Error(t, err)
			assert.Equal(t, http.StatusInternalServerError, h.LastStatus())
		})
	}
}

// The guest's progress is credited at sign-in; when that cannot be stored, the
// visitor is not signed in half-way — the answer is a server error and the
// cookie that carries their progress is left for the next attempt.
func TestSignInWithGuestProgressAnswersAStorageFailureAsAServerError(t *testing.T) {
	h := testsupport.NewGuestHarness(t, progress.GuestProgress{}.Record(jsCourseSlug, thirdLessonSlug))
	input := signInAsActingUser(t, h)
	breakWrites(t, h, lessonprogress.Table, "INSERT")

	_, err := h.Client.CreateSession(t.Context(), input)

	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, h.LastStatus())
	assert.False(t, clearsGuestCookie(h.ResponseCookies()), "the progress is not consumed")
}

// outboxTable is the table the domain events are written to, unquoted.
func outboxTable() string {
	return strings.Trim(wmsql.DefaultPostgreSQLSchema{}.MessagesTable(events.Topic), `"`)
}

// A session cookie that is not a token this server signed is no session: the
// visitor is anonymous on a public read and refused on a protected one, never
// an error page.
func TestAGarbledSessionCookieIsNoSession(t *testing.T) {
	router := newAuthRouter(t)
	garbled := &http.Cookie{Name: "JWT", Value: "not-a-token"}

	me := do(t, router, http.MethodGet, "/api/me", "", garbled)
	require.Equal(t, http.StatusOK, me.StatusCode)
	var current struct {
		User *struct{} `json:"user"`
	}
	require.NoError(t, json.NewDecoder(me.Body).Decode(&current))
	assert.Nil(t, current.User)

	lead := do(t, router, http.MethodPost, "/api/leads",
		`{"contactMethod":"phone","contactValue":"+10000000000","ymClientId":null}`, garbled)
	assert.Equal(t, http.StatusUnauthorized, lead.StatusCode)
}

// On an operation that needs a session, a valid cookie whose user cannot be
// loaded is a server error, not a sign-out: the visitor is not told to sign in
// again because the database is down. (A public read never fails on it; it
// answers the visitor as anonymous, by design.)
func TestAValidSessionIsAServerErrorWhenTheUserCannotBeLoaded(t *testing.T) {
	answerServerErrors(t, []failure{
		{"the user behind the cookie", reads(user.Table), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.CreateLead(ctx, leadFromTheSite())
			return err
		}},
	})
}
