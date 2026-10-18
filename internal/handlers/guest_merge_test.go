package handlers_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/api"
	"hexletbasics/internal/auth"
	"hexletbasics/internal/events"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

// Signing up with guest progress credits exactly the lessons the visitor
// finished.
func TestSignUpCreditsTheLessonsTheGuestFinished(t *testing.T) {
	t.Parallel()
	h := testsupport.NewGuestHarness(t, progressCookie(firstLessonSlug, secondLessonSlug))
	ctx := context.Background()

	res, err := h.Client.CreateUser(ctx, &api.SignUpInput{
		Email:    "newcomer@example.com",
		Password: "correct horse battery staple",
	})
	require.NoError(t, err)

	created := res

	assert.Equal(t, http.StatusCreated, h.LastStatus())

	finished := finishedLessonSlugs(t, h, int(created.Response.ID), jsCourseSlug)
	assert.ElementsMatch(t, []string{firstLessonSlug, secondLessonSlug}, finished,
		"every lesson the cookie names is credited")

	assert.Len(t, publishedOf[events.CourseStarted](t, h), 1, "the enrollment is new")
	assert.Len(t, publishedOf[events.LessonFinished](t, h), 2, "one per credited lesson")
	assert.True(t, clearsGuestCookie(h.ResponseCookies()), "the cookie is consumed")
}

// Signing in with the browser ahead of the account adds what it finished: the
// merge is a union, and a later lesson does not pull the earlier ones with it.
func TestSignInAddsTheGuestsLessonsToTheAccount(t *testing.T) {
	t.Parallel()
	h := testsupport.NewGuestHarness(t, progressCookie(thirdLessonSlug))
	ctx := context.Background()

	res, err := h.Client.CreateSession(ctx, signInAsActingUser(t, h))
	require.NoError(t, err)

	session := res
	finished := finishedLessonSlugs(t, h, int(session.Response.ID), jsCourseSlug)
	assert.ElementsMatch(t, []string{firstLessonSlug, thirdLessonSlug}, finished,
		"the account had only the first; the browser's third is added, the second is not implied")

	assert.Empty(t, publishedOf[events.CourseStarted](t, h), "the learner was already enrolled")
	assert.Len(t, publishedOf[events.LessonFinished](t, h), 1, "only what the merge actually created")
	assert.True(t, clearsGuestCookie(h.ResponseCookies()))
}

// An account that already has everything the cookie names is untouched.
func TestSignInLeavesAnAccountAheadOfTheCookieAlone(t *testing.T) {
	t.Parallel()
	h := testsupport.NewGuestHarness(t, progressCookie(firstLessonSlug))
	ctx := context.Background()
	before := h.DB.LessonProgress.Query().CountX(ctx)

	_, err := h.Client.CreateSession(ctx, signInAsActingUser(t, h))
	require.NoError(t, err)

	assert.Equal(t, before, h.DB.LessonProgress.Query().CountX(ctx), "nothing is written")
	assert.Empty(t, publishedOf[events.LessonFinished](t, h), "and no progress fact is published")
	assert.Empty(t, publishedOf[events.CourseStarted](t, h))
}

// A cookie naming a lesson the current version no longer contains resolves to
// nothing, so that course resets rather than corrupting the account.
func TestSignInIgnoresAGuestLessonMissingFromTheCurrentVersion(t *testing.T) {
	t.Parallel()
	h := testsupport.NewGuestHarness(t, progress.GuestProgress{}.Record(jsCourseSlug, "removed-in-a-later-build"))
	ctx := context.Background()
	before := h.DB.LessonProgress.Query().CountX(ctx)

	_, err := h.Client.CreateSession(ctx, signInAsActingUser(t, h))
	require.NoError(t, err)

	assert.Equal(t, before, h.DB.LessonProgress.Query().CountX(ctx))
}

// Signing in without a guest cookie behaves exactly as it did before.
func TestSignInWithoutAGuestCookieMergesNothing(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	before := h.DB.LessonProgress.Query().CountX(ctx)

	_, err := h.Client.CreateSession(ctx, signInAsActingUser(t, h))
	require.NoError(t, err)

	assert.Equal(t, before, h.DB.LessonProgress.Query().CountX(ctx))
	assert.False(t, clearsGuestCookie(h.ResponseCookies()), "no cookie to clear")
}

// finishedLessonSlugs is the learner's finished lessons in course order.
func finishedLessonSlugs(t *testing.T, h *testsupport.Harness, userID int, courseSlug string) []string {
	t.Helper()
	ctx := t.Context()
	crs := h.DB.Course.Query().Where(course.SlugEQ(courseSlug)).OnlyX(ctx)

	rows := h.DB.LessonProgress.Query().
		Where(
			lessonprogress.UserID(userID),
			lessonprogress.CourseID(crs.ID),
			lessonprogress.StateEQ("finished"),
		).
		AllX(ctx)

	slugs := make([]string, 0, len(rows))
	for _, row := range rows {
		lesson := h.DB.CourseLesson.GetX(ctx, row.LessonID)
		slugs = append(slugs, *lesson.Slug)
	}

	assertEnrollmentExists(t, h, userID, crs.ID)

	return slugs
}

func assertEnrollmentExists(t *testing.T, h *testsupport.Harness, userID, courseID int) {
	t.Helper()
	count := h.DB.Enrollment.Query().
		Where(enrollment.UserID(userID), enrollment.CourseID(courseID)).
		CountX(t.Context())
	require.Equal(t, 1, count, "the merge creates the enrollment when it is absent")
}

// signInAsActingUser returns the credentials of the fixture learner the harness
// acts as, who has a password, for the flows that sign in rather than reuse the
// harness session.
func signInAsActingUser(t *testing.T, h *testsupport.Harness) *api.SessionInput {
	t.Helper()

	return &api.SessionInput{
		Email:    lo.FromPtr(h.DB.User.GetX(t.Context(), actingUserID(t, h)).Email),
		Password: fixture.Password,
	}
}

// clearsGuestCookie inspects the raw Set-Cookie headers rather than the
// client-decoded list: the generated client reads that header as one
// comma-separated array and reports only the first cookie.
// startsSession reports whether the response set a session cookie that names a
// session, not one that ends it.
func startsSession(setCookie []string) bool {
	for _, raw := range setCookie {
		cookie, err := http.ParseSetCookie(raw)
		if err == nil && cookie.Name == auth.CookieName && cookie.Value != "" && cookie.MaxAge >= 0 {
			return true
		}
	}

	return false
}

func clearsGuestCookie(setCookie []string) bool {
	return strings.Contains(strings.Join(setCookie, ","), progress.GuestCookieName+"=")
}
