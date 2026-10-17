package handlers_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

// The acting learner is enrolled in JavaScript with the first of three lessons
// finished, so the dashboard reports them as still working on it — with the
// same completion the course page reports, because the same module computes it.
func TestDashboardListsStartedCoursesWithProgress(t *testing.T) {
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")
	ctx := context.Background()

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)
	dashboard, ok := res.(*api.MyDashboard)
	require.True(t, ok, "got %T", res)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	require.Len(t, dashboard.StartedEnrollments, 1)
	assert.Empty(t, dashboard.FinishedEnrollments)

	started := dashboard.StartedEnrollments[0]
	assert.Equal(t, int32(33), started.Completion)
	assert.Equal(t, secondLessonSlug, started.Progress.NextLessonSlug.Value)

	card, ok := dashboard.LandingPagesByCourseId[strconv.Itoa(int(started.CourseId))]
	require.True(t, ok, "each enrollment is paired with its catalogue card")
	assert.Equal(t, jsCourseSlug, card.Course.Slug)
}

// A finished course moves to the other list rather than disappearing.
func TestDashboardSeparatesFinishedCourses(t *testing.T) {
	h := testsupport.NewHarnessAs(t, fixtures.UserLearnerFinishedJavascript)
	testsupport.SpeakTo(h, "ru")
	ctx := context.Background()

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)
	dashboard := res.(*api.MyDashboard)

	assert.Empty(t, dashboard.StartedEnrollments)
	require.Len(t, dashboard.FinishedEnrollments, 1)
	assert.Equal(t, int32(100), dashboard.FinishedEnrollments[0].Completion)
}

// Both lists at once, each carrying its own course's card.
func TestDashboardReturnsStartedAndFinishedTogether(t *testing.T) {
	h := testsupport.NewHarnessAs(t, fixtures.UserLearnerJavascriptAndGo)
	testsupport.SpeakTo(h, "ru")
	ctx := context.Background()

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)
	dashboard := res.(*api.MyDashboard)

	require.Len(t, dashboard.StartedEnrollments, 1)
	require.Len(t, dashboard.FinishedEnrollments, 1)
	assert.Len(t, dashboard.LandingPagesByCourseId, 2)
}

// A course with no main catalogue entry has no card to render, so it is
// not on the dashboard at all — the legacy scope. The learner's other course,
// which has one, still is.
func TestDashboardSkipsCoursesWithoutAMainCard(t *testing.T) {
	h := testsupport.NewHarnessAs(t, fixtures.UserLearnerHiddenCards)
	testsupport.SpeakTo(h, "ru")
	ctx := context.Background()
	cardless := h.DB.Course.GetX(ctx, fixtures.ID(fixtures.CourseSingleLesson))

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)
	dashboard := res.(*api.MyDashboard)

	require.Len(t, dashboard.StartedEnrollments, 1, "only the course with a card")
	assert.NotEqual(t, int32(cardless.ID), dashboard.StartedEnrollments[0].CourseId)
	assert.NotContains(t, dashboard.LandingPagesByCourseId, strconv.Itoa(cardless.ID))
}

// The card follows the legacy dashboard scope, `LandingPage.web.where(main:
// true)`: a published page in the request locale of a completed course. A draft
// main page, or one written for another locale, is no card, and the course drops
// off the dashboard with it.
func TestDashboardSkipsCoursesWhoseMainCardIsNotPublicInTheLocale(t *testing.T) {
	h := testsupport.NewHarnessAs(t, fixtures.UserLearnerHiddenCards)
	testsupport.SpeakTo(h, "ru")
	draft := fixtures.ID(fixtures.CourseWithDraftCard)
	foreign := fixtures.ID(fixtures.CourseWithEnglishCard)

	res, err := h.Client.GetMyDashboard(t.Context())
	require.NoError(t, err)
	dashboard := res.(*api.MyDashboard)

	require.Len(t, dashboard.StartedEnrollments, 1, "only JavaScript has a public card")
	assert.NotContains(t, dashboard.LandingPagesByCourseId, strconv.Itoa(draft))
	assert.NotContains(t, dashboard.LandingPagesByCourseId, strconv.Itoa(foreign))
}

// Unlike the catalogue, the dashboard never required `listed`: a learner's own
// course stays on their dashboard even when it is hidden from the public lists.
func TestDashboardKeepsACourseWhoseMainCardIsUnlisted(t *testing.T) {
	h := testsupport.NewHarnessAs(t, fixtures.UserLearnerUnlistedCard)
	testsupport.SpeakTo(h, "ru")
	hidden := fixtures.ID(fixtures.CourseWithUnlistedCard)

	res, err := h.Client.GetMyDashboard(t.Context())
	require.NoError(t, err)
	dashboard := res.(*api.MyDashboard)

	assert.Contains(t, dashboard.LandingPagesByCourseId, strconv.Itoa(hidden))
}

// A visitor is refused, cookie or no cookie: the dashboard is the incentive to
// have an account. The refusal is a declared response rather than a transport
// error, so the client returns it as one.
func TestDashboardRefusesAGuest(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	res, err := h.Client.GetMyDashboard(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus())
	assert.IsType(t, &api.ProblemDetails{}, res, "no dashboard is served")
}
