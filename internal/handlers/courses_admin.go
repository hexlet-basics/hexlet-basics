package handlers

import (
	"context"
	"errors"

	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/localization"
)

// Admin course endpoints. The API embeds the current version, so the generated
// reads load it (ent/admincrud.go, from the AdminCRUD annotation). The cover
// image (coverAttachmentId) and the derived repositoryUrl are not written here:
// the cover is deferred until the Attachments uploader lands, and repositoryUrl
// is computed from the slug on read. There is no delete (parity with legacy).

func (s *Server) AdminListCourses(ctx context.Context, params api.AdminListCoursesParams) (api.AdminListCoursesRes, error) {
	return s.db.Course.AdminList(ctx, params.Page, params.PerPage, s.conv.ToCourses)
}

// AdminListCourseOptions serves selects: every course as id and name, with no
// page.
func (s *Server) AdminListCourseOptions(ctx context.Context) (api.AdminListCourseOptionsRes, error) {
	rows, err := s.courseDirectory.Options(ctx)
	if err != nil {
		return nil, err
	}

	options := make(api.AdminListCourseOptionsOKApplicationJSON, 0, len(rows))
	for _, row := range rows {
		options = append(options, api.CourseOption{ID: int32(row.ID), Name: row.Name})
	}
	return &options, nil
}

func (s *Server) AdminGetCourse(ctx context.Context, params api.AdminGetCourseParams) (api.AdminGetCourseRes, error) {
	return s.db.Course.AdminGet(ctx, params.ID, s.conv.ToCourse)
}

func (s *Server) AdminCreateCourse(ctx context.Context, req *api.CourseInput) (api.AdminCreateCourseRes, error) {
	return s.db.Course.AdminCreate(ctx, req, s.conv.ToCourse)
}

func (s *Server) AdminUpdateCourse(ctx context.Context, req *api.CourseInput, params api.AdminUpdateCourseParams) (api.AdminUpdateCourseRes, error) {
	return s.db.Course.AdminUpdate(ctx, params.ID, req, s.conv.ToCourse)
}

// AdminCreateCourseVersion starts a build of a new course version: it creates the
// version row in `created` state and enqueues the exercise-loader job, which
// fetches the course repo, parses it, and (on success) promotes the new version
// live. Maps 1:1 to the legacy VersionsController#create → ExerciseLoaderJob. The
// 201 body is the freshly-created (not yet built) version, so the admin UI can
// poll its state.
func (s *Server) AdminCreateCourseVersion(ctx context.Context, params api.AdminCreateCourseVersionParams) (api.AdminCreateCourseVersionRes, error) {
	version, err := s.starter.Start(ctx, int(params.ID))
	if errors.Is(err, domain.ErrNotFound) {
		return &api.NotFoundError{Message: s.i18n.Text(ctx, localization.CourseNotFound)}, nil
	}
	if err != nil {
		return nil, err
	}

	body := apiconv.CourseVersionFromEnt(version)
	return &body, nil
}

// AdminReviewCourse asks for an AI review of every lesson of the course's
// current version, all locales included.
func (s *Server) AdminReviewCourse(ctx context.Context, params api.AdminReviewCourseParams) (api.AdminReviewCourseRes, error) {
	err := s.lessonReviews.ReviewCourse(ctx, int(params.ID))
	if errors.Is(err, domain.ErrNotFound) {
		return &api.NotFoundError{Message: s.i18n.Text(ctx, localization.CourseNotFound)}, nil
	}
	if err != nil {
		return nil, err
	}
	return &api.AdminReviewCourseNoContent{}, nil
}
