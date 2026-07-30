package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/api"
	"hexletbasics/internal/inputconv"
)

// withLandingCourse loads the associated course, which supplies the landing
// page's courseSlug, duration and membersCount at the API boundary.
func withLandingCourse(q *ent.LandingPageQuery) *ent.LandingPageQuery {
	return q.WithCourse()
}

// URL stays `/admin/language_landing_pages`; the domain concept is a course
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
	row, err := s.db.LandingPage.Create().
		SetLanguageID(int(req.CourseId)).
		SetNillableSlug(inputconv.Ptr(req.Slug)).
		SetNillableName(inputconv.Ptr(req.Name)).
		SetNillableMain(inputconv.Ptr(req.Main)).
		SetNillableListed(inputconv.Ptr(req.Listed)).
		SetNillableFooter(inputconv.Ptr(req.Footer)).
		SetNillableFooterName(inputconv.Ptr(req.FooterName)).
		SetNillableState(inputconv.StringPtr(req.State)).
		SetNillableOrder(inputconv.Ptr(req.Order)).
		SetNillableLandingPageToRedirectID(inputconv.IntPtr(req.LandingPageToRedirectId)).
		SetNillableMetaTitle(inputconv.Ptr(req.MetaTitle)).
		SetNillableMetaDescription(inputconv.Ptr(req.MetaDescription)).
		SetNillableHeader(inputconv.Ptr(req.Header)).
		SetNillableDescription(inputconv.Ptr(req.Description)).
		SetNillableUsedInHeader(inputconv.Ptr(req.UsedInHeader)).
		SetNillableUsedInDescription(inputconv.Ptr(req.UsedInDescription)).
		SetNillableOutcomesHeader(inputconv.Ptr(req.OutcomesHeader)).
		SetNillableOutcomesDescription(inputconv.Ptr(req.OutcomesDescription)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.getAdminCourseLandingPage(ctx, int32(row.ID))
}

func (s *Server) AdminUpdateCourseLandingPage(ctx context.Context, req *api.CourseLandingPageInput, params api.AdminUpdateCourseLandingPageParams) (api.AdminUpdateCourseLandingPageRes, error) {
	upd := s.db.LandingPage.UpdateOneID(int(params.ID)).SetLanguageID(int(req.CourseId))

	applyNil(req.Slug.Null, req.Slug.Value, upd.SetSlug, upd.ClearSlug)
	applyNil(req.Name.Null, req.Name.Value, upd.SetName, upd.ClearName)
	applyNil(req.Main.Null, req.Main.Value, upd.SetMain, upd.ClearMain)
	applyNil(req.Listed.Null, req.Listed.Value, upd.SetListed, upd.ClearListed)
	applyNil(req.Footer.Null, req.Footer.Value, upd.SetFooter, upd.ClearFooter)
	applyNil(req.FooterName.Null, req.FooterName.Value, upd.SetFooterName, upd.ClearFooterName)
	applyNil(req.Order.Null, req.Order.Value, upd.SetOrder, upd.ClearOrder)
	applyNil(req.MetaTitle.Null, req.MetaTitle.Value, upd.SetMetaTitle, upd.ClearMetaTitle)
	applyNil(req.MetaDescription.Null, req.MetaDescription.Value, upd.SetMetaDescription, upd.ClearMetaDescription)
	applyNil(req.Header.Null, req.Header.Value, upd.SetHeader, upd.ClearHeader)
	applyNil(req.Description.Null, req.Description.Value, upd.SetDescription, upd.ClearDescription)
	applyNil(req.UsedInHeader.Null, req.UsedInHeader.Value, upd.SetUsedInHeader, upd.ClearUsedInHeader)
	applyNil(req.UsedInDescription.Null, req.UsedInDescription.Value, upd.SetUsedInDescription, upd.ClearUsedInDescription)
	applyNil(req.OutcomesHeader.Null, req.OutcomesHeader.Value, upd.SetOutcomesHeader, upd.ClearOutcomesHeader)
	applyNil(req.OutcomesDescription.Null, req.OutcomesDescription.Value, upd.SetOutcomesDescription, upd.ClearOutcomesDescription)
	// Typed columns whose API value needs a small conversion before the setter.
	applyNil(req.State.Null, string(req.State.Value), upd.SetState, upd.ClearState)
	applyNil(req.LandingPageToRedirectId.Null, int(req.LandingPageToRedirectId.Value), upd.SetLandingPageToRedirectID, upd.ClearLandingPageToRedirectID)

	row, err := upd.Save(ctx)
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
