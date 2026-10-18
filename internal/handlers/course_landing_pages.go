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

func (s *Server) AdminListCourseLandingPages(ctx context.Context, params api.AdminListCourseLandingPagesParams) (*api.CourseLandingPagePage, error) {
	return s.db.CourseLandingPage.AdminList(ctx, params.Page, params.PerPage, params.Q, s.conv.ToCourseLandingPages)
}

func (s *Server) AdminGetCourseLandingPage(ctx context.Context, params api.AdminGetCourseLandingPageParams) (*api.CourseLandingPage, error) {
	return s.db.CourseLandingPage.AdminGet(ctx, params.ID, s.conv.ToCourseLandingPage)
}

// AdminCreateCourseLandingPage creates a landing page as a draft: the state is
// not an input (it moves only through the transition operations), so the column
// default applies.
func (s *Server) AdminCreateCourseLandingPage(ctx context.Context, req *api.CourseLandingPageInput) (*api.CourseLandingPage, error) {
	return s.db.CourseLandingPage.AdminCreate(ctx, req, s.conv.ToCourseLandingPage)
}

func (s *Server) AdminUpdateCourseLandingPage(ctx context.Context, req *api.CourseLandingPageInput, params api.AdminUpdateCourseLandingPageParams) (*api.CourseLandingPage, error) {
	return s.db.CourseLandingPage.AdminUpdate(ctx, params.ID, req, s.conv.ToCourseLandingPage)
}

func (s *Server) AdminDeleteCourseLandingPage(ctx context.Context, params api.AdminDeleteCourseLandingPageParams) error {
	if err := s.db.CourseLandingPage.AdminDelete(ctx, params.ID); err != nil {
		return err
	}
	return nil
}

// A state is not an update (ADR-0020, state transitions addendum): the
// Publication State of a landing page moves only through the three operations
// below, each with its own Permission, and the legal moves are the
// landingpages module's (the shared transitions.Publish, Unpublish and
// Archive). A new landing page is a draft (the column default).

// AdminPublishCourseLandingPage makes a draft or archived landing page public.
func (s *Server) AdminPublishCourseLandingPage(ctx context.Context, params api.AdminPublishCourseLandingPageParams) (*api.CourseLandingPage, error) {
	return moveState(ctx, int(params.ID), s.landingPageStates.Publish, s.getAdminCourseLandingPage)
}

// AdminUnpublishCourseLandingPage takes a published landing page back to draft.
func (s *Server) AdminUnpublishCourseLandingPage(ctx context.Context, params api.AdminUnpublishCourseLandingPageParams) (*api.CourseLandingPage, error) {
	return moveState(ctx, int(params.ID), s.landingPageStates.Unpublish, s.getAdminCourseLandingPage)
}

// AdminArchiveCourseLandingPage archives a draft or published landing page.
func (s *Server) AdminArchiveCourseLandingPage(ctx context.Context, params api.AdminArchiveCourseLandingPageParams) (*api.CourseLandingPage, error) {
	return moveState(ctx, int(params.ID), s.landingPageStates.Archive, s.getAdminCourseLandingPage)
}

// getAdminCourseLandingPage is the read model a state move echoes back.
func (s *Server) getAdminCourseLandingPage(ctx context.Context, id int) (*api.CourseLandingPage, error) {
	return s.db.CourseLandingPage.AdminGet(ctx, int32(id), s.conv.ToCourseLandingPage)
}
