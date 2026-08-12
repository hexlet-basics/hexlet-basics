package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courseversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
)

// GetCourse returns everything the public course landing page needs: the
// course, its landing copy, the lessons of its current version, and — for a
// signed-in learner — where they stand.
//
// The lesson list model stays the shared one; per-learner state travels in the
// progress payload instead, because the same list model feeds the admin lesson
// index where a learner's state is meaningless.
//
// An anonymous request is not an error: it returns the same payload with no
// enrollment, which is exactly what the page renders for a visitor today.
func (s *Server) GetCourse(ctx context.Context, params api.GetCourseParams) (api.GetCourseRes, error) {
	crs, err := s.db.Course.Query().
		Where(course.SlugEQ(params.Slug)).
		WithCurrentVersion().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	landing, err := s.mainLandingPage(ctx, crs)
	if err != nil {
		return nil, err
	}

	lessons, err := s.currentLessonList(ctx, crs)
	if err != nil {
		return nil, err
	}

	view := &api.CourseView{
		Course:      s.conv.ToCourse(crs),
		LandingPage: landing,
		Lessons:     lessons,
		Enrollment:  api.NilEnrollment{Null: true},
		Progress:    api.NilCourseProgress{Null: true},
	}

	// Everyone gets a position, because everyone has one: a visitor who has
	// finished nothing stands at the beginning with the first lesson open, and
	// that is a fact about the course, not about having an account. Deriving it
	// client-side instead would be a second implementation of the gate.
	learner := submittingLearner(ctx)
	state, err := s.progress.CourseState(ctx, learner, crs.ID)
	if err != nil {
		return nil, err
	}
	if len(state.Lessons) > 0 {
		view.Progress = api.NewNilCourseProgress(apiconv.ToCourseProgress(state))
	}

	if !learner.SignedIn() {
		return view, nil
	}

	enrolled, err := s.db.Enrollment.Query().
		Where(enrollment.UserID(learner.UserID), enrollment.CourseID(crs.ID)).
		Only(ctx)
	switch {
	case ent.IsNotFound(err):
		// Not started: the absence of a record, not a state a record can be in.
		return view, nil
	case err != nil:
		return nil, err
	}
	view.Enrollment = api.NewNilEnrollment(apiconv.ToEnrollment(enrolled, state))
	return view, nil
}

// mainLandingPage is the course's own landing copy: the page flagged main, and
// otherwise the oldest one, mirroring how legacy picks the canonical page.
func (s *Server) mainLandingPage(ctx context.Context, crs *ent.Course) (api.NilCourseLandingPage, error) {
	page, err := s.db.LandingPage.Query().
		Where(landingpage.CourseID(crs.ID)).
		Order(ent.Desc(landingpage.FieldMain), ent.Asc(landingpage.FieldID)).
		WithCourse().
		First(ctx)
	switch {
	case ent.IsNotFound(err):
		return api.NilCourseLandingPage{Null: true}, nil
	case err != nil:
		return api.NilCourseLandingPage{}, err
	}
	return api.NewNilCourseLandingPage(s.conv.ToCourseLandingPage(page)), nil
}

// currentLessonList is the lessons of the course's current version, in the
// request locale, ordered as the course orders them.
func (s *Server) currentLessonList(ctx context.Context, crs *ent.Course) ([]api.CourseLessonListItem, error) {
	if crs.CurrentVersionID == nil {
		return []api.CourseLessonListItem{}, nil
	}

	infos, err := s.db.CourseLessonTranslation.Query().
		Where(
			courselessontranslation.LocaleEQ(s.i18n.Locale(ctx)),
			courselessontranslation.HasCourseVersionWith(courseversion.ID(*crs.CurrentVersionID)),
		).
		WithLesson().
		Order(ent.Asc(courselessontranslation.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return s.conv.ToCourseLessonListItems(infos), nil
}
