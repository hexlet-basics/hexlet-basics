package handlers

import (
	"context"
	"errors"
	"time"

	"hexletbasics/internal/api"
	"hexletbasics/internal/domain"
)

// GetActiveBanner is the banner the pages that carry one show: the published
// banner of the request locale whose window is open now, as legacy's
// `Banner.active.first`. An open end is a null bound, so a banner without
// dates is always live. Several live ones resolve to the one that started last,
// then the newest created: NULLS LAST on the start, because Postgres sorts nulls
// first under a bare DESC and an undated banner would otherwise beat a
// scheduled one. None live is a null banner, not an error: the page just shows
// nothing.
func (s *Server) GetActiveBanner(ctx context.Context) (*api.ActiveBanner, error) {
	row, err := s.bannerDirectory.Active(ctx, s.i18n.Locale(ctx), time.Now())
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return &api.ActiveBanner{Banner: api.NilBanner{Null: true}}, nil
	case err != nil:
		return nil, err
	}
	return &api.ActiveBanner{Banner: api.NewNilBanner(s.conv.ToBanner(row))}, nil
}

// The admin operations are generated (ent/admincrud.go, from the AdminCRUD and
// AdminInput schema annotations); a missing id is ent's not-found error, which
// the central ErrorHandler maps to 404. Banners carry no uniqueness constraint,
// so there is no 409 path.

func (s *Server) AdminListBanners(ctx context.Context, params api.AdminListBannersParams) (*api.BannerPage, error) {
	return s.db.Banner.AdminList(ctx, params.Page, params.PerPage, params.Q, s.conv.ToBanners)
}

func (s *Server) AdminGetBanner(ctx context.Context, params api.AdminGetBannerParams) (*api.Banner, error) {
	return s.db.Banner.AdminGet(ctx, params.ID, s.conv.ToBanner)
}

// AdminCreateBanner creates a banner as a draft: the state is not an input (it
// moves only through the transition operations), so the column default applies.
func (s *Server) AdminCreateBanner(ctx context.Context, req *api.BannerInput) (*api.Banner, error) {
	return s.db.Banner.AdminCreate(ctx, req, s.conv.ToBanner)
}

func (s *Server) AdminUpdateBanner(ctx context.Context, req *api.BannerInput, params api.AdminUpdateBannerParams) (*api.Banner, error) {
	return s.db.Banner.AdminUpdate(ctx, params.ID, req, s.conv.ToBanner)
}

func (s *Server) AdminDeleteBanner(ctx context.Context, params api.AdminDeleteBannerParams) error {
	if err := s.db.Banner.AdminDelete(ctx, params.ID); err != nil {
		return err
	}
	return nil
}

// A state is not an update (ADR-0020, state transitions addendum): the
// Publication State of a banner moves only through the three operations below,
// each with its own Permission, and the legal moves are the banners module's
// (the shared transitions.Publish, Unpublish and Archive). A new banner is a
// draft (the column default), so it reaches the page only when someone
// publishes it.

// AdminPublishBanner makes a draft or archived banner eligible for the pages.
func (s *Server) AdminPublishBanner(ctx context.Context, params api.AdminPublishBannerParams) (*api.Banner, error) {
	return moveState(ctx, int(params.ID), s.bannerStates.Publish, s.getAdminBanner)
}

// AdminUnpublishBanner takes a published banner back to draft.
func (s *Server) AdminUnpublishBanner(ctx context.Context, params api.AdminUnpublishBannerParams) (*api.Banner, error) {
	return moveState(ctx, int(params.ID), s.bannerStates.Unpublish, s.getAdminBanner)
}

// AdminArchiveBanner archives a draft or published banner.
func (s *Server) AdminArchiveBanner(ctx context.Context, params api.AdminArchiveBannerParams) (*api.Banner, error) {
	return moveState(ctx, int(params.ID), s.bannerStates.Archive, s.getAdminBanner)
}

// getAdminBanner is the read model a state move echoes back.
func (s *Server) getAdminBanner(ctx context.Context, id int) (*api.Banner, error) {
	return s.db.Banner.AdminGet(ctx, int32(id), s.conv.ToBanner)
}
