package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/banner"
	"hexletbasics/internal/api"
)

// AdminListBanners returns a page of banners, newest first.
func (s *Server) AdminListBanners(ctx context.Context, params api.AdminListBannersParams) (api.AdminListBannersRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.BannerQuery { return s.db.Banner.Query().Order(ent.Desc(banner.FieldID)) },
		s.conv.ToBanners,
		func(items []api.Banner, total, page, perPage int32) *api.BannerPage {
			return &api.BannerPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminGetBanner returns a single banner by id. A missing id returns ent's
// not-found error, which the central ErrorHandler maps to 404.
func (s *Server) AdminGetBanner(ctx context.Context, params api.AdminGetBannerParams) (api.AdminGetBannerRes, error) {
	return getOne(ctx, int(params.ID), s.db.Banner.Get, s.conv.ToBanner)
}

// AdminCreateBanner creates a banner. The input-to-builder mapping is the
// generated SetInput (adminput.tmpl), driven by the schema annotations.
//
// Banners carry no uniqueness constraint, so there is no 409 path (unlike the
// course-category resource). An empty `body` is rejected by the contract's
// minLength at decode time, surfaced as 400 by the central ErrorHandler.
func (s *Server) AdminCreateBanner(ctx context.Context, req *api.BannerInput) (api.AdminCreateBannerRes, error) {
	row, err := s.db.Banner.Create().SetInput(req).Save(ctx)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToBanner(row)
	return &item, nil
}

// AdminUpdateBanner updates a banner. A missing id returns ent's not-found error
// (mapped to 404 centrally). The generated SetInput keeps the legacy
// assign_attributes semantics: a null nullable field clears the column.
func (s *Server) AdminUpdateBanner(ctx context.Context, req *api.BannerInput, params api.AdminUpdateBannerParams) (api.AdminUpdateBannerRes, error) {
	row, err := s.db.Banner.UpdateOneID(int(params.ID)).SetInput(req).Save(ctx)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToBanner(row)
	return &item, nil
}

// AdminDeleteBanner removes a banner by id. A missing id returns ent's not-found
// error (mapped to 404 centrally).
func (s *Server) AdminDeleteBanner(ctx context.Context, params api.AdminDeleteBannerParams) (api.AdminDeleteBannerRes, error) {
	if err := s.db.Banner.DeleteOneID(int(params.ID)).Exec(ctx); err != nil {
		return nil, err
	}
	return &api.AdminDeleteBannerNoContent{}, nil
}
