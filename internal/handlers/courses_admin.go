package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/transitions"
)

// Admin course endpoints. The API embeds the
// current version, so reads load it. The cover image (coverAttachmentId) and
// the derived repositoryUrl are not written here — the cover is deferred until
// the Attachments uploader lands, and repositoryUrl is computed from the slug on
// read. There is no delete (parity with legacy).

func withCourseVersion(q *ent.CourseQuery) *ent.CourseQuery {
	return q.WithCurrentVersion()
}

func (s *Server) AdminListCourses(ctx context.Context, params api.AdminListCoursesParams) (api.AdminListCoursesRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.CourseQuery { return withCourseVersion(s.db.Course.Query()).Order(ent.Desc(course.FieldID)) },
		s.conv.ToCourses,
		func(items []api.Course, total, page, perPage int32) *api.CoursePage {
			return &api.CoursePage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

func (s *Server) AdminGetCourse(ctx context.Context, params api.AdminGetCourseParams) (api.AdminGetCourseRes, error) {
	return s.getAdminCourse(ctx, params.ID)
}

func (s *Server) getAdminCourse(ctx context.Context, id int32) (*api.Course, error) {
	return getOne(ctx, int(id),
		func(ctx context.Context, id int) (*ent.Course, error) {
			return withCourseVersion(s.db.Course.Query().Where(course.ID(id))).Only(ctx)
		},
		s.conv.ToCourse,
	)
}

// AdminCreateCourse creates a course, always as a draft: the Readiness is not an
// input (a state is not an update), so it is set here and moves only through the
// mark-as operations. A course left without one would be a stateless row.
func (s *Server) AdminCreateCourse(ctx context.Context, req *api.CourseInput) (api.AdminCreateCourseRes, error) {
	row, err := s.db.Course.Create().SetInput(req).SetReadiness(string(transitions.ReadinessDraft)).Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.getAdminCourse(ctx, int32(row.ID))
}

func (s *Server) AdminUpdateCourse(ctx context.Context, req *api.CourseInput, params api.AdminUpdateCourseParams) (api.AdminUpdateCourseRes, error) {
	row, err := s.db.Course.UpdateOneID(int(params.ID)).SetInput(req).Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.getAdminCourse(ctx, int32(row.ID))
}

// A state is not an update (ADR-0018, state transitions addendum): the Readiness
// of a course moves only through the three operations below, each with its own
// Permission. Any change between two distinct values is legal, so the tables in
// internal/transitions differ from the Publication State's, but the refusal of a
// repeated value is the same 409.
//
// The readiness column is all a transition changes. What it controls is read
// from it by the queries, with nothing copied or cached: the public catalogue,
// the course category pages, the sitemap and the lessons list ask for completed
// courses (landingpages.Completed), so marking a course completed lists it and
// marking it back withdraws it. The course loader writes the readiness of a
// course version (its spec's `progress`), never the course's.

// AdminMarkCourseDraft takes a course back to a draft.
func (s *Server) AdminMarkCourseDraft(ctx context.Context, params api.AdminMarkCourseDraftParams) (api.AdminMarkCourseDraftRes, error) {
	return s.transitionCourse(ctx, int(params.ID), transitions.MarkDraft)
}

// AdminMarkCourseInDevelopment marks a course as being developed.
func (s *Server) AdminMarkCourseInDevelopment(ctx context.Context, params api.AdminMarkCourseInDevelopmentParams) (api.AdminMarkCourseInDevelopmentRes, error) {
	return s.transitionCourse(ctx, int(params.ID), transitions.MarkInDevelopment)
}

// AdminMarkCourseCompleted marks a course as finished, which lists it publicly.
func (s *Server) AdminMarkCourseCompleted(ctx context.Context, params api.AdminMarkCourseCompletedParams) (api.AdminMarkCourseCompletedRes, error) {
	return s.transitionCourse(ctx, int(params.ID), transitions.MarkCompleted)
}

// transitionCourse checks the move against the table and writes the new
// readiness through UpdateOne, so the audit hook records the old and the new
// value. A missing id surfaces as ent not-found (404), before any write. The
// answer is the admin course view the update gives.
func (s *Server) transitionCourse(ctx context.Context, id int, move transitions.Transition[transitions.Readiness]) (*api.Course, error) {
	row, err := s.db.Course.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := refuseTransitionNamed(move.Check(transitions.ReadinessOf(row.Readiness)), readinessNames); err != nil {
		return nil, err
	}
	if err := s.db.Course.UpdateOneID(id).SetReadiness(string(move.To)).Exec(ctx); err != nil {
		return nil, err
	}
	return s.getAdminCourse(ctx, int32(id))
}

// AdminCreateCourseVersion starts a build of a new course version: it creates the
// version row in `created` state and enqueues the exercise-loader job, which
// fetches the course repo, parses it, and (on success) promotes the new version
// live. Maps 1:1 to the legacy VersionsController#create → ExerciseLoaderJob. The
// 201 body is the freshly-created (not yet built) version, so the admin UI can
// poll its state.
func (s *Server) AdminCreateCourseVersion(ctx context.Context, params api.AdminCreateCourseVersionParams) (api.AdminCreateCourseVersionRes, error) {
	c, err := s.db.Course.Get(ctx, int(params.ID))
	if ent.IsNotFound(err) {
		return &api.NotFoundError{Message: s.i18n.Text(ctx, localization.CourseNotFound)}, nil
	}
	if err != nil {
		return nil, err
	}

	version, err := s.starter.Start(ctx, c.ID)
	if err != nil {
		return nil, err
	}

	body := apiconv.CourseVersionFromEnt(version)
	return &body, nil
}

// AdminReviewCourse enqueues an AI review job for every lesson info of the
// course's current version, all locales included — legacy
// `language.current_lesson_infos.find_each { ReviewLessonJob.perform_later }`.
// A course without a current version simply enqueues nothing.
func (s *Server) AdminReviewCourse(ctx context.Context, params api.AdminReviewCourseParams) (api.AdminReviewCourseRes, error) {
	c, err := s.db.Course.Get(ctx, int(params.ID))
	if ent.IsNotFound(err) {
		return &api.NotFoundError{Message: s.i18n.Text(ctx, localization.CourseNotFound)}, nil
	}
	if err != nil {
		return nil, err
	}

	var infoIDs []int
	if c.CurrentVersionID != nil {
		infoIDs, err = s.db.CourseLessonTranslation.Query().
			Where(courselessontranslation.CourseVersionID(*c.CurrentVersionID)).
			Order(ent.Asc(courselessontranslation.FieldID)).
			IDs(ctx)
		if err != nil {
			return nil, err
		}
	}

	if err := s.enqueueLessonReviews(ctx, infoIDs); err != nil {
		return nil, err
	}
	return &api.AdminReviewCourseNoContent{}, nil
}
