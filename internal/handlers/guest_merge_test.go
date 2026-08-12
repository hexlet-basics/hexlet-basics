package handlers_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/api"
	"hexletbasics/internal/events"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

// Signing up with guest progress credits the whole prefix: sequential
// progression means a visitor could not have finished the second lesson without
// the first, so the prefix IS the set — rows are written for a completion the
// rule implies rather than one that was individually observed.
func TestSignUpCreditsTheGuestPrefix(t *testing.T) {
	h := testsupport.NewGuestHarness(t, progress.GuestProgress{}.Record(jsCourseSlug, secondLessonSlug))
	ctx := context.Background()

	res, err := h.Client.CreateUser(ctx, &api.SignUpInput{
		Email:    "newcomer@example.com",
		Password: "correct horse battery staple",
	})
	require.NoError(t, err)
	created, ok := res.(*api.UserHeaders)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, http.StatusCreated, h.LastStatus())

	finished := finishedLessonSlugs(t, h, int(created.Response.ID), jsCourseSlug)
	assert.ElementsMatch(t, []string{firstLessonSlug, secondLessonSlug}, finished,
		"the lesson before the guest's furthest one is credited too")

	assert.Len(t, publishedOf[events.CourseStarted](h), 1, "the enrollment is new")
	assert.Len(t, publishedOf[events.LessonFinished](h), 2, "one per credited lesson")
	assert.True(t, clearsGuestCookie(h.ResponseCookies()), "the cookie is consumed")
}

// Signing in with the browser ahead of the account moves the account forward.
func TestSignInTakesTheFurtherOfTheTwoPositions(t *testing.T) {
	h := testsupport.NewGuestHarness(t, progress.GuestProgress{}.Record(jsCourseSlug, thirdLessonSlug))
	ctx := context.Background()

	res, err := h.Client.CreateSession(ctx, signInAsActingUser(t, h))
	require.NoError(t, err)
	session := res.(*api.UserHeaders)
	finished := finishedLessonSlugs(t, h, int(session.Response.ID), jsCourseSlug)
	assert.ElementsMatch(t, []string{firstLessonSlug, secondLessonSlug, thirdLessonSlug}, finished,
		"the account had only the first; the browser's third pulls the prefix with it")

	assert.Empty(t, publishedOf[events.CourseStarted](h), "the learner was already enrolled")
	assert.Len(t, publishedOf[events.LessonFinished](h), 2, "only what the merge actually created")
	assert.True(t, clearsGuestCookie(h.ResponseCookies()))
}

// An account already further along is untouched: the merge takes the higher
// position, which is the one already stored.
func TestSignInLeavesAnAccountAheadOfTheCookieAlone(t *testing.T) {
	h := testsupport.NewGuestHarness(t, progress.GuestProgress{}.Record(jsCourseSlug, firstLessonSlug))
	ctx := context.Background()
	before := h.DB.LessonProgress.Query().CountX(ctx)

	_, err := h.Client.CreateSession(ctx, signInAsActingUser(t, h))
	require.NoError(t, err)

	assert.Equal(t, before, h.DB.LessonProgress.Query().CountX(ctx), "nothing is written")
	assert.Empty(t, publishedOf[events.LessonFinished](h), "and no progress fact is published")
	assert.Empty(t, publishedOf[events.CourseStarted](h))
}

// A cookie naming a lesson the current version no longer contains resolves to
// nothing, so that course resets rather than corrupting the account.
func TestSignInIgnoresAGuestLessonMissingFromTheCurrentVersion(t *testing.T) {
	h := testsupport.NewGuestHarness(t, progress.GuestProgress{}.Record(jsCourseSlug, "removed-in-a-later-build"))
	ctx := context.Background()
	before := h.DB.LessonProgress.Query().CountX(ctx)

	_, err := h.Client.CreateSession(ctx, signInAsActingUser(t, h))
	require.NoError(t, err)

	assert.Equal(t, before, h.DB.LessonProgress.Query().CountX(ctx))
}

// Signing in without a guest cookie behaves exactly as it did before.
func TestSignInWithoutAGuestCookieMergesNothing(t *testing.T) {
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

// signInAsActingUser gives the fixture learner a password and returns the
// credentials for it: the fixtures carry no digest, because the harness
// normally mints a session directly.
func signInAsActingUser(t *testing.T, h *testsupport.Harness) *api.SessionInput {
	t.Helper()
	return &api.SessionInput{
		Email:    testsupport.GivePassword(t, h, actingUserID(t, h)),
		Password: testsupport.HarnessUserPassword,
	}
}

// clearsGuestCookie inspects the raw Set-Cookie headers rather than the
// client-decoded list: the generated client reads that header as one
// comma-separated array and reports only the first cookie.
func clearsGuestCookie(setCookie []string) bool {
	return strings.Contains(strings.Join(setCookie, ","), progress.GuestCookieName+"=")
}
