package handlers_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/enrollment"
	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// The acting learner is enrolled in JavaScript with the first of three lessons
// finished, so the dashboard reports them as still working on it — with the
// same completion the course page reports, because the same module computes it.
func TestDashboardListsStartedCoursesWithProgress(t *testing.T) {
	h := testsupport.NewHarness(t)
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
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	finishLesson(t, h, secondLessonSlug)
	finishLesson(t, h, thirdLessonSlug)
	markEnrollmentFinished(t, h, jsCourseSlug)

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)
	dashboard := res.(*api.MyDashboard)

	assert.Empty(t, dashboard.StartedEnrollments)
	require.Len(t, dashboard.FinishedEnrollments, 1)
	assert.Equal(t, int32(100), dashboard.FinishedEnrollments[0].Completion)
}

// Both lists at once, each carrying its own course's card.
func TestDashboardReturnsStartedAndFinishedTogether(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	second := seedCourse(t, h, "dashboard-extra", "intro")
	enrollActingUser(t, h, *second.CourseID, "finished")
	giveCatalogueCardForID(t, h, *second.CourseID)

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)
	dashboard := res.(*api.MyDashboard)

	require.Len(t, dashboard.StartedEnrollments, 1)
	require.Len(t, dashboard.FinishedEnrollments, 1)
	assert.Len(t, dashboard.LandingPagesByCourseId, 2)
}

// A course with no listed main catalogue entry has no card to render, so it is
// not on the dashboard at all — the legacy scope. The learner's other course,
// which has one, still is.
func TestDashboardSkipsCoursesWithoutAListedMainCard(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	cardless := seedCourse(t, h, "dashboard-cardless", "intro")
	enrollActingUser(t, h, *cardless.CourseID, "started")

	res, err := h.Client.GetMyDashboard(ctx)
	require.NoError(t, err)
	dashboard := res.(*api.MyDashboard)

	require.Len(t, dashboard.StartedEnrollments, 1, "only the course with a card")
	assert.NotEqual(t, int32(*cardless.CourseID), dashboard.StartedEnrollments[0].CourseId)
	assert.NotContains(t, dashboard.LandingPagesByCourseId, strconv.Itoa(*cardless.CourseID))
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

// giveCatalogueCardForID makes a course dashboard-visible by giving it the
// listed main landing page the card is drawn from. The fixture JavaScript
// course already has one; seeded courses do not.
func giveCatalogueCardForID(t *testing.T, h *testsupport.Harness, courseID int) {
	t.Helper()
	h.DB.LandingPage.Create().
		SetCourseID(courseID).
		SetSlug("card-" + strconv.Itoa(courseID)).
		SetName("Card").
		SetMain(true).
		SetListed(true).
		SetState("published").
		SaveX(t.Context())
}

func enrollActingUser(t *testing.T, h *testsupport.Harness, courseID int, state string) *ent.Enrollment {
	t.Helper()
	return h.DB.Enrollment.Create().
		SetUserID(actingUserID(t, h)).
		SetCourseID(courseID).
		SetState(state).
		SaveX(t.Context())
}

func markEnrollmentFinished(t *testing.T, h *testsupport.Harness, slug string) {
	t.Helper()
	ctx := t.Context()
	crs := h.DB.Course.Query().Where(course.SlugEQ(slug)).OnlyX(ctx)
	h.DB.Enrollment.Update().
		Where(
			enrollment.UserID(actingUserID(t, h)),
			enrollment.CourseID(crs.ID),
		).
		SetState("finished").
		ExecX(ctx)
}
