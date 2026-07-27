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
func (s *Server) AdminListReviews(ctx context.Context, params api.AdminListReviewsParams) (*api.ReviewPage, error) {
	page := newPagination(params.Page, params.PerPage)

	total, err := s.db.Review.Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := withReviewEdges(s.db.Review.Query()).
		Order(ent.Desc(review.FieldID)).
		Offset(page.Offset()).
		Limit(page.Limit()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return &api.ReviewPage{
		Items:   s.conv.ToReviews(rows),
		Total:   int32(total),
		Page:    page.Page,
		PerPage: page.PerPage,
	}, nil
}

// AdminGetReview returns a single review by id. A missing id returns ent's
// not-found error, which the central ErrorHandler maps to 404.
func (s *Server) AdminGetReview(ctx context.Context, params api.AdminGetReviewParams) (*api.Review, error) {
	row, err := withReviewEdges(s.db.Review.Query().Where(review.ID(int(params.ID)))).Only(ctx)
	if err != nil {
		return nil, err
	}
	item := s.conv.ToReview(row)
	return &item, nil
}

// AdminCreateReview creates a review. courseId (language_id) and userId are
// NOT NULL in the DB; a null in the input surfaces as a constraint error via
// the central ErrorHandler. The created row is reloaded with its edges so the
// response carries the full embedded course and user.
func (s *Server) AdminCreateReview(ctx context.Context, req *api.ReviewInput) (*api.Review, error) {
	create := s.db.Review.Create().
		SetNillableBody(nilStringPtr(req.Body)).
		SetNillableFirstName(nilStringPtr(req.FirstName)).
		SetNillableLastName(nilStringPtr(req.LastName)).
		SetNillablePinned(nilBoolPtr(req.Pinned)).
		SetNillableState(nilReviewStatePtr(req.State))
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

// nilReviewStatePtr resolves ogen's NilReviewState to a *string for ent's
// SetNillableState, where nil leaves the nullable column unset (null) on create.
func nilReviewStatePtr(v api.NilReviewState) *string {
	if v.Null {
		return nil
	}
	s := string(v.Value)
	return &s
}
