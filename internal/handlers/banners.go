package handlers

import (
	"context"
	"time"

	"hexletbasics/ent"
	"hexletbasics/ent/banner"
	"hexletbasics/internal/api"
)

// AdminListBanners returns a page of banners, newest first.
func (s *Server) AdminListBanners(ctx context.Context, params api.AdminListBannersParams) (*api.BannerPage, error) {
	page := newPagination(params.Page, params.PerPage)

	total, err := s.db.Banner.Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Banner.Query().
		Order(ent.Desc(banner.FieldID)).
		Offset(page.Offset()).
		Limit(page.Limit()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return &api.BannerPage{
		Items:   s.conv.ToBanners(rows),
		Total:   int32(total),
		Page:    page.Page,
		PerPage: page.PerPage,
	}, nil
}

// AdminGetBanner returns a single banner by id. A missing id returns ent's
// not-found error, which the central ErrorHandler maps to 404.
func (s *Server) AdminGetBanner(ctx context.Context, params api.AdminGetBannerParams) (*api.Banner, error) {
	row, err := s.db.Banner.Get(ctx, int(params.ID))
	if err != nil {
		return nil, err
	}
	item := s.conv.ToBanner(row)
	return &item, nil
}

// AdminCreateBanner creates a banner.
//
// Banners carry no uniqueness constraint, so there is no 409 path (unlike the
// course-category resource). An empty `body` is rejected by the contract's
// minLength at decode time, surfaced as 400 by the central ErrorHandler.
func (s *Server) AdminCreateBanner(ctx context.Context, req *api.BannerInput) (*api.Banner, error) {
	row, err := s.db.Banner.Create().
		SetState(string(req.State)).
		SetBackground(string(req.Background)).
		SetLocale(string(req.Locale)).
		SetBody(req.Body).
		SetNillableURL(nilStringPtr(req.URL)).
		SetNillableStartsAt(nilTimePtr(req.StartsAt)).
		SetNillableFinishesAt(nilTimePtr(req.FinishesAt)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToBanner(row)
	return &item, nil
}

// AdminUpdateBanner updates a banner. A missing id returns ent's not-found error
// (mapped to 404 centrally). A null nullable field clears the column, matching
// the legacy assign_attributes semantics; a value sets it.
func (s *Server) AdminUpdateBanner(ctx context.Context, req *api.BannerInput, params api.AdminUpdateBannerParams) (*api.Banner, error) {
	upd := s.db.Banner.UpdateOneID(int(params.ID)).
		SetState(string(req.State)).
		SetBackground(string(req.Background)).
		SetLocale(string(req.Locale)).
		SetBody(req.Body)

	if req.URL.Null {
		upd.ClearURL()
	} else {
		upd.SetURL(req.URL.Value)
	}
	if req.StartsAt.Null {
		upd.ClearStartsAt()
	} else {
		upd.SetStartsAt(req.StartsAt.Value)
	}
	if req.FinishesAt.Null {
		upd.ClearFinishesAt()
	} else {
		upd.SetFinishesAt(req.FinishesAt.Value)
	}

	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToBanner(row)
	return &item, nil
}

// AdminDeleteBanner removes a banner by id. A missing id returns ent's not-found
// error (mapped to 404 centrally).
func (s *Server) AdminDeleteBanner(ctx context.Context, params api.AdminDeleteBannerParams) error {
	return s.db.Banner.DeleteOneID(int(params.ID)).Exec(ctx)
}

// nilStringPtr resolves ogen's NilString to a *string for ent's SetNillable*,
// where nil leaves the nullable column unset (null) on create.
func nilStringPtr(v api.NilString) *string {
	if v.Null {
		return nil
	}
	return &v.Value
}

// nilTimePtr resolves ogen's NilDateTime to a *time.Time for ent's SetNillable*.
func nilTimePtr(v api.NilDateTime) *time.Time {
	if v.Null {
		return nil
	}
	return &v.Value
}
