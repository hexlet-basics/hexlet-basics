package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/languagelesson"
	"hexletbasics/ent/languagelessonmember"
	"hexletbasics/ent/languagelessonreview"
	"hexletbasics/ent/languagelessonversioninfo"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
)

// These three admin lists (legacy `admin/language_lessons`,
// `admin/language_lesson_members`, `admin/language_lesson_reviews`) all read the
// lesson/version-info graph and are READ-ONLY (index only). They follow the
// blog_posts shape: a filtered base query for the page + count, then a few
// batched by-id queries to enrich cross-table fields, rather than ent joins
// (ent has no edges declared on these read-only schemas).
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
func (s *Server) AdminListCourseLessons(ctx context.Context, params api.AdminListCourseLessonsParams) (*api.CourseLessonListItemPage, error) {
	page := newPagination(params.Page, params.PerPage)

	// `.current` = infos whose version is the CURRENT version of a COMPLETED
	// course. A language_version_id belongs to exactly one course, so matching on
	// the version id alone reproduces the legacy join without also carrying the
	// `info.language_id = courses.id` predicate.
	currentVersionIDs, err := s.completedCurrentVersionIDs(ctx)
	if err != nil {
		return nil, err
	}

	// One base builder reused for count and page so both carry the same filter;
	// a bare Count would over-report (it would ignore locale/version scoping).
	base := func() *ent.LanguageLessonVersionInfoQuery {
		return s.db.LanguageLessonVersionInfo.Query().Where(
			languagelessonversioninfo.LocaleEQ(defaultAdminLocale),
			languagelessonversioninfo.LanguageVersionIDIn(currentVersionIDs...),
		)
	}

	total, err := base().Count(ctx)
	if err != nil {
		return nil, err
	}

	// Legacy orders by the Info id (ransack `sf=id`), NOT by the lesson id that
	// this row reports as its `id`, so we order on the info's field explicitly.
	infos, err := base().
		Order(ent.Desc(languagelessonversioninfo.FieldID)).
		Offset(page.Offset()).
		Limit(page.Limit()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	lessonIDs := make([]int, len(infos))
	for i, info := range infos {
		lessonIDs[i] = info.LanguageLessonID
	}
	slugByLesson, err := s.lessonSlugs(ctx, lessonIDs)
	if err != nil {
		return nil, err
	}

	items := make([]api.CourseLessonListItem, len(infos))
	for i, info := range infos {
		items[i] = api.CourseLessonListItem{
			ID:          int32(info.LanguageLessonID),
			Name:        apiconv.NilStringFromPtr(info.Name),
			Description: apiconv.NilStringFromPtr(info.Description),
			Slug:        slugByLesson[info.LanguageLessonID],
		}
	}

	return &api.CourseLessonListItemPage{Items: items, Total: int32(total), Page: page.Page, PerPage: page.PerPage}, nil
}

// AdminListCourseLessonMembers lists per-lesson memberships, newest first
// (legacy `admin/language_lesson_members`). The course slug, lesson slug, and
// lesson name are enriched from the member's language/lesson.
func (s *Server) AdminListCourseLessonMembers(ctx context.Context, params api.AdminListCourseLessonMembersParams) (*api.CourseLessonMemberPage, error) {
	page := newPagination(params.Page, params.PerPage)

	total, err := s.db.LanguageLessonMember.Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	members, err := s.db.LanguageLessonMember.Query().
		Order(ent.Desc(languagelessonmember.FieldID)).
		Offset(page.Offset()).
		Limit(page.Limit()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	courseIDs := make([]int, 0, len(members))
	lessonIDs := make([]int, 0, len(members))
	for _, m := range members {
		courseIDs = append(courseIDs, m.LanguageID)
		lessonIDs = append(lessonIDs, m.LessonID)
	}
	courseSlugByID, err := s.courseSlugs(ctx, courseIDs)
	if err != nil {
		return nil, err
	}
	lessonSlugByID, err := s.lessonSlugs(ctx, lessonIDs)
	if err != nil {
		return nil, err
	}
	lessonNameByID, err := s.lessonNames(ctx, lessonIDs)
	if err != nil {
		return nil, err
	}

	items := make([]api.CourseLessonMember, len(members))
	for i, m := range members {
		items[i] = api.CourseLessonMember{
			ID:               int32(m.ID),
			UserId:           int32(m.UserID),
			State:            memberState(m.State),
			MessagesCount:    apiconv.NilInt32FromPtr(m.MessagesCount),
			CreatedAt:        m.CreatedAt,
			CourseSlug:       courseSlugByID[m.LanguageID],
			CourseLessonSlug: lessonSlugByID[m.LessonID],
			CourseLessonName: lessonNameByID[m.LessonID],
		}
	}

	return &api.CourseLessonMemberPage{Items: items, Total: int32(total), Page: page.Page, PerPage: page.PerPage}, nil
}

// AdminListCourseLessonReviews lists AI lesson reviews for the current locale
// that actually have a summary, newest first (legacy
// `admin/language_lesson_reviews`, `where(locale:).with_summary`). The lesson
// slug/natural-order and course slug are enriched by id.
func (s *Server) AdminListCourseLessonReviews(ctx context.Context, params api.AdminListCourseLessonReviewsParams) (*api.CourseLessonReviewPage, error) {
	page := newPagination(params.Page, params.PerPage)

	// Same filter on count and page: an empty-summary review (legacy skips these,
	// they mark lessons with no student questions) must not inflate the total.
	base := func() *ent.LanguageLessonReviewQuery {
		return s.db.LanguageLessonReview.Query().Where(
			languagelessonreview.LocaleEQ(defaultAdminLocale),
			languagelessonreview.SummaryNEQ(""),
		)
	}

	total, err := base().Count(ctx)
	if err != nil {
		return nil, err
	}

	reviews, err := base().
		Order(ent.Desc(languagelessonreview.FieldID)).
		Offset(page.Offset()).
		Limit(page.Limit()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	courseIDs := make([]int, 0, len(reviews))
	lessonIDs := make([]int, 0, len(reviews))
	for _, r := range reviews {
		courseIDs = append(courseIDs, r.LanguageID)
		lessonIDs = append(lessonIDs, r.LanguageLessonID)
	}
	courseSlugByID, err := s.courseSlugs(ctx, courseIDs)
	if err != nil {
		return nil, err
	}
	lessons, err := s.lessonsByID(ctx, lessonIDs)
	if err != nil {
		return nil, err
	}

	items := make([]api.CourseLessonReview, len(reviews))
	for i, r := range reviews {
		lesson := lessons[r.LanguageLessonID]
		items[i] = api.CourseLessonReview{
			ID:                        int32(r.ID),
			Locale:                    r.Locale,
			CourseId:                  int32(r.LanguageID),
			CourseLessonId:            int32(r.LanguageLessonID),
			CourseLessonVersionId:     int32(r.LanguageLessonVersionID),
			CourseLessonVersionInfoId: int32(r.LanguageLessonVersionInfoID),
			Summary:                   r.Summary,
			Slug:                      lessonSlug(lesson),
			LessonNaturalOrder:        lessonNaturalOrder(lesson),
			CourseSlug:                courseSlugByID[r.LanguageID],
			CreatedAt:                 r.CreatedAt,
		}
	}

	return &api.CourseLessonReviewPage{Items: items, Total: int32(total), Page: page.Page, PerPage: page.PerPage}, nil
}

// completedCurrentVersionIDs returns the current_version_id of every completed
// course (nil ids skipped). This is the version set that defines a lesson info
// as "current" for the admin lesson list.
func (s *Server) completedCurrentVersionIDs(ctx context.Context) ([]int, error) {
	courses, err := s.db.Course.Query().Where(course.ProgressEQ("completed")).All(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(courses))
	for _, c := range courses {
		if c.CurrentVersionID != nil {
			ids = append(ids, *c.CurrentVersionID)
		}
	}
	return ids, nil
}

// lessonSlugs maps lesson id -> slug for the given lesson ids (missing/null
// slugs are absent from the map, yielding the zero string).
func (s *Server) lessonSlugs(ctx context.Context, ids []int) (map[int]string, error) {
	lessons, err := s.lessonsByID(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make(map[int]string, len(lessons))
	for id, l := range lessons {
		out[id] = lessonSlug(l)
	}
	return out, nil
}

// lessonNames maps lesson id -> its name for the default locale. Legacy is
// locale-agnostic here (`localed_info = infos.first`, the lowest-id info across
// all versions/locales); we instead take the default-locale info, consistent
// with the pinned admin locale above. A lesson with no such info yields "".
func (s *Server) lessonNames(ctx context.Context, lessonIDs []int) (map[int]string, error) {
	if len(lessonIDs) == 0 {
		return map[int]string{}, nil
	}
	infos, err := s.db.LanguageLessonVersionInfo.Query().
		Where(
			languagelessonversioninfo.LocaleEQ(defaultAdminLocale),
			languagelessonversioninfo.LanguageLessonIDIn(lessonIDs...),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[int]string, len(infos))
	for _, info := range infos {
		if _, seen := out[info.LanguageLessonID]; seen {
			continue // a lesson can have several versions' infos; keep the first.
		}
		out[info.LanguageLessonID] = apiconv.StringFromPtr(info.Name)
	}
	return out, nil
}

// lessonsByID batch-loads lessons keyed by id (used for slug + natural order).
func (s *Server) lessonsByID(ctx context.Context, ids []int) (map[int]*ent.LanguageLesson, error) {
	if len(ids) == 0 {
		return map[int]*ent.LanguageLesson{}, nil
	}
	lessons, err := s.db.LanguageLesson.Query().
		Where(languagelesson.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[int]*ent.LanguageLesson, len(lessons))
	for _, l := range lessons {
		out[l.ID] = l
	}
	return out, nil
}

// courseSlugs maps course (language) id -> slug for the given ids.
func (s *Server) courseSlugs(ctx context.Context, ids []int) (map[int]string, error) {
	if len(ids) == 0 {
		return map[int]string{}, nil
	}
	courses, err := s.db.Course.Query().
		Where(course.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[int]string, len(courses))
	for _, c := range courses {
		out[c.ID] = apiconv.StringFromPtr(c.Slug)
	}
	return out, nil
}

// lessonSlug / lessonNaturalOrder safely read a (possibly missing) lesson's
// nullable columns. The contract types natural order as non-null, but the
// baseline column is nullable, so a missing lesson or NULL reads as 0.
func lessonSlug(l *ent.LanguageLesson) string {
	if l == nil {
		return ""
	}
	return apiconv.StringFromPtr(l.Slug)
}

func lessonNaturalOrder(l *ent.LanguageLesson) int32 {
	if l == nil || l.NaturalOrder == nil {
		return 0
	}
	return int32(*l.NaturalOrder)
}

// memberState bridges the nullable state column to the API enum. A started
// membership is the initial AASM state, so a NULL/blank column reads as started.
func memberState(v *string) api.MemberState {
	if v == nil || *v == "" {
		return api.MemberStateStarted
	}
	return api.MemberState(*v)
}
