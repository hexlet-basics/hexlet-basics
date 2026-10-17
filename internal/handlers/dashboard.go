package handlers

import (
	"context"
	"strconv"

	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/dashboard"
	"hexletbasics/internal/progress"
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
// Only courses with a listed, main catalogue entry appear, matching the legacy
// scope: the dashboard renders course cards, and a course with no card to
// render would be a row the client cannot draw.
func (s *Server) GetMyDashboard(ctx context.Context) (api.GetMyDashboardRes, error) {
	dashboard, err := s.MyDashboard(ctx, MustAuthenticatedUser(ctx).ID)
	if err != nil {
		return nil, err
	}
	return dashboard, nil
}

// MyDashboard is the dashboard of one user, for callers that identified the
// user themselves: the MCP learner tool takes the id from its bearer token,
// where the HTTP operation takes it from the session. It is the same code, so
// the two cannot answer differently.
func (s *Server) MyDashboard(ctx context.Context, userID int) (*api.MyDashboard, error) {
	courses, err := s.dashboard.Courses(ctx, userID, s.i18n.Locale(ctx))
	if err != nil {
		return nil, err
	}

	dashboard := &api.MyDashboard{
		StartedEnrollments:     []api.Enrollment{},
		FinishedEnrollments:    []api.Enrollment{},
		LandingPagesByCourseId: api.MyDashboardLandingPagesByCourseId{},
	}
	for _, c := range courses.Started {
		dashboard.StartedEnrollments = append(dashboard.StartedEnrollments, apiconv.ToEnrollment(c.Enrollment, c.State))
		s.addCard(dashboard, c)
	}
	for _, c := range courses.Finished {
		dashboard.FinishedEnrollments = append(dashboard.FinishedEnrollments, apiconv.ToEnrollment(c.Enrollment, c.State))
		s.addCard(dashboard, c)
	}
	return dashboard, nil
}

// addCard keys the course's card by course id as a string: JSON object keys are
// strings, and the client pairs each card with the enrollment's courseId.
func (s *Server) addCard(out *api.MyDashboard, c dashboard.Course) {
	out.LandingPagesByCourseId[strconv.Itoa(c.Enrollment.CourseID)] = s.conv.ToCatalogItem(c.Card)
}

// MyCourseProgress is where the user stands in one Course in the glossary's
// three conditions, for the MCP learner tools. The progress module computes it
// from the same rows as the course page, with Started added.
func (s *Server) MyCourseProgress(ctx context.Context, userID int, courseSlug string) (*progress.LearnerCourse, error) {
	return s.progress.LearnerCourse(ctx, userID, courseSlug)
}
