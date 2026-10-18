package handlers

import (
	"context"

	"hexletbasics/internal/api"
	"hexletbasics/internal/qnaitem"
)

// QnA items are nested resources: every operation is scoped to a parent
// (a course category or a landing page). List returns a bare array ordered by
// id, and update/delete are scoped to the parent as well as the id, so an item
// under a different parent reads as 404. The qnaitems module owns both rules.

func qnaFieldsOf(in *api.QnaItemInput) qnaitem.Fields {
	return qnaitem.Fields{Question: in.Question, Answer: in.Answer}
}

// --- Course category QnA items ---------------------------------------------

func (s *Server) AdminListCategoryQnaItems(ctx context.Context, params api.AdminListCategoryQnaItemsParams) ([]api.QnaItem, error) {
	rows, err := s.categoryQna.List(ctx, int(params.CategoryId))
	if err != nil {
		return nil, err
	}
	return s.conv.ToCategoryQnaItems(rows), nil
}

func (s *Server) AdminCreateCategoryQnaItem(ctx context.Context, req *api.QnaItemInput, params api.AdminCreateCategoryQnaItemParams) (*api.QnaItem, error) {
	row, err := s.categoryQna.Create(ctx, int(params.CategoryId), qnaFieldsOf(req))
	if err != nil {
		return nil, err
	}
	item := s.conv.ToCategoryQnaItem(row)
	return &item, nil
}

func (s *Server) AdminUpdateCategoryQnaItem(ctx context.Context, req *api.QnaItemInput, params api.AdminUpdateCategoryQnaItemParams) (*api.QnaItem, error) {
	row, err := s.categoryQna.Update(ctx, int(params.CategoryId), int(params.ID), qnaFieldsOf(req))
	if err != nil {
		return nil, err
	}
	item := s.conv.ToCategoryQnaItem(row)
	return &item, nil
}

func (s *Server) AdminDeleteCategoryQnaItem(ctx context.Context, params api.AdminDeleteCategoryQnaItemParams) error {
	if err := s.categoryQna.Delete(ctx, int(params.CategoryId), int(params.ID)); err != nil {
		return err
	}
	return nil
}

// --- Course landing page QnA items -----------------------------------------

func (s *Server) AdminListLandingPageQnaItems(ctx context.Context, params api.AdminListLandingPageQnaItemsParams) ([]api.QnaItem, error) {
	rows, err := s.landingPageQna.List(ctx, int(params.LandingPageId))
	if err != nil {
		return nil, err
	}
	return s.conv.ToLandingPageQnaItems(rows), nil
}

func (s *Server) AdminCreateLandingPageQnaItem(ctx context.Context, req *api.QnaItemInput, params api.AdminCreateLandingPageQnaItemParams) (*api.QnaItem, error) {
	row, err := s.landingPageQna.Create(ctx, int(params.LandingPageId), qnaFieldsOf(req))
	if err != nil {
		return nil, err
	}
	item := s.conv.ToLandingPageQnaItem(row)
	return &item, nil
}

func (s *Server) AdminUpdateLandingPageQnaItem(ctx context.Context, req *api.QnaItemInput, params api.AdminUpdateLandingPageQnaItemParams) (*api.QnaItem, error) {
	row, err := s.landingPageQna.Update(ctx, int(params.LandingPageId), int(params.ID), qnaFieldsOf(req))
	if err != nil {
		return nil, err
	}
	item := s.conv.ToLandingPageQnaItem(row)
	return &item, nil
}

func (s *Server) AdminDeleteLandingPageQnaItem(ctx context.Context, params api.AdminDeleteLandingPageQnaItemParams) error {
	if err := s.landingPageQna.Delete(ctx, int(params.LandingPageId), int(params.ID)); err != nil {
		return err
	}
	return nil
}
