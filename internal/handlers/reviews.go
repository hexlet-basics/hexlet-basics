package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/domain"
)

// ListPublicReviews is the public reviews page (legacy `Web::ReviewsController
// #index`): published reviews in the request locale, newest (highest id) first,
// pagy's 20 per page. `pinned` is deliberately not an ordering key — legacy
// ordered by id alone. Review locales are ru/en only, so an es request gets an
// empty page, not an error, as legacy's `with_locale` scope gave it.
func (s *Server) ListPublicReviews(ctx context.Context, params api.ListPublicReviewsParams) (*api.ReviewPage, error) {
	locale := s.i18n.Locale(ctx)
	return listFrom(ctx, params.Page, params.PerPage,
		func(ctx context.Context, page domain.Page) ([]*ent.Review, int, error) {
			return s.reviewDirectory.Public(ctx, locale, page)
		},
		s.conv.ToReviews,
		func(items []api.Review, total, page, perPage int32) *api.ReviewPage {
			return &api.ReviewPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// The admin operations are generated (ent/admincrud.go, from the AdminCRUD and
// AdminInput schema annotations): the review is read with its course (and the
// course's current version) and its user, so a response embeds them in full.
// courseId and userId are NOT NULL in the DB; a null in the input surfaces as a
// constraint error (409 centrally), because the schema marks them SetOnly and
// SetInput skips nulls.

// AdminListReviews returns a page of reviews, newest first.
func (s *Server) AdminListReviews(ctx context.Context, params api.AdminListReviewsParams) (api.AdminListReviewsRes, error) {
	return s.db.Review.AdminList(ctx, params.Page, params.PerPage, s.conv.ToReviews)
}

func (s *Server) AdminGetReview(ctx context.Context, params api.AdminGetReviewParams) (api.AdminGetReviewRes, error) {
	return s.db.Review.AdminGet(ctx, params.ID, s.conv.ToReview)
}

func (s *Server) AdminCreateReview(ctx context.Context, req *api.ReviewInput) (api.AdminCreateReviewRes, error) {
	return s.db.Review.AdminCreate(ctx, req, s.conv.ToReview)
}

func (s *Server) AdminUpdateReview(ctx context.Context, req *api.ReviewInput, params api.AdminUpdateReviewParams) (api.AdminUpdateReviewRes, error) {
	return s.db.Review.AdminUpdate(ctx, params.ID, req, s.conv.ToReview)
}

func (s *Server) AdminDeleteReview(ctx context.Context, params api.AdminDeleteReviewParams) (api.AdminDeleteReviewRes, error) {
	if err := s.db.Review.AdminDelete(ctx, params.ID); err != nil {
		return nil, err
	}
	return &api.AdminDeleteReviewNoContent{}, nil
}
