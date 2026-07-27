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

// URL stays `/admin/language_landing_pages`; the domain concept is a course
// landing page. The `outcomesImage` asset is deferred (returned null) and the
// input's `outcomesImageAttachmentId` is ignored until the Attachments uploader
// lands — same deferral as the course cover.

func (s *Server) AdminListCourseLandingPages(ctx context.Context, params api.AdminListCourseLandingPagesParams) (*api.CourseLandingPagePage, error) {
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

func (s *Server) AdminGetCourseLandingPage(ctx context.Context, params api.AdminGetCourseLandingPageParams) (*api.CourseLandingPage, error) {
	return getOne(ctx, int(params.ID),
		func(ctx context.Context, id int) (*ent.LandingPage, error) {
			return withLandingCourse(s.db.LandingPage.Query().Where(landingpage.ID(id))).Only(ctx)
		},
		s.conv.ToCourseLandingPage,
	)
}

func (s *Server) AdminCreateCourseLandingPage(ctx context.Context, req *api.CourseLandingPageInput) (*api.CourseLandingPage, error) {
	row, err := s.db.LandingPage.Create().
		SetLanguageID(int(req.CourseId)).
		SetNillableSlug(nilStringPtr(req.Slug)).
		SetNillableName(nilStringPtr(req.Name)).
		SetNillableMain(nilBoolPtr(req.Main)).
		SetNillableListed(nilBoolPtr(req.Listed)).
		SetNillableFooter(nilBoolPtr(req.Footer)).
		SetNillableFooterName(nilStringPtr(req.FooterName)).
		SetNillableState(nilLandingPageStatePtr(req.State)).
		SetNillableOrder(nilStringPtr(req.Order)).
		SetNillableLandingPageToRedirectID(nilIntPtr(req.LandingPageToRedirectId)).
		SetNillableMetaTitle(nilStringPtr(req.MetaTitle)).
		SetNillableMetaDescription(nilStringPtr(req.MetaDescription)).
		SetNillableHeader(nilStringPtr(req.Header)).
		SetNillableDescription(nilStringPtr(req.Description)).
		SetNillableUsedInHeader(nilStringPtr(req.UsedInHeader)).
		SetNillableUsedInDescription(nilStringPtr(req.UsedInDescription)).
		SetNillableOutcomesHeader(nilStringPtr(req.OutcomesHeader)).
		SetNillableOutcomesDescription(nilStringPtr(req.OutcomesDescription)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.AdminGetCourseLandingPage(ctx, api.AdminGetCourseLandingPageParams{ID: int32(row.ID)})
}

func (s *Server) AdminUpdateCourseLandingPage(ctx context.Context, req *api.CourseLandingPageInput, params api.AdminUpdateCourseLandingPageParams) (*api.CourseLandingPage, error) {
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
	return s.AdminGetCourseLandingPage(ctx, api.AdminGetCourseLandingPageParams{ID: int32(row.ID)})
}

func (s *Server) AdminDeleteCourseLandingPage(ctx context.Context, params api.AdminDeleteCourseLandingPageParams) error {
	return s.db.LandingPage.DeleteOneID(int(params.ID)).Exec(ctx)
}

// nilIntPtr resolves ogen's NilInt32 to a *int for ent's SetNillable* on int
// columns, where nil leaves the nullable column unset (null) on create.
func nilIntPtr(v api.NilInt32) *int {
	if v.Null {
		return nil
	}
	n := int(v.Value)
	return &n
}

// nilLandingPageStatePtr resolves ogen's NilLandingPageState to a *string for
// ent's SetNillableState.
func nilLandingPageStatePtr(v api.NilLandingPageState) *string {
	if v.Null {
		return nil
	}
	s := string(v.Value)
	return &s
}
