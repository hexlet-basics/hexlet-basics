package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/api"
)

// withLandingCourse loads the associated course, which supplies the landing
// page's courseSlug, duration and membersCount at the API boundary.
func withLandingCourse(q *ent.LandingPageQuery) *ent.LandingPageQuery {
	return q.WithCourse()
}

// The route is `/admin/course_landing_pages`; the domain concept is a course
// landing page. The `outcomesImage` asset is deferred (returned null) and the
// input's `outcomesImageAttachmentId` is ignored until the Attachments uploader
// lands — same deferral as the course cover.

func (s *Server) AdminListCourseLandingPages(ctx context.Context, params api.AdminListCourseLandingPagesParams) (api.AdminListCourseLandingPagesRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.LandingPageQuery {
			return withLandingCourse(s.db.LandingPage.Query()).Order(ent.Desc(landingpage.FieldID))
		},
		s.conv.ToCourseLandingPages,
		func(items []api.CourseLandingPage, total, page, perPage int32) *api.CourseLandingPagePage {
			return &api.CourseLandingPagePage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

func (s *Server) AdminGetCourseLandingPage(ctx context.Context, params api.AdminGetCourseLandingPageParams) (api.AdminGetCourseLandingPageRes, error) {
	return s.getAdminCourseLandingPage(ctx, params.ID)
}

func (s *Server) getAdminCourseLandingPage(ctx context.Context, id int32) (*api.CourseLandingPage, error) {
	return getOne(ctx, int(id),
		func(ctx context.Context, id int) (*ent.LandingPage, error) {
			return withLandingCourse(s.db.LandingPage.Query().Where(landingpage.ID(id))).Only(ctx)
		},
		s.conv.ToCourseLandingPage,
	)
}

func (s *Server) AdminCreateCourseLandingPage(ctx context.Context, req *api.CourseLandingPageInput) (api.AdminCreateCourseLandingPageRes, error) {
	row, err := s.db.LandingPage.Create().SetInput(req).Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.getAdminCourseLandingPage(ctx, int32(row.ID))
}

func (s *Server) AdminUpdateCourseLandingPage(ctx context.Context, req *api.CourseLandingPageInput, params api.AdminUpdateCourseLandingPageParams) (api.AdminUpdateCourseLandingPageRes, error) {
	row, err := s.db.LandingPage.UpdateOneID(int(params.ID)).SetInput(req).Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.getAdminCourseLandingPage(ctx, int32(row.ID))
}

func (s *Server) AdminDeleteCourseLandingPage(ctx context.Context, params api.AdminDeleteCourseLandingPageParams) (api.AdminDeleteCourseLandingPageRes, error) {
	if err := s.db.LandingPage.DeleteOneID(int(params.ID)).Exec(ctx); err != nil {
		return nil, err
	}
	return &api.AdminDeleteCourseLandingPageNoContent{}, nil
}
