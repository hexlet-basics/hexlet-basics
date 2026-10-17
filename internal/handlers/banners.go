package handlers

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"

	"hexletbasics/ent"
	"hexletbasics/ent/banner"
	"hexletbasics/internal/api"
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
	now := time.Now()
	row, err := s.db.Banner.Query().
		Where(
			banner.StateEQ("published"),
			banner.LocaleEQ(s.i18n.Locale(ctx)),
			banner.Or(banner.StartsAtIsNil(), banner.StartsAtLTE(now)),
			banner.Or(banner.FinishesAtIsNil(), banner.FinishesAtGTE(now)),
		).
		Order(
			banner.ByStartsAt(sql.OrderDesc(), sql.OrderNullsLast()),
			banner.ByCreatedAt(sql.OrderDesc()),
			banner.ByID(sql.OrderDesc()),
		).
		First(ctx)
	switch {
	case ent.IsNotFound(err):
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

func (s *Server) AdminListBanners(ctx context.Context, params api.AdminListBannersParams) (api.AdminListBannersRes, error) {
	return s.db.Banner.AdminList(ctx, params.Page, params.PerPage, s.conv.ToBanners)
}

func (s *Server) AdminGetBanner(ctx context.Context, params api.AdminGetBannerParams) (api.AdminGetBannerRes, error) {
	return s.db.Banner.AdminGet(ctx, params.ID, s.conv.ToBanner)
}

func (s *Server) AdminCreateBanner(ctx context.Context, req *api.BannerInput) (api.AdminCreateBannerRes, error) {
	return s.db.Banner.AdminCreate(ctx, req, s.conv.ToBanner)
}

func (s *Server) AdminUpdateBanner(ctx context.Context, req *api.BannerInput, params api.AdminUpdateBannerParams) (api.AdminUpdateBannerRes, error) {
	return s.db.Banner.AdminUpdate(ctx, params.ID, req, s.conv.ToBanner)
}

func (s *Server) AdminDeleteBanner(ctx context.Context, params api.AdminDeleteBannerParams) (api.AdminDeleteBannerRes, error) {
	if err := s.db.Banner.AdminDelete(ctx, params.ID); err != nil {
		return nil, err
	}
	return &api.AdminDeleteBannerNoContent{}, nil
}
