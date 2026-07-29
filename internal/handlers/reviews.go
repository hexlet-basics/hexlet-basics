package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/review"
	"hexletbasics/internal/api"
	"hexletbasics/internal/inputconv"
)

// withReviewEdges loads the course (with its current version, which ToCourse
// needs) and the user, both embedded in full by the Review API model.
func withReviewEdges(q *ent.ReviewQuery) *ent.ReviewQuery {
	return q.
		WithCourse(func(cq *ent.CourseQuery) { cq.WithCurrentVersion() }).
		WithUser()
}

// AdminListReviews returns a page of reviews, newest first.
func (s *Server) AdminListReviews(ctx context.Context, params api.AdminListReviewsParams) (*api.ReviewPage, error) {
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
func (s *Server) AdminGetReview(ctx context.Context, params api.AdminGetReviewParams) (*api.Review, error) {
	return getOne(ctx, int(params.ID),
		func(ctx context.Context, id int) (*ent.Review, error) {
			return withReviewEdges(s.db.Review.Query().Where(review.ID(id))).Only(ctx)
		},
		s.conv.ToReview,
	)
}

// AdminCreateReview creates a review. courseId (language_id) and userId are
// NOT NULL in the DB; a null in the input surfaces as a constraint error via
// the central ErrorHandler. The created row is reloaded with its edges so the
// response carries the full embedded course and user.
func (s *Server) AdminCreateReview(ctx context.Context, req *api.ReviewInput) (*api.Review, error) {
	create := s.db.Review.Create().
		SetNillableBody(inputconv.Ptr(req.Body)).
		SetNillableFirstName(inputconv.Ptr(req.FirstName)).
		SetNillableLastName(inputconv.Ptr(req.LastName)).
		SetNillablePinned(inputconv.Ptr(req.Pinned)).
		SetNillableState(inputconv.StringPtr(req.State))
	if !req.CourseId.Null {
		create.SetLanguageID(int(req.CourseId.Value))
	}
	if !req.UserId.Null {
		create.SetUserID(int(req.UserId.Value))
	}

	row, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.AdminGetReview(ctx, api.AdminGetReviewParams{ID: int32(row.ID)})
}

// AdminUpdateReview updates a review. A missing id returns ent's not-found error
// (mapped to 404 centrally). A null nullable field clears the column, matching
// the legacy assign_attributes semantics; the course/user associations are only
// reassigned when a non-null id is supplied.
func (s *Server) AdminUpdateReview(ctx context.Context, req *api.ReviewInput, params api.AdminUpdateReviewParams) (*api.Review, error) {
	upd := s.db.Review.UpdateOneID(int(params.ID))

	if req.Body.Null {
		upd.ClearBody()
	} else {
		upd.SetBody(req.Body.Value)
	}
	if req.FirstName.Null {
		upd.ClearFirstName()
	} else {
		upd.SetFirstName(req.FirstName.Value)
	}
	if req.LastName.Null {
		upd.ClearLastName()
	} else {
		upd.SetLastName(req.LastName.Value)
	}
	if req.Pinned.Null {
		upd.ClearPinned()
	} else {
		upd.SetPinned(req.Pinned.Value)
	}
	if req.State.Null {
		upd.ClearState()
	} else {
		upd.SetState(string(req.State.Value))
	}
	if !req.CourseId.Null {
		upd.SetLanguageID(int(req.CourseId.Value))
	}
	if !req.UserId.Null {
		upd.SetUserID(int(req.UserId.Value))
	}

	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}
	return s.AdminGetReview(ctx, api.AdminGetReviewParams{ID: int32(row.ID)})
}

// AdminDeleteReview removes a review by id. A missing id returns ent's not-found
// error (mapped to 404 centrally).
func (s *Server) AdminDeleteReview(ctx context.Context, params api.AdminDeleteReviewParams) error {
	return s.db.Review.DeleteOneID(int(params.ID)).Exec(ctx)
}
