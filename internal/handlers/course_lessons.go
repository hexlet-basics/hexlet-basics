package handlers

import (
	"context"
	"errors"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/domain"
)

// These three admin lists (legacy `admin/course_lessons`,
// `admin/lesson_progress`, `admin/course_lesson_reviews`) read the
// lesson/version-info graph and are READ-ONLY (index only). The courselessons
// module owns what each lists; apiconv maps the graph it returns to the API
// types.
//
// Locale: legacy filters these by the request's `I18n.locale` (the admin UI
// language), which is `s.i18n.Locale(ctx)` here.

// AdminListCourseLessons lists lessons as the admin index shows them: one row
// per lesson for the request locale, from the current version's translation.
func (s *Server) AdminListCourseLessons(ctx context.Context, params api.AdminListCourseLessonsParams) (*api.CourseLessonListItemPage, error) {
	locale := s.i18n.Locale(ctx)
	return listFrom(ctx, params.Page, params.PerPage,
		func(ctx context.Context, page domain.Page) ([]*ent.CourseLessonTranslation, int, error) {
			return s.lessonDirectory.Lessons(ctx, locale, params.Q.Value, page)
		},
		s.conv.ToCourseLessonListItems,
		func(items []api.CourseLessonListItem, total, page, perPage int32) *api.CourseLessonListItemPage {
			return &api.CourseLessonListItemPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminListLessonProgress lists per-lesson progress rows, newest first, in the
// request locale.
func (s *Server) AdminListLessonProgress(ctx context.Context, params api.AdminListLessonProgressParams) (*api.LessonProgressPage, error) {
	locale := s.i18n.Locale(ctx)
	return listFrom(ctx, params.Page, params.PerPage,
		func(ctx context.Context, page domain.Page) ([]*ent.LessonProgress, int, error) {
			return s.lessonDirectory.Progress(ctx, locale, page)
		},
		s.conv.ToLessonProgressList,
		func(items []api.LessonProgress, total, page, perPage int32) *api.LessonProgressPage {
			return &api.LessonProgressPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminListCourseLessonReviews lists the AI lesson reviews of the request locale
// that have a summary, newest first.
func (s *Server) AdminListCourseLessonReviews(ctx context.Context, params api.AdminListCourseLessonReviewsParams) (*api.CourseLessonReviewPage, error) {
	locale := s.i18n.Locale(ctx)
	return listFrom(ctx, params.Page, params.PerPage,
		func(ctx context.Context, page domain.Page) ([]*ent.CourseLessonReview, int, error) {
			return s.lessonDirectory.Reviews(ctx, locale, page)
		},
		s.conv.ToCourseLessonReviews,
		func(items []api.CourseLessonReview, total, page, perPage int32) *api.CourseLessonReviewPage {
			return &api.CourseLessonReviewPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminReviewCourseLesson asks for an AI review of every translation of one
// lesson: all locales and all versions, so the historical summaries refresh too.
func (s *Server) AdminReviewCourseLesson(ctx context.Context, params api.AdminReviewCourseLessonParams) error {
	err := s.lessonReviews.ReviewLesson(ctx, int(params.ID))
	if errors.Is(err, domain.ErrNotFound) {
		return errLessonNotFound
	}
	if err != nil {
		return err
	}
	return nil
}
