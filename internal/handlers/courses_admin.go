package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/localization"
)

// Admin course endpoints (legacy `/admin/languages`). The API embeds the
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

func (s *Server) AdminCreateCourse(ctx context.Context, req *api.CourseInput) (api.AdminCreateCourseRes, error) {
	row, err := s.db.Course.Create().SetInput(req).Save(ctx)
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
			Where(courselessontranslation.LanguageVersionID(*c.CurrentVersionID)).
			Order(ent.Asc(courselessontranslation.FieldID)).
			IDs(ctx)
		if err != nil {
			return nil, err
		}
	}

	if err := s.reviews.EnqueueLessonReviews(ctx, infoIDs); err != nil {
		return nil, err
	}
	return &api.AdminReviewCourseNoContent{}, nil
}
