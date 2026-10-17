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

// AdminCreateCourse creates a course as a draft: the Readiness is not an input
// (a state is not an update), so the column default applies and it moves only
// through the mark-as operations.
func (s *Server) AdminCreateCourse(ctx context.Context, req *api.CourseInput) (api.AdminCreateCourseRes, error) {
	return s.db.Course.AdminCreate(ctx, req, s.conv.ToCourse)
}

func (s *Server) AdminUpdateCourse(ctx context.Context, req *api.CourseInput, params api.AdminUpdateCourseParams) (api.AdminUpdateCourseRes, error) {
	return s.db.Course.AdminUpdate(ctx, params.ID, req, s.conv.ToCourse)
}

// A state is not an update (ADR-0020, state transitions addendum): the Readiness
// of a course moves only through the three operations below, each with its own
// Permission, and the legal moves are the courses module's. Any change between
// two distinct values is legal, so the tables in internal/transitions differ
// from the Publication State's, but the refusal of a repeated value is the same
// 409.

// AdminMarkCourseDraft takes a course back to a draft.
func (s *Server) AdminMarkCourseDraft(ctx context.Context, params api.AdminMarkCourseDraftParams) (api.AdminMarkCourseDraftRes, error) {
	return moveReadiness(ctx, int(params.ID), s.courseReadiness.MarkDraft, s.getAdminCourse)
}

// AdminMarkCourseInDevelopment marks a course as being developed.
func (s *Server) AdminMarkCourseInDevelopment(ctx context.Context, params api.AdminMarkCourseInDevelopmentParams) (api.AdminMarkCourseInDevelopmentRes, error) {
	return moveReadiness(ctx, int(params.ID), s.courseReadiness.MarkInDevelopment, s.getAdminCourse)
}

// AdminMarkCourseCompleted marks a course as finished, which lists it publicly.
func (s *Server) AdminMarkCourseCompleted(ctx context.Context, params api.AdminMarkCourseCompletedParams) (api.AdminMarkCourseCompletedRes, error) {
	return moveReadiness(ctx, int(params.ID), s.courseReadiness.MarkCompleted, s.getAdminCourse)
}

// getAdminCourse is the admin course view a readiness move echoes back, the one
// the update gives.
func (s *Server) getAdminCourse(ctx context.Context, id int) (*api.Course, error) {
	return s.db.Course.AdminGet(ctx, int32(id), s.conv.ToCourse)
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
