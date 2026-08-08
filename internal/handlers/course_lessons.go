package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselessonreview"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courseversion"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/api"
	"hexletbasics/internal/localization"
)

// These three admin lists (legacy `admin/language_lessons`,
// `admin/language_lesson_members`, `admin/language_lesson_reviews`) all read the
// lesson/version-info graph and are READ-ONLY (index only). Their ent queries
// eager-load the graph needed by the API read models; apiconv owns the mapping
// from that graph to generated API types.
//
// Locale: legacy filters these by the request's `I18n.locale` (the admin UI
// language). There is no request locale at the ogen handler boundary yet (the
// known admin-locale design gap), so we pin the default locale here. When
// request locale reaches handlers, thread it through in place of this constant.
const defaultAdminLocale = "en"

// AdminListCourseLessons lists lessons as the admin index shows them: one row
// per lesson for the current locale, drawn from the current version's info
// (legacy `Language::Lesson::Version::Info.current.with_locale`). name and
// description come from the info row; the slug is joined from the lesson.
func (s *Server) AdminListCourseLessons(ctx context.Context, params api.AdminListCourseLessonsParams) (api.AdminListCourseLessonsRes, error) {
	page := newPagination(params.Page, params.PerPage)

	// One base builder reused for count and page so both carry the same filter;
	// a bare Count would over-report (it would ignore locale/version scoping).
	base := func() *ent.CourseLessonTranslationQuery {
		return s.db.CourseLessonTranslation.Query().Where(
			courselessontranslation.LocaleEQ(defaultAdminLocale),
			// Legacy `.current`: the info's course version must be selected as
			// current by a completed course.
			courselessontranslation.HasCourseVersionWith(
				courseversion.HasCurrentCoursesWith(course.ProgressEQ("completed")),
			),
		)
	}

	total, err := base().Count(ctx)
	if err != nil {
		return nil, err
	}

	// Legacy orders by the Info id (ransack `sf=id`), NOT by the lesson id that
	// this row reports as its `id`, so we order on the info's field explicitly.
	infos, err := base().
		WithLesson().
		Order(ent.Desc(courselessontranslation.FieldID)).
		Offset(page.Offset()).
		Limit(page.Limit()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return &api.CourseLessonListItemPage{
		Items:   s.conv.ToCourseLessonListItems(infos),
		Total:   int32(total),
		Page:    page.Page,
		PerPage: page.PerPage,
	}, nil
}

// AdminListCourseLessonMembers lists per-lesson memberships, newest first
// (legacy `admin/language_lesson_members`). The course slug, lesson slug, and
// lesson name are enriched from the member's language/lesson.
func (s *Server) AdminListCourseLessonMembers(ctx context.Context, params api.AdminListCourseLessonMembersParams) (api.AdminListCourseLessonMembersRes, error) {
	page := newPagination(params.Page, params.PerPage)

	total, err := s.db.LessonProgress.Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	members, err := s.db.LessonProgress.Query().
		WithCourse().
		WithLesson(func(q *ent.CourseLessonQuery) {
			q.WithInfos(func(q *ent.CourseLessonTranslationQuery) {
				q.Where(courselessontranslation.LocaleEQ(defaultAdminLocale)).
					Order(ent.Asc(courselessontranslation.FieldID))
			})
		}).
		Order(ent.Desc(lessonprogress.FieldID)).
		Offset(page.Offset()).
		Limit(page.Limit()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return &api.CourseLessonMemberPage{
		Items:   s.conv.ToCourseLessonMembers(members),
		Total:   int32(total),
		Page:    page.Page,
		PerPage: page.PerPage,
	}, nil
}

// AdminListCourseLessonReviews lists AI lesson reviews for the current locale
// that actually have a summary, newest first (legacy
// `admin/language_lesson_reviews`, `where(locale:).with_summary`). The lesson
// slug/natural-order and course slug are enriched by id.
func (s *Server) AdminListCourseLessonReviews(ctx context.Context, params api.AdminListCourseLessonReviewsParams) (api.AdminListCourseLessonReviewsRes, error) {
	page := newPagination(params.Page, params.PerPage)

	// Same filter on count and page: an empty-summary review (legacy skips these,
	// they mark lessons with no student questions) must not inflate the total.
	base := func() *ent.CourseLessonReviewQuery {
		return s.db.CourseLessonReview.Query().Where(
			courselessonreview.LocaleEQ(defaultAdminLocale),
			courselessonreview.SummaryNEQ(""),
		)
	}

	total, err := base().Count(ctx)
	if err != nil {
		return nil, err
	}

	reviews, err := base().
		WithCourse().
		WithLesson().
		Order(ent.Desc(courselessonreview.FieldID)).
		Offset(page.Offset()).
		Limit(page.Limit()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return &api.CourseLessonReviewPage{
		Items:   s.conv.ToCourseLessonReviews(reviews),
		Total:   int32(total),
		Page:    page.Page,
		PerPage: page.PerPage,
	}, nil
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
		Where(courselessontranslation.LanguageLessonID(lesson.ID)).
		Order(ent.Asc(courselessontranslation.FieldID)).
		IDs(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.reviews.EnqueueLessonReviews(ctx, infoIDs); err != nil {
		return nil, err
	}
	return &api.AdminReviewCourseLessonNoContent{}, nil
}
