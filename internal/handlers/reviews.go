package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/review"
	"hexletbasics/internal/api"
)

// withReviewEdges loads the course (with its current version, which ToCourse
// needs) and the user, both embedded in full by the Review API model.
func withReviewEdges(q *ent.ReviewQuery) *ent.ReviewQuery {
	return q.
		WithCourse(func(cq *ent.CourseQuery) { cq.WithCurrentVersion() }).
		WithUser()
}

// AdminListReviews returns a page of reviews, newest first.
func (s *Server) AdminListReviews(ctx context.Context, params api.AdminListReviewsParams) (api.AdminListReviewsRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.ReviewQuery { return withReviewEdges(s.db.Review.Query()).Order(ent.Desc(review.FieldID)) },
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
			return withReviewEdges(s.db.Review.Query().Where(review.ID(id))).Only(ctx)
		},
		s.conv.ToReview,
	)
}

// AdminCreateReview creates a review. courseId (language_id) and userId are
// NOT NULL in the DB; a null in the input surfaces as a constraint error via
// the central ErrorHandler (the schema marks them SetOnly, so SetInput skips
// nulls). The created row is reloaded with its edges so the response carries
// the full embedded course and user.
func (s *Server) AdminCreateReview(ctx context.Context, req *api.ReviewInput) (api.AdminCreateReviewRes, error) {
	row, err := s.db.Review.Create().SetInput(req).Save(ctx)
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
