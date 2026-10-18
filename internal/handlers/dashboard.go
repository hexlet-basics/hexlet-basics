package handlers

import (
	"context"

	"hexletbasics/internal/api"
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
// Only courses with a main catalogue entry appear, matching the legacy scope:
// the dashboard renders course cards, so a course with none is not listed and
// each listed course is complete.
func (s *Server) GetMyDashboard(ctx context.Context) (*api.MyDashboard, error) {
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

	return &api.MyDashboard{
		Started:  s.conv.ToDashboardCourses(courses.Started),
		Finished: s.conv.ToDashboardCourses(courses.Finished),
	}, nil
}

// MyCourseProgress is where the user stands in one Course in the glossary's
// three conditions, for the MCP learner tools. The progress module computes it
// from the same rows as the course page, with Started added.
func (s *Server) MyCourseProgress(ctx context.Context, userID int, courseSlug string) (*progress.LearnerCourse, error) {
	return s.progress.LearnerCourse(ctx, userID, courseSlug)
}
