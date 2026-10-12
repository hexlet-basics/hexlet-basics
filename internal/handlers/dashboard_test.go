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
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")
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
	testsupport.SpeakTo(h, "ru")
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

// A course with no main catalogue entry has no card to render, so it is
// not on the dashboard at all — the legacy scope. The learner's other course,
// which has one, still is.
func TestDashboardSkipsCoursesWithoutAMainCard(t *testing.T) {
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")
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

// The card follows the legacy dashboard scope, `LandingPage.web.where(main:
// true)`: a published page in the request locale of a completed course. A draft
// main page, or one written for another locale, is no card, and the course drops
// off the dashboard with it.
func TestDashboardSkipsCoursesWhoseMainCardIsNotPublicInTheLocale(t *testing.T) {
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")

	draft := seedCourse(t, h, "dashboard-draft", "intro")
	enrollActingUser(t, h, *draft.CourseID, "started")
	giveCatalogueCardForID(t, h, *draft.CourseID, func(c *ent.CourseLandingPageCreate) { c.SetState("draft") })

	foreign := seedCourse(t, h, "dashboard-foreign", "intro")
	enrollActingUser(t, h, *foreign.CourseID, "started")
	giveCatalogueCardForID(t, h, *foreign.CourseID, func(c *ent.CourseLandingPageCreate) { c.SetLocale("en") })

	res, err := h.Client.GetMyDashboard(t.Context())
	require.NoError(t, err)
	dashboard := res.(*api.MyDashboard)

	require.Len(t, dashboard.StartedEnrollments, 1, "only the fixture course has a public card")
	assert.NotContains(t, dashboard.LandingPagesByCourseId, strconv.Itoa(*draft.CourseID))
	assert.NotContains(t, dashboard.LandingPagesByCourseId, strconv.Itoa(*foreign.CourseID))
}

// Unlike the catalogue, the dashboard never required `listed`: a learner's own
// course stays on their dashboard even when it is hidden from the public lists.
func TestDashboardKeepsACourseWhoseMainCardIsUnlisted(t *testing.T) {
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")

	hidden := seedCourse(t, h, "dashboard-unlisted", "intro")
	enrollActingUser(t, h, *hidden.CourseID, "started")
	giveCatalogueCardForID(t, h, *hidden.CourseID, func(c *ent.CourseLandingPageCreate) { c.SetListed(false) })

	res, err := h.Client.GetMyDashboard(t.Context())
	require.NoError(t, err)
	dashboard := res.(*api.MyDashboard)

	assert.Contains(t, dashboard.LandingPagesByCourseId, strconv.Itoa(*hidden.CourseID))
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
// main landing page the card is drawn from, in ru, the locale the fixture JavaScript card is written in,
// and finishing the course the way every public one is. The fixture JavaScript
// course already has one; seeded courses do not. adjust edits the page before
// it is saved, so a test can make it a draft or move it to another locale.
func giveCatalogueCardForID(t *testing.T, h *testsupport.Harness, courseID int, adjust ...func(*ent.CourseLandingPageCreate)) {
	t.Helper()
	h.DB.Course.UpdateOneID(courseID).SetReadiness("completed").ExecX(t.Context())
	create := h.DB.CourseLandingPage.Create().
		SetCourseID(courseID).
		SetSlug("card-" + strconv.Itoa(courseID)).
		SetName("Card").
		SetMain(true).
		SetListed(true).
		SetLocale("ru").
		SetState("published")
	for _, edit := range adjust {
		edit(create)
	}
	create.SaveX(t.Context())
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
