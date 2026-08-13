package handlers

import (
	"context"
	"net/url"
	"path"
	"strings"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
)

// GetCourseLesson returns the lesson player payload: the lesson's prose for the
// request locale, the code and tests of the course's current build, and where
// the visitor stands in the course.
//
// Reading never starts the lesson. Legacy enrolled a learner as a side effect
// of rendering this page; the frontend now preloads routes on hover, so a read
// with that effect would enroll them in every lesson they pointed at and emit a
// fact for each (ADR-0012). Starting is its own command.
//
// It answers everyone, including a visitor on a lesson they may not take yet:
// theory is public and is what search engines index. Whether the editor's
// actions are open is decided from the progress payload, by the client, from
// the server's answer.
func (s *Server) GetCourseLesson(
	ctx context.Context,
	params api.GetCourseLessonParams,
) (api.GetCourseLessonRes, error) {
	crs, err := s.db.Course.Query().
		Where(course.SlugEQ(params.CourseSlug)).
		WithCurrentVersion().
		Only(ctx)
	if ent.IsNotFound(err) {
		return s.lessonNotFound(ctx), nil
	}
	if err != nil {
		return nil, err
	}
	if crs.CurrentVersionID == nil {
		// Nothing is built, so there is no lesson content to serve.
		return s.lessonNotFound(ctx), nil
	}

	lesson, err := s.db.CourseLesson.Query().
		Where(
			courselesson.SlugEQ(params.Slug),
			courselesson.CourseID(crs.ID),
		).
		Only(ctx)
	if ent.IsNotFound(err) {
		return s.lessonNotFound(ctx), nil
	}
	if err != nil {
		return nil, err
	}

	// Resolved against the CURRENT version: a lesson only a retired build
	// contained is gone from the course, and serving its old content would show
	// a page whose code and tests no longer exist.
	version, err := s.db.CourseLessonVersion.Query().
		Where(
			courselessonversion.LessonID(lesson.ID),
			courselessonversion.CourseVersionID(*crs.CurrentVersionID),
		).
		Only(ctx)
	if ent.IsNotFound(err) {
		return s.lessonNotFound(ctx), nil
	}
	if err != nil {
		return nil, err
	}

	locale := s.i18n.Locale(ctx)
	translation, err := s.db.CourseLessonTranslation.Query().
		Where(
			courselessontranslation.CourseLessonID(lesson.ID),
			courselessontranslation.CourseVersionID(*crs.CurrentVersionID),
			courselessontranslation.LocaleEQ(locale),
		).
		Only(ctx)
	if ent.IsNotFound(err) {
		// A lesson not translated into this locale is not a lesson this visitor
		// can read, which is how legacy treated it too.
		return s.lessonNotFound(ctx), nil
	}
	if err != nil {
		return nil, err
	}

	state, err := s.progress.CourseState(ctx, currentLearner(ctx), crs.ID)
	if err != nil {
		return nil, err
	}

	// The same page the course landing read resolves, because the player titles
	// this page with the landing copy's name rather than the course's own.
	landing, err := s.mainLandingPage(ctx, crs)
	if err != nil {
		return nil, err
	}

	// The same list the course landing read returns: the player's navigation tab
	// joins it to the progress payload by slug, so the two screens name and order
	// the lessons identically.
	lessons, err := s.currentLessonList(ctx, crs)
	if err != nil {
		return nil, err
	}

	view := &api.CourseLessonView{
		Lesson: apiconv.ToCourseLesson(apiconv.LessonContent{
			Course:        crs,
			Lesson:        lesson,
			Version:       version,
			Translation:   translation,
			SourceCodeURL: s.lessonSourceURL(version, locale),
		}, s.conv.ToCourse(crs)),
		LandingPage: landing,
		Lessons:     lessons,
		Progress:    api.NilCourseProgress{Null: true},
	}
	if len(state.Lessons) > 0 {
		view.Progress = api.NewNilCourseProgress(apiconv.ToCourseProgress(state))
	}
	return view, nil
}

// lessonSourceURL points at the lesson's README in the content repository, so a
// reader can propose a fix to what they are reading.
//
// `path_to_code` is repository-relative and starts at the course repo, so the
// link is the org base plus that path plus the locale's README — with `modules`
// rewritten to GitHub's blob path, exactly as the legacy helper did. A lesson
// whose version carries no path has no link rather than a broken one.
func (s *Server) lessonSourceURL(version *ent.CourseLessonVersion, locale string) string {
	if version.PathToCode == nil || *version.PathToCode == "" {
		return ""
	}

	base, err := url.JoinPath(s.cfg.CourseRepoBaseURL, *version.PathToCode, locale, "README.md")
	if err != nil {
		return ""
	}
	return strings.Replace(base, path.Join("/", "modules"), "/blob/main/modules", 1)
}
