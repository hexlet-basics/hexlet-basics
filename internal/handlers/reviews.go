package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/review"
	"hexletbasics/internal/api"
	"hexletbasics/internal/reviews"
	"hexletbasics/internal/transitions"
)

// ListPublicReviews is the public reviews page (legacy `Web::ReviewsController
// #index`): published reviews in the request locale, newest (highest id) first,
// pagy's 20 per page. `pinned` is deliberately not an ordering key — legacy
// ordered by id alone. Review locales are ru/en only, so an es request gets an
// empty page, not an error, as legacy's `with_locale` scope gave it.
func (s *Server) ListPublicReviews(ctx context.Context, params api.ListPublicReviewsParams) (*api.ReviewPage, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.ReviewQuery { return reviews.Public(s.db, s.i18n.Locale(ctx)) },
		s.conv.ToReviews,
		func(items []api.Review, total, page, perPage int32) *api.ReviewPage {
			return &api.ReviewPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminListReviews returns a page of reviews, newest first.
func (s *Server) AdminListReviews(ctx context.Context, params api.AdminListReviewsParams) (api.AdminListReviewsRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.ReviewQuery { return reviews.WithEdges(s.db.Review.Query()).Order(ent.Desc(review.FieldID)) },
		s.conv.ToReviews,
		func(items []api.Review, total, page, perPage int32) *api.ReviewPage {
			return &api.ReviewPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminGetReview returns a single review by id. A missing id returns ent's
// not-found error, which the central ErrorHandler maps to 404.
func (s *Server) AdminGetReview(ctx context.Context, params api.AdminGetReviewParams) (api.AdminGetReviewRes, error) {
	return s.getAdminReview(ctx, params.ID)
}

func (s *Server) getAdminReview(ctx context.Context, id int32) (*api.Review, error) {
	return getOne(ctx, int(id),
		func(ctx context.Context, id int) (*ent.Review, error) {
			return reviews.WithEdges(s.db.Review.Query().Where(review.ID(id))).Only(ctx)
		},
		s.conv.ToReview,
	)
}

// AdminCreateReview creates a review, always as a draft: the state is not part
// of the input, it moves only through the transition operations. courseId (language_id) and userId are
// NOT NULL in the DB; a null in the input surfaces as a constraint error via
// the central ErrorHandler (the schema marks them SetOnly, so SetInput skips
// nulls). The created row is reloaded with its edges so the response carries
// the full embedded course and user.
func (s *Server) AdminCreateReview(ctx context.Context, req *api.ReviewInput) (api.AdminCreateReviewRes, error) {
	row, err := s.db.Review.Create().SetInput(req).SetState(string(transitions.Draft)).Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.getAdminReview(ctx, int32(row.ID))
}

// AdminUpdateReview updates a review. A missing id returns ent's not-found error
// (mapped to 404 centrally). The generated SetInput keeps the legacy semantics:
// a null nullable field clears the column, while the course/user associations
// are only reassigned when a non-null id is supplied (SetOnly).
func (s *Server) AdminUpdateReview(ctx context.Context, req *api.ReviewInput, params api.AdminUpdateReviewParams) (api.AdminUpdateReviewRes, error) {
	row, err := s.db.Review.UpdateOneID(int(params.ID)).SetInput(req).Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.getAdminReview(ctx, int32(row.ID))
}

// AdminDeleteReview removes a review by id. A missing id returns ent's not-found
// error (mapped to 404 centrally).
func (s *Server) AdminDeleteReview(ctx context.Context, params api.AdminDeleteReviewParams) (api.AdminDeleteReviewRes, error) {
	if err := s.db.Review.DeleteOneID(int(params.ID)).Exec(ctx); err != nil {
		return nil, err
	}
	return &api.AdminDeleteReviewNoContent{}, nil
}

// A state is not an update (ADR-0018, state transitions addendum): the
// Publication State of a review moves only through the three operations below,
// each with its own Permission, and the legal moves are the shared
// transitions.Publish, Unpublish and Archive. A new review is a draft.
//
// The state column is all a transition changes: the public reviews page
// (reviews.Public) lists the published reviews of the locale, so publishing is
// what puts a review there. `pinned` is a separate flag of the review, edited
// with it, and legacy kept no cache or stamp behind the state.

// AdminPublishReview makes a draft or archived review public.
func (s *Server) AdminPublishReview(ctx context.Context, params api.AdminPublishReviewParams) (api.AdminPublishReviewRes, error) {
	return s.transitionReview(ctx, int(params.ID), transitions.Publish)
}

// AdminUnpublishReview takes a published review back to draft.
func (s *Server) AdminUnpublishReview(ctx context.Context, params api.AdminUnpublishReviewParams) (api.AdminUnpublishReviewRes, error) {
	return s.transitionReview(ctx, int(params.ID), transitions.Unpublish)
}

// AdminArchiveReview archives a draft or published review.
func (s *Server) AdminArchiveReview(ctx context.Context, params api.AdminArchiveReviewParams) (api.AdminArchiveReviewRes, error) {
	return s.transitionReview(ctx, int(params.ID), transitions.Archive)
}

// transitionReview checks the move against the table and writes the new state
// through UpdateOne, so the audit hook records the old and the new state. A
// review with no state (nullable in the legacy schema) counts as a draft. A
// missing id surfaces as ent not-found (404), before any write.
func (s *Server) transitionReview(ctx context.Context, id int, move transitions.Transition[transitions.PublicationState]) (*api.Review, error) {
	row, err := s.db.Review.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := refuseTransition(move.Check(transitions.CurrentOf(row.State))); err != nil {
		return nil, err
	}

	if err := s.db.Review.UpdateOneID(id).SetState(string(move.To)).Exec(ctx); err != nil {
		return nil, err
	}
	return s.getAdminReview(ctx, int32(id))
}
