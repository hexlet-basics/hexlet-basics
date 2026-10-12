package handlers

import (
	"context"
	"strconv"

	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
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
	u, ok := AuthenticatedUser(ctx)
	if !ok {
		return nil, errUnauthenticated
	}

	courses, err := s.dashboard.Courses(ctx, u.ID, s.i18n.Locale(ctx))
	if err != nil {
		return nil, err
	}

	dashboard := &api.MyDashboard{
		StartedEnrollments:     []api.Enrollment{},
		FinishedEnrollments:    []api.Enrollment{},
		LandingPagesByCourseId: api.MyDashboardLandingPagesByCourseId{},
	}
	for _, c := range courses {
		item := apiconv.ToEnrollment(c.Enrollment, c.State)
		if c.State.State == progress.StateFinished {
			dashboard.FinishedEnrollments = append(dashboard.FinishedEnrollments, item)
		} else {
			dashboard.StartedEnrollments = append(dashboard.StartedEnrollments, item)
		}
		// Keyed by course id as a string: JSON object keys are strings, and the
		// client pairs each card with the enrollment's courseId.
		dashboard.LandingPagesByCourseId[strconv.Itoa(c.Enrollment.CourseID)] = s.conv.ToCatalogItem(c.Card)
	}
	return dashboard, nil
}
