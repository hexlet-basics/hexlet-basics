package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
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
func (s *Server) AdminListReviews(ctx context.Context, params api.AdminListReviewsParams) (*api.ReviewPage, error) {
	return s.db.Review.AdminList(ctx, params.Page, params.PerPage, params.Q, s.conv.ToReviews)
}

func (s *Server) AdminGetReview(ctx context.Context, params api.AdminGetReviewParams) (*api.Review, error) {
	return s.db.Review.AdminGet(ctx, params.ID, s.conv.ToReview)
}

func (s *Server) AdminCreateReview(ctx context.Context, req *api.ReviewInput) (*api.Review, error) {
	return s.db.Review.AdminCreate(ctx, req, s.conv.ToReview)
}

func (s *Server) AdminUpdateReview(ctx context.Context, req *api.ReviewInput, params api.AdminUpdateReviewParams) (*api.Review, error) {
	return s.db.Review.AdminUpdate(ctx, params.ID, req, s.conv.ToReview)
}

func (s *Server) AdminDeleteReview(ctx context.Context, params api.AdminDeleteReviewParams) error {
	if err := s.db.Review.AdminDelete(ctx, params.ID); err != nil {
		return err
	}

	return nil
}

// A state is not an update (ADR-0020, state transitions addendum): the
// Publication State of a review moves only through the three operations below,
// each with its own Permission, and the legal moves are the reviews module's
// (the shared transitions.Publish, Unpublish and Archive). A new review is a
// draft (the column default).

// AdminPublishReview makes a draft or archived review public.
func (s *Server) AdminPublishReview(ctx context.Context, params api.AdminPublishReviewParams) (*api.Review, error) {
	return moveState(ctx, int(params.ID), s.reviewStates.Publish, s.getAdminReview)
}

// AdminUnpublishReview takes a published review back to draft.
func (s *Server) AdminUnpublishReview(ctx context.Context, params api.AdminUnpublishReviewParams) (*api.Review, error) {
	return moveState(ctx, int(params.ID), s.reviewStates.Unpublish, s.getAdminReview)
}

// AdminArchiveReview archives a draft or published review.
func (s *Server) AdminArchiveReview(ctx context.Context, params api.AdminArchiveReviewParams) (*api.Review, error) {
	return moveState(ctx, int(params.ID), s.reviewStates.Archive, s.getAdminReview)
}

// getAdminReview is the read model a state move echoes back.
func (s *Server) getAdminReview(ctx context.Context, id int) (*api.Review, error) {
	return s.db.Review.AdminGet(ctx, apiconv.Int32FromInt(id), s.conv.ToReview)
}
