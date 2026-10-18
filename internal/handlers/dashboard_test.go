package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// The acting learner is enrolled in JavaScript with the first of three lessons
// finished, so the dashboard reports them as still working on it — with the
// same completion the course page reports, because the same module computes it.
func TestDashboardListsStartedCoursesWithProgress(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")

	ctx := context.Background()

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)

	dashboard := res

	assert.Equal(t, http.StatusOK, h.LastStatus())

	require.Len(t, dashboard.Started, 1)
	assert.Empty(t, dashboard.Finished)

	started := dashboard.Started[0]
	assert.Equal(t, int32(33), started.Progress.Completion)
	assert.Equal(t, jsCourseSlug, started.Card.Course.Slug, "each course carries its own card")
	assert.Equal(t, started.Enrollment.CourseId, started.Card.Course.ID)
}

// The resume link is labelled with the lesson's name, as the course page names
// it, never with its slug. This learner has finished nothing, so they resume at
// the first lesson, which has a Russian translation.
func TestDashboardNamesTheNextLessonAsTheCoursePageDoes(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserQuotaUsedToday)
	testsupport.SpeakTo(h, "ru")

	ctx := t.Context()

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)

	next := res.Started[0].Progress.NextLesson.Value

	page, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsLandingSlug})
	require.NoError(t, err)

	lesson, ok := lo.Find(page.Lessons, func(l api.LearnerLesson) bool { return l.Slug == next.Slug })
	require.True(t, ok)

	assert.Equal(t, firstLessonSlug, next.Slug)
	require.False(t, next.Name.Null, "the lesson is translated into the locale")
	assert.Equal(t, lesson.Name.Value, next.Name.Value)
	assert.NotEqual(t, next.Slug, next.Name.Value)
}

// A lesson not translated into the locale has no page there to resume at, so the
// dashboard shows the course without a resume link rather than a slug where a
// name belongs.
func TestDashboardGivesAnUntranslatedNextLessonNoLink(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")

	res, err := h.Client.GetMyDashboard(t.Context())
	require.NoError(t, err)

	course := res.Started[0]

	assert.True(t, course.Progress.NextLesson.Null, "the second lesson has no ru translation")
	assert.InDelta(t, 33, course.Progress.Completion, 0, "the standing is unaffected")
}

// A finished course has nowhere to resume.
func TestDashboardGivesAFinishedCourseNoNextLesson(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserLearnerFinishedJavascript)
	testsupport.SpeakTo(h, "ru")

	res, err := h.Client.GetMyDashboard(t.Context())
	require.NoError(t, err)

	assert.True(t, res.Finished[0].Progress.NextLesson.Null)
}

// A finished course moves to the other list rather than disappearing.
func TestDashboardSeparatesFinishedCourses(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserLearnerFinishedJavascript)
	testsupport.SpeakTo(h, "ru")

	ctx := context.Background()

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)

	dashboard := res

	assert.Empty(t, dashboard.Started)
	require.Len(t, dashboard.Finished, 1)
	assert.Equal(t, int32(100), dashboard.Finished[0].Progress.Completion)
}

// Both lists at once, each carrying its own course's card.
func TestDashboardReturnsStartedAndFinishedTogether(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserLearnerJavascriptAndGo)
	testsupport.SpeakTo(h, "ru")

	ctx := context.Background()

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)

	dashboard := res

	require.Len(t, dashboard.Started, 1)
	require.Len(t, dashboard.Finished, 1)
	assert.NotEqual(t, dashboard.Started[0].Card.Course.ID, dashboard.Finished[0].Card.Course.ID,
		"each carries its own course's card")
}

// A course with no main catalogue entry has no card to render, so it is
// not on the dashboard at all — the legacy scope. The learner's other course,
// which has one, still is.
func TestDashboardSkipsCoursesWithoutAMainCard(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserLearnerHiddenCards)
	testsupport.SpeakTo(h, "ru")

	ctx := context.Background()
	cardless := h.DB.Course.GetX(ctx, fixture.ID(fixture.CourseSingleLesson))

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)

	dashboard := res

	require.Len(t, dashboard.Started, 1, "only the course with a card")
	assert.NotEqual(t, int32(cardless.ID), dashboard.Started[0].Enrollment.CourseId)
}

// The card follows the legacy dashboard scope, `LandingPage.web.where(main:
// true)`: a published page in the request locale of a completed course. A draft
// main page, or one written for another locale, is no card, and the course drops
// off the dashboard with it.
func TestDashboardSkipsCoursesWhoseMainCardIsNotPublicInTheLocale(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserLearnerHiddenCards)
	testsupport.SpeakTo(h, "ru")

	draft := fixture.ID(fixture.CourseWithDraftCard)
	foreign := fixture.ID(fixture.CourseWithEnglishCard)

	res, err := h.Client.GetMyDashboard(t.Context())
	require.NoError(t, err)

	dashboard := res

	require.Len(t, dashboard.Started, 1, "only JavaScript has a public card")
	assert.NotEqual(t, int32(draft), dashboard.Started[0].Enrollment.CourseId)
	assert.NotEqual(t, int32(foreign), dashboard.Started[0].Enrollment.CourseId)
}

// Unlike the catalogue, the dashboard never required `listed`: a learner's own
// course stays on their dashboard even when it is hidden from the public lists.
func TestDashboardKeepsACourseWhoseMainCardIsUnlisted(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarnessAs(t, fixture.UserLearnerUnlistedCard)
	testsupport.SpeakTo(h, "ru")

	hidden := fixture.ID(fixture.CourseWithUnlistedCard)

	res, err := h.Client.GetMyDashboard(t.Context())
	require.NoError(t, err)

	dashboard := res

	assert.True(t, lo.ContainsBy(dashboard.Started, func(c api.DashboardCourse) bool {
		return c.Enrollment.CourseId == int32(hidden)
	}))
}

// A visitor is refused, cookie or no cookie: the dashboard is the incentive to
// have an account. The refusal is a declared response rather than a transport
// error, so the client returns it as one.
func TestDashboardRefusesAGuest(t *testing.T) {
	t.Parallel()
	h := testsupport.NewAnonymousHarness(t)

	_, err := h.Client.GetMyDashboard(context.Background())
	requireProblem(t, err, http.StatusUnauthorized)
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus())
}
