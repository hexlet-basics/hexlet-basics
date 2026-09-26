package handlers

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/coursemoduletranslation"
	"hexletbasics/ent/coursemoduleversion"
	"hexletbasics/ent/courseversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/landingpage"
	"hexletbasics/ent/landingpageqnaitem"
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

	modules, err := s.currentModuleList(ctx, crs)
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
		Lessons:     lessons,
		Modules:     modules,
		QnaItems:    qnaItems,
		Enrollment:  api.NilEnrollment{Null: true},
		Progress:    api.NilCourseProgress{Null: true},
	}

	// Everyone gets a position, because everyone has one: a visitor who has
	// finished nothing stands at the beginning with the first lesson open, and
	// that is a fact about the course, not about having an account. Deriving it
	// client-side instead would be a second implementation of the gate.
	learner := currentLearner(ctx)
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
//
// `main` is nullable and most pages leave it null, so the ordering has to say
// NULLS LAST explicitly: Postgres sorts nulls first under a bare DESC, which
// would hand the canonical slot to whichever page is merely unflagged.
func (s *Server) mainLandingPage(ctx context.Context, crs *ent.Course) (api.NilCourseLandingPage, error) {
	page, err := s.db.LandingPage.Query().
		Where(landingpage.CourseID(crs.ID)).
		Order(
			landingpage.ByMain(sql.OrderDesc(), sql.OrderNullsLast()),
			landingpage.ByID(),
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

// currentLessonList is the lessons of the course's current version, in the
// request locale, ordered as the course orders them.
//
// Course order is the Position on the lesson's version row — the same number the
// progress payload's positions come from, so the two arrays the pages join line
// up entry for entry. Ordering by the row's own id instead would be legacy crc32
// ids in a random sequence, and ordering by the lesson's `natural_order` would
// read a column the loader stopped writing.
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
	return s.conv.ToCourseLessonListItems(infos), nil
}

// currentModuleList is the modules of the course's current version, in the
// request locale, each with its lessons in course order — the learning
// program legacy's course page rendered as an accordion.
//
// Modules are ordered by their version row's Position, as the loader built
// them; the translation id breaks ties. A module's lessons are ordered by the
// same natural order currentLessonList uses, so a module's slugs read in the
// order the flat list shows them.
func (s *Server) currentModuleList(ctx context.Context, crs *ent.Course) ([]api.CourseModuleListItem, error) {
	if crs.CurrentVersionID == nil {
		return []api.CourseModuleListItem{}, nil
	}

	rows, err := s.db.CourseModuleTranslation.Query().
		Where(
			coursemoduletranslation.LocaleEQ(s.i18n.Locale(ctx)),
			coursemoduletranslation.CourseVersionID(*crs.CurrentVersionID),
		).
		WithVersion(func(q *ent.CourseModuleVersionQuery) {
			q.WithLessonVersions(func(q *ent.CourseLessonVersionQuery) {
				q.WithLesson().Order(
					courselessonversion.ByNaturalOrder(sql.OrderNullsLast()),
					courselessonversion.ByID(),
				)
			})
		}).
		Order(
			coursemoduletranslation.ByVersionField(coursemoduleversion.FieldOrder, sql.OrderNullsLast()),
			// Qualify the id: the edge ordering joins language_module_versions,
			// so a bare `id` is ambiguous.
			func(s *sql.Selector) { s.OrderBy(s.C(coursemoduletranslation.FieldID)) },
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return s.conv.ToCourseModuleListItems(rows), nil
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
