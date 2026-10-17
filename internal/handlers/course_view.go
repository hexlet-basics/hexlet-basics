package handlers

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/courseversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/landingpageqnaitem"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/curriculum"
)

// GetCourse returns everything the public course landing page needs: the
// course, its landing copy, the lessons of its current version, and — for a
// signed-in learner — where they stand.
//
// The slug is a landing page's, not the course's, in the request locale: that
// is what legacy's `languages#show` resolved (`LandingPage.with_locale.find_by
// (slug:)`), and what every public link to the page carries (`/languages/
// javascript-ru`), so existing URLs keep working (ADR-0002). A course has one
// landing page per locale, which is why the same slug under another locale is a
// 404 rather than the other locale's copy. The page shown is the one the URL
// names, not the course's main page, so a secondary page keeps its own copy.
//
// Each lesson in the list carries the visitor's standing in it, so the page
// joins nothing; the course's own standing (Completion, the Next Lesson) travels
// as the progress payload.
//
// An anonymous request is not an error: it returns the same payload with no
// enrollment, which is exactly what the page renders for a visitor today.
func (s *Server) GetCourse(ctx context.Context, params api.GetCourseParams) (api.GetCourseRes, error) {
	page, err := s.db.CourseLandingPage.Query().
		Where(
			courselandingpage.SlugEQ(params.Slug),
			courselandingpage.LocaleEQ(s.i18n.Locale(ctx)),
		).
		WithCourse(func(q *ent.CourseQuery) { q.WithCurrentVersion() }).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	crs := page.Edges.Course
	landing := api.NewNilCourseLandingPage(s.conv.ToCourseLandingPage(page))

	movedTo, err := s.archivedPageTarget(ctx, page)
	if err != nil {
		return nil, err
	}

	// Everyone gets a standing, because everyone has one: a visitor who has
	// finished nothing stands at the beginning, and that is a fact about the
	// course, not about having an account. Deriving it client-side instead would
	// be a second implementation of what finished means.
	learner := currentLearner(ctx)
	state, err := s.progress.CourseState(ctx, learner, crs.ID)
	if err != nil {
		return nil, err
	}

	infos, err := s.currentLessonInfos(ctx, crs)
	if err != nil {
		return nil, err
	}
	lessons := apiconv.ToLearnerLessons(infos, state)

	moduleRows, err := s.currentModuleRows(ctx, crs)
	if err != nil {
		return nil, err
	}

	qnaItems, err := s.landingQnaItems(ctx, landing)
	if err != nil {
		return nil, err
	}

	view := &api.CourseView{
		Course:      s.conv.ToCourse(crs),
		LandingPage: landing,
		MovedTo:     movedTo,
		Lessons:     lessons,
		Modules:     apiconv.ToCourseModules(moduleRows, lessons),
		QnaItems:    qnaItems,
		Enrollment:  api.NilEnrollment{Null: true},
		Progress:    api.NilCourseProgress{Null: true},
	}
	if len(state.Lessons) > 0 {
		view.Progress = api.NewNilCourseProgress(
			apiconv.ToCourseProgress(state, apiconv.NextLessonRef(state.NextLessonSlug, lessons)),
		)
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
	view.Enrollment = api.NewNilEnrollment(apiconv.ToEnrollment(enrolled))
	return view, nil
}

// archivedPageTarget is where an archived landing page's URL moved, as legacy's
// redirect_archived_language chose it: the page the archived one names, else the
// course's main page in the request locale. Null for a page that is not
// archived. An archived page with neither has nowhere to go, which is a 404.
func (s *Server) archivedPageTarget(ctx context.Context, page *ent.CourseLandingPage) (api.NilString, error) {
	if page.State == nil || *page.State != "archived" {
		return api.NilString{Null: true}, nil
	}

	query := s.db.CourseLandingPage.Query()
	if page.LandingPageToRedirectID != nil {
		query = query.Where(courselandingpage.ID(*page.LandingPageToRedirectID))
	} else {
		query = query.Where(
			courselandingpage.CourseID(page.CourseID),
			courselandingpage.LocaleEQ(s.i18n.Locale(ctx)),
			courselandingpage.Main(true),
			courselandingpage.IDNEQ(page.ID),
		)
	}
	target, err := query.First(ctx)
	if err != nil {
		return api.NilString{}, err
	}
	return api.NewNilString(lo.FromPtr(target.Slug)), nil
}

// mainLandingPage is the course's own landing copy in the request locale: the
// page flagged main, and otherwise the oldest one, mirroring how legacy picks the
// canonical page. The locale matters beyond the copy: the course page is served
// only under its landing page's own locale (GetCourse), so a link built from a
// page of another locale would lead to a 404.
//
// `main` is nullable and most pages leave it null, so the ordering has to say
// NULLS LAST explicitly: Postgres sorts nulls first under a bare DESC, which
// would hand the canonical slot to whichever page is merely unflagged.
func (s *Server) mainLandingPage(ctx context.Context, crs *ent.Course) (api.NilCourseLandingPage, error) {
	page, err := s.db.CourseLandingPage.Query().
		Where(
			courselandingpage.CourseID(crs.ID),
			courselandingpage.LocaleEQ(s.i18n.Locale(ctx)),
		).
		Order(
			courselandingpage.ByMain(sql.OrderDesc(), sql.OrderNullsLast()),
			courselandingpage.ByID(),
		).
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

// currentLessonInfos is the translations of the lessons of the course's current
// version in the request locale, each with its lesson loaded, ordered as the
// course orders them.
//
// Course order is the Position on the lesson's version row — the same number the
// progress state's positions come from, so a lesson row and its standing line up.
// Ordering by the row's own id instead would be legacy crc32 ids in a random
// sequence, and ordering by the lesson's `natural_order` would read a column the
// loader stopped writing.
func (s *Server) currentLessonInfos(ctx context.Context, crs *ent.Course) ([]*ent.CourseLessonTranslation, error) {
	if crs.CurrentVersionID == nil {
		return nil, nil
	}

	infos, err := s.db.CourseLessonTranslation.Query().
		Where(
			courselessontranslation.LocaleEQ(s.i18n.Locale(ctx)),
			courselessontranslation.HasCourseVersionWith(courseversion.ID(*crs.CurrentVersionID)),
		).
		WithLesson().
		Order(
			courselessontranslation.ByVersionField(
				courselessonversion.FieldNaturalOrder,
				sql.OrderNullsLast(),
			),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

// currentModuleRows is the modules of the course's current version, in the
// request locale, each with its lessons in course order — the learning
// program legacy's course page rendered as an accordion.
//
// Modules are ordered by their version row's Position, as the loader built
// them; the translation id breaks ties. A module's lessons are ordered by the
// same natural order currentLessonInfos uses, so a module's lessons read in the
// order the flat list shows them.
func (s *Server) currentModuleRows(ctx context.Context, crs *ent.Course) ([]*ent.CourseModuleTranslation, error) {
	if crs.CurrentVersionID == nil {
		return nil, nil
	}
	return curriculum.Modules(ctx, s.db, s.i18n.Locale(ctx), *crs.CurrentVersionID)
}

// landingQnaItems is the landing page's questions and answers, oldest first,
// as legacy's `qna_items` association returned them. A course without landing
// copy has none.
func (s *Server) landingQnaItems(ctx context.Context, landing api.NilCourseLandingPage) ([]api.QnaItem, error) {
	page, ok := landing.Get()
	if !ok {
		return []api.QnaItem{}, nil
	}
	items, err := s.db.LandingPageQnaItem.Query().
		Where(landingpageqnaitem.CourseLandingPageID(int(page.ID))).
		Order(landingpageqnaitem.ByID()).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return s.conv.ToLandingPageQnaItems(items), nil
}
