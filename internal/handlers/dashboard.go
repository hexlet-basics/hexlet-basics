package handlers

import (
	"context"

	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/dashboard"
)

// GetMyDashboard returns the learner's courses, split into what they are
// working on and what they have finished, each carrying the same progress the
// course page reports — it is computed by the same module, so the two cannot
// disagree.
//
// Authenticated only, and deliberately so: this is where the incentive to
// create an account lives, and a guest carrying a progress cookie is refused
// rather than served their own dashboard. The contract declares the session, so
// ogen refuses the request before this handler runs.
//
// Only courses with a main catalogue entry appear, matching the legacy scope:
// the dashboard renders course cards, so a course with none is not listed and
// each listed course is complete.
func (s *Server) GetMyDashboard(ctx context.Context) (api.GetMyDashboardRes, error) {
	u := MustAuthenticatedUser(ctx)

	courses, err := s.dashboard.Courses(ctx, u.ID, s.i18n.Locale(ctx))
	if err != nil {
		return nil, err
	}

	return &api.MyDashboard{
		Started:  s.dashboardCourses(courses.Started),
		Finished: s.dashboardCourses(courses.Finished),
	}, nil
}

// dashboardCourses projects each course: the enrollment, where the learner
// stands in it, and the card to draw.
func (s *Server) dashboardCourses(courses []dashboard.Course) []api.DashboardCourse {
	items := make([]api.DashboardCourse, 0, len(courses))
	for _, c := range courses {
		next := api.NilLessonRef{Null: true}
		if c.NextLesson != nil {
			next = apiconv.NamedLessonRef(c.State.NextLessonSlug, c.NextLesson.ID, &c.NextLesson.Name)
		}
		items = append(items, api.DashboardCourse{
			Enrollment: apiconv.ToEnrollment(c.Enrollment),
			Progress:   apiconv.ToCourseProgress(c.State, next),
			Card:       s.conv.ToCatalogItem(c.Card),
		})
	}
	return items
}
