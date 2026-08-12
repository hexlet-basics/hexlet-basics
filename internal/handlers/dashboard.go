package handlers

import (
	"context"
	"strconv"

	"hexletbasics/ent"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/landingpage"
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

	// Newest enrollment first: a learner returning after an absence finds what
	// they most recently touched at the top, and the order is total.
	enrollments, err := s.db.Enrollment.Query().
		Where(enrollment.UserID(u.ID)).
		Order(ent.Desc(enrollment.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	dashboard := &api.MyDashboard{
		StartedEnrollments:     []api.Enrollment{},
		FinishedEnrollments:    []api.Enrollment{},
		LandingPagesByCourseId: api.MyDashboardLandingPagesByCourseId{},
	}

	for _, enrolled := range enrollments {
		card, err := s.catalogueCard(ctx, enrolled.CourseID)
		if err != nil {
			return nil, err
		}
		if card == nil {
			continue
		}

		state, err := s.progress.CourseState(ctx, u.ID, enrolled.CourseID)
		if err != nil {
			return nil, err
		}

		item := apiconv.ToEnrollment(enrolled, state)
		if state.State == progress.StateFinished {
			dashboard.FinishedEnrollments = append(dashboard.FinishedEnrollments, item)
		} else {
			dashboard.StartedEnrollments = append(dashboard.StartedEnrollments, item)
		}
		// Keyed by course id as a string: JSON object keys are strings, and the
		// client pairs each card with the enrollment's courseId.
		dashboard.LandingPagesByCourseId[strconv.Itoa(enrolled.CourseID)] = *card
	}
	return dashboard, nil
}

// catalogueCard is the course's listed main landing page, or nil when it has
// none — the legacy dashboard scope.
func (s *Server) catalogueCard(ctx context.Context, courseID int) (*api.CourseCatalogItem, error) {
	page, err := s.db.LandingPage.Query().
		Where(
			landingpage.CourseID(courseID),
			landingpage.MainEQ(true),
			landingpage.ListedEQ(true),
		).
		Order(ent.Asc(landingpage.FieldID)).
		WithCourse().
		First(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, nil
	case err != nil:
		return nil, err
	}

	item := s.conv.ToCatalogItem(page)
	return &item, nil
}
