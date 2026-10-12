package handlers

import (
	"context"
	"strconv"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/enrollment"
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

	// The card is the course's main landing page in the request locale, legacy
	// `LandingPage.web.where(main: true)`: the dashboard scope, which unlike the
	// public lists does not require `listed`. One query serves every enrollment.
	cards, err := s.catalog.MainPages(ctx, s.i18n.Locale(ctx),
		lo.Map(enrollments, func(e *ent.Enrollment, _ int) int { return e.CourseID })...)
	if err != nil {
		return nil, err
	}

	states, err := s.progress.CourseStates(ctx, u.ID,
		lo.Map(enrollments, func(e *ent.Enrollment, _ int) int { return e.CourseID }))
	if err != nil {
		return nil, err
	}

	for _, enrolled := range enrollments {
		page, ok := cards[enrolled.CourseID]
		if !ok {
			continue
		}
		state := states[enrolled.CourseID]

		item := apiconv.ToEnrollment(enrolled, state)
		if state.State == progress.StateFinished {
			dashboard.FinishedEnrollments = append(dashboard.FinishedEnrollments, item)
		} else {
			dashboard.StartedEnrollments = append(dashboard.StartedEnrollments, item)
		}
		// Keyed by course id as a string: JSON object keys are strings, and the
		// client pairs each card with the enrollment's courseId.
		dashboard.LandingPagesByCourseId[strconv.Itoa(enrolled.CourseID)] = s.conv.ToCatalogItem(page)
	}
	return dashboard, nil
}
