package handlers

import (
	"context"

	"hexletbasics/internal/api"
)

// The route is `/admin/course_landing_pages`; the domain concept is a course
// landing page. The admin operations are generated (ent/admincrud.go); the
// schema's AdminCRUD annotation loads the course that supplies the page's
// courseSlug, duration and membersCount. The `outcomesImage` asset is deferred
// (returned null) and the input's `outcomesImageAttachmentId` is ignored until
// the Attachments uploader lands — same deferral as the course cover.

func (s *Server) AdminListCourseLandingPages(ctx context.Context, params api.AdminListCourseLandingPagesParams) (api.AdminListCourseLandingPagesRes, error) {
	return s.db.CourseLandingPage.AdminList(ctx, params.Page, params.PerPage, s.conv.ToCourseLandingPages)
}

func (s *Server) AdminGetCourseLandingPage(ctx context.Context, params api.AdminGetCourseLandingPageParams) (api.AdminGetCourseLandingPageRes, error) {
	return s.db.CourseLandingPage.AdminGet(ctx, params.ID, s.conv.ToCourseLandingPage)
}

func (s *Server) AdminCreateCourseLandingPage(ctx context.Context, req *api.CourseLandingPageInput) (api.AdminCreateCourseLandingPageRes, error) {
	return s.db.CourseLandingPage.AdminCreate(ctx, req, s.conv.ToCourseLandingPage)
}

func (s *Server) AdminUpdateCourseLandingPage(ctx context.Context, req *api.CourseLandingPageInput, params api.AdminUpdateCourseLandingPageParams) (api.AdminUpdateCourseLandingPageRes, error) {
	return s.db.CourseLandingPage.AdminUpdate(ctx, params.ID, req, s.conv.ToCourseLandingPage)
}

func (s *Server) AdminDeleteCourseLandingPage(ctx context.Context, params api.AdminDeleteCourseLandingPageParams) (api.AdminDeleteCourseLandingPageRes, error) {
	if err := s.db.CourseLandingPage.AdminDelete(ctx, params.ID); err != nil {
		return nil, err
	}
	return &api.AdminDeleteCourseLandingPageNoContent{}, nil
}
