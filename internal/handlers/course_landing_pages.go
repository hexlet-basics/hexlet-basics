package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/internal/api"
	"hexletbasics/internal/transitions"
)

// withLandingCourse loads the associated course, which supplies the landing
// page's courseSlug, duration and membersCount at the API boundary.
func withLandingCourse(q *ent.CourseLandingPageQuery) *ent.CourseLandingPageQuery {
	return q.WithCourse()
}

// The route is `/admin/course_landing_pages`; the domain concept is a course
// landing page. The `outcomesImage` asset is deferred (returned null) and the
// input's `outcomesImageAttachmentId` is ignored until the Attachments uploader
// lands — same deferral as the course cover.

func (s *Server) AdminListCourseLandingPages(ctx context.Context, params api.AdminListCourseLandingPagesParams) (api.AdminListCourseLandingPagesRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.CourseLandingPageQuery {
			return withLandingCourse(s.db.CourseLandingPage.Query()).Order(ent.Desc(courselandingpage.FieldID))
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
		func(ctx context.Context, id int) (*ent.CourseLandingPage, error) {
			return withLandingCourse(s.db.CourseLandingPage.Query().Where(courselandingpage.ID(id))).Only(ctx)
		},
		s.conv.ToCourseLandingPage,
	)
}

// AdminCreateCourseLandingPage creates a landing page, always as a draft: the
// state is not an input, it moves only through the transition operations.
func (s *Server) AdminCreateCourseLandingPage(ctx context.Context, req *api.CourseLandingPageInput) (api.AdminCreateCourseLandingPageRes, error) {
	row, err := s.db.CourseLandingPage.Create().SetInput(req).SetState(string(transitions.Draft)).Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.getAdminCourseLandingPage(ctx, int32(row.ID))
}

func (s *Server) AdminUpdateCourseLandingPage(ctx context.Context, req *api.CourseLandingPageInput, params api.AdminUpdateCourseLandingPageParams) (api.AdminUpdateCourseLandingPageRes, error) {
	row, err := s.db.CourseLandingPage.UpdateOneID(int(params.ID)).SetInput(req).Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.getAdminCourseLandingPage(ctx, int32(row.ID))
}

func (s *Server) AdminDeleteCourseLandingPage(ctx context.Context, params api.AdminDeleteCourseLandingPageParams) (api.AdminDeleteCourseLandingPageRes, error) {
	if err := s.db.CourseLandingPage.DeleteOneID(int(params.ID)).Exec(ctx); err != nil {
		return nil, err
	}
	return &api.AdminDeleteCourseLandingPageNoContent{}, nil
}

// A state is not an update (ADR-0017, state transitions addendum): the
// Publication State of a landing page moves only through the three operations
// below, each with its own Permission, and the legal moves are the shared
// transitions.Publish, Unpublish and Archive. A new landing page is a draft.
//
// The state column is all a transition changes, and every public read follows
// from it through landingpages.Published (the catalog, the category pages, the
// Yandex feed and the sitemap), so publishing a page also lists it under its
// Course Categories through the category items it already has: nothing is
// copied or cached. Legacy additionally validated a description on a published
// page; the Go update never enforced that or any other presence rule of the
// page, so it is not carried over by the transition either.

// AdminPublishCourseLandingPage makes a draft or archived landing page public.
func (s *Server) AdminPublishCourseLandingPage(ctx context.Context, params api.AdminPublishCourseLandingPageParams) (api.AdminPublishCourseLandingPageRes, error) {
	return s.transitionCourseLandingPage(ctx, int(params.ID), transitions.Publish)
}

// AdminUnpublishCourseLandingPage takes a published landing page back to draft.
func (s *Server) AdminUnpublishCourseLandingPage(ctx context.Context, params api.AdminUnpublishCourseLandingPageParams) (api.AdminUnpublishCourseLandingPageRes, error) {
	return s.transitionCourseLandingPage(ctx, int(params.ID), transitions.Unpublish)
}

// AdminArchiveCourseLandingPage archives a draft or published landing page.
func (s *Server) AdminArchiveCourseLandingPage(ctx context.Context, params api.AdminArchiveCourseLandingPageParams) (api.AdminArchiveCourseLandingPageRes, error) {
	return s.transitionCourseLandingPage(ctx, int(params.ID), transitions.Archive)
}

// transitionCourseLandingPage checks the move against the table and writes the
// new state through UpdateOne, so the audit hook records the old and the new
// state. A page with no state (nullable in the legacy schema) counts as a
// draft. A missing id surfaces as ent not-found (404), before any write.
func (s *Server) transitionCourseLandingPage(ctx context.Context, id int, move transitions.Transition[transitions.PublicationState]) (*api.CourseLandingPage, error) {
	row, err := s.db.CourseLandingPage.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := refuseTransition(move.Check(transitions.CurrentOf(row.State))); err != nil {
		return nil, err
	}

	if err := s.db.CourseLandingPage.UpdateOneID(id).SetState(string(move.To)).Exec(ctx); err != nil {
		return nil, err
	}
	return s.getAdminCourseLandingPage(ctx, int32(id))
}
