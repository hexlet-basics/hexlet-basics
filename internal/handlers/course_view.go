package handlers

import (
	"context"

	"hexletbasics/internal/progress"

	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
)

// GetCourse returns everything the public course landing page needs: the
// course, its landing copy, the lessons of its current version, and where the
// visitor stands. What the page shows and when it is not found are decided in
// coursepages; this handler converts the result to the contract.
//
// Each lesson in the list carries the visitor's standing in it, so the page
// joins nothing; the course's own standing (Completion, the Next Lesson) travels
// as the progress payload.
//
// An anonymous request is not an error: it returns the same payload with no
// enrollment, which is exactly what the page renders for a visitor today.
func (s *Server) GetCourse(ctx context.Context, params api.GetCourseParams) (api.GetCourseRes, error) {
	page, err := s.coursePages.Course(ctx, currentLearner(ctx), s.i18n.Locale(ctx), params.Slug)
	if err != nil {
		return nil, err
	}

	lessons := apiconv.ToLearnerLessons(page.Infos, page.State)
	view := &api.CourseView{
		Course:      s.conv.ToCourse(page.Course),
		LandingPage: api.NewNilCourseLandingPage(s.conv.ToCourseLandingPage(page.Landing)),
		MovedTo:     api.NilString{Null: true},
		Lessons:     lessons,
		Modules:     apiconv.ToCourseModules(page.Modules, lessons),
		QnaItems:    s.conv.ToLandingPageQnaItems(page.QnaItems),
		Enrollment:  api.NilEnrollment{Null: true},
		Progress:    s.courseProgress(page.State, lessons),
	}
	if page.MovedTo != nil {
		view.MovedTo = api.NewNilString(*page.MovedTo)
	}
	if page.Enrollment != nil {
		view.Enrollment = api.NewNilEnrollment(apiconv.ToEnrollment(page.Enrollment))
	}
	return view, nil
}

// courseProgress is the visitor's standing as the contract carries it: null for
// a course with no lessons, where there is nothing to stand in.
func (s *Server) courseProgress(state *progress.CourseState, lessons []api.LearnerLesson) api.NilCourseProgress {
	if len(state.Lessons) == 0 {
		return api.NilCourseProgress{Null: true}
	}
	return api.NewNilCourseProgress(
		apiconv.ToCourseProgress(state, apiconv.NextLessonRef(state.NextLessonSlug, lessons)),
	)
}
