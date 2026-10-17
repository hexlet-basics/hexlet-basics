package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/courselessonreview"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courseversion"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/api"
	"hexletbasics/internal/landingpages"
	"hexletbasics/internal/localization"
)

// These three admin lists (legacy `admin/course_lessons`,
// `admin/lesson_progress`, `admin/course_lesson_reviews`) all read the
// lesson/version-info graph and are READ-ONLY (index only). Their ent queries
// eager-load the graph needed by the API read models; apiconv owns the mapping
// from that graph to generated API types.
//
//
// Locale: legacy filters these by the request's `I18n.locale` (the admin UI
// language), which is `s.i18n.Locale(ctx)` here.

// AdminListCourseLessons lists lessons as the admin index shows them: one row
// per lesson for the current locale, drawn from the current version's info
// (legacy `Language::Lesson::Version::Info.current.with_locale`). name and
// description come from the info row; the slug is joined from the lesson.
func (s *Server) AdminListCourseLessons(ctx context.Context, params api.AdminListCourseLessonsParams) (api.AdminListCourseLessonsRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		s.db.CourseLessonTranslation.Query().
			Where(
				courselessontranslation.LocaleEQ(s.i18n.Locale(ctx)),
				// Legacy `.current`: the info's course version must be selected as
				// current by a completed course.
				courselessontranslation.HasCourseVersionWith(
					courseversion.HasCurrentCoursesWith(landingpages.Completed()),
				),
			).
			WithLesson().
			// Legacy orders by the Info id (ransack `sf=id`), NOT by the lesson id
			// that this row reports as its `id`, so we order on the info's field.
			Order(ent.Desc(courselessontranslation.FieldID)),
		s.conv.ToCourseLessonListItems,
		func(items []api.CourseLessonListItem, total, page, perPage int32) *api.CourseLessonListItemPage {
			return &api.CourseLessonListItemPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminListLessonProgress lists per-lesson progress rows, newest first (legacy
// `admin/lesson_progress`). The course slug, lesson
// slug, and lesson name are enriched from the row's course/lesson.
func (s *Server) AdminListLessonProgress(ctx context.Context, params api.AdminListLessonProgressParams) (api.AdminListLessonProgressRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		s.db.LessonProgress.Query().
			WithCourse().
			WithLesson(func(q *ent.CourseLessonQuery) {
				q.WithInfos(func(q *ent.CourseLessonTranslationQuery) {
					q.Where(courselessontranslation.LocaleEQ(s.i18n.Locale(ctx))).
						Order(ent.Asc(courselessontranslation.FieldID))
				})
			}).
			Order(ent.Desc(lessonprogress.FieldID)),
		s.conv.ToLessonProgressList,
		func(items []api.LessonProgress, total, page, perPage int32) *api.LessonProgressPage {
			return &api.LessonProgressPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminListCourseLessonReviews lists AI lesson reviews for the current locale
// that actually have a summary, newest first (legacy
// `admin/course_lesson_reviews`, `where(locale:).with_summary`). The lesson
// slug/natural-order and course slug are enriched by id.
func (s *Server) AdminListCourseLessonReviews(ctx context.Context, params api.AdminListCourseLessonReviewsParams) (api.AdminListCourseLessonReviewsRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		// An empty-summary review (legacy skips these, they mark lessons with no
		// student questions) must not inflate the total, so the filter lives on
		// the one query the count and the page share.
		s.db.CourseLessonReview.Query().
			Where(
				courselessonreview.LocaleEQ(s.i18n.Locale(ctx)),
				courselessonreview.SummaryNEQ(""),
			).
			WithCourse().
			WithLesson().
			Order(ent.Desc(courselessonreview.FieldID)),
		s.conv.ToCourseLessonReviews,
		func(items []api.CourseLessonReview, total, page, perPage int32) *api.CourseLessonReviewPage {
			return &api.CourseLessonReviewPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminReviewCourseLesson enqueues an AI review job for every info of one
// lesson — all locales AND all versions, so historical summaries refresh too
// (legacy `lesson.infos.find_each { ReviewLessonJob.perform_later }`).
func (s *Server) AdminReviewCourseLesson(ctx context.Context, params api.AdminReviewCourseLessonParams) (api.AdminReviewCourseLessonRes, error) {
	lesson, err := s.db.CourseLesson.Get(ctx, int(params.ID))
	if ent.IsNotFound(err) {
		return &api.NotFoundError{Message: s.i18n.Text(ctx, localization.LessonNotFound)}, nil
	}
	if err != nil {
		return nil, err
	}

	infoIDs, err := s.db.CourseLessonTranslation.Query().
		Where(courselessontranslation.CourseLessonID(lesson.ID)).
		Order(ent.Asc(courselessontranslation.FieldID)).
		IDs(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.enqueueLessonReviews(ctx, infoIDs); err != nil {
		return nil, err
	}
	return &api.AdminReviewCourseLessonNoContent{}, nil
}
