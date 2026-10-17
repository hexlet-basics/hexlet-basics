package handlers

import (
	"context"
	"errors"

	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/coursepages"
)

// GetCourseLesson returns the lesson player payload: the lesson's prose for the
// request locale, the code and tests of the course's current build, and where
// the visitor stands in the course. What the page shows and what makes a lesson
// not found are decided in coursepages; this handler converts the result to the
// contract.
//
// Reading never starts the lesson (ADR-0012), and it answers everyone: whether
// the editor's actions are open is decided by the client from the progress
// payload.
func (s *Server) GetCourseLesson(
	ctx context.Context,
	params api.GetCourseLessonParams,
) (api.GetCourseLessonRes, error) {
	page, err := s.coursePages.Lesson(ctx, currentLearner(ctx), s.i18n.Locale(ctx), params.CourseSlug, params.Slug)
	if errors.Is(err, coursepages.ErrNotFound) {
		return s.lessonNotFound(ctx), nil
	}
	if err != nil {
		return nil, err
	}

	lessons := apiconv.ToLearnerLessons(page.Infos, page.State)
	landing := api.NilCourseLandingPage{Null: true}
	if page.Landing != nil {
		landing = api.NewNilCourseLandingPage(s.conv.ToCourseLandingPage(page.Landing))
	}
	return &api.CourseLessonView{
		Lesson: apiconv.ToCourseLesson(apiconv.LessonContent{
			Course:        page.Course,
			Lesson:        page.Lesson,
			Version:       page.Version,
			Translation:   page.Translation,
			SourceCodeURL: page.SourceURL,
		}, s.conv.ToCourse(page.Course)),
		LandingPage: landing,
		Lessons:     lessons,
		Adjacent:    apiconv.AdjacentLessons(params.Slug, lessons),
		Progress:    s.courseProgress(page.State, lessons),
	}, nil
}
