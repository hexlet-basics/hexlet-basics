package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/categoryqnaitem"
	"hexletbasics/ent/landingpageqnaitem"
	"hexletbasics/internal/api"
)

// QnA items are nested resources: every operation is scoped to a parent
// (a course category or a landing page). List returns a bare array ordered by
// id (legacy `qna_items.order(:id)`), and update/delete filter by both the
// parent and the id so an item under a different parent reads as 404 — matching
// the legacy `parent.qna_items.find(id)` scoping.

// --- Course category QnA items ---------------------------------------------

func (s *Server) AdminListCategoryQnaItems(ctx context.Context, params api.AdminListCategoryQnaItemsParams) ([]api.QnaItem, error) {
	return listAll(ctx,
		s.db.CategoryQnaItem.Query().
			Where(categoryqnaitem.LanguageCategoryID(int(params.CategoryId))).
			Order(ent.Asc(categoryqnaitem.FieldID)).
			All,
		s.conv.ToCategoryQnaItems,
	)
}

func (s *Server) AdminCreateCategoryQnaItem(ctx context.Context, req *api.QnaItemInput, params api.AdminCreateCategoryQnaItemParams) (*api.QnaItem, error) {
	row, err := s.db.CategoryQnaItem.Create().
		SetLanguageCategoryID(int(params.CategoryId)).
		SetQuestion(req.Question).
		SetAnswer(req.Answer).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	item := s.conv.ToCategoryQnaItem(row)
	return &item, nil
}

func (s *Server) AdminUpdateCategoryQnaItem(ctx context.Context, req *api.QnaItemInput, params api.AdminUpdateCategoryQnaItemParams) (*api.QnaItem, error) {
	// Scoped lookup first: a missing id, or one under a different category,
	// returns ent's not-found error (404) before any write.
	row, err := s.db.CategoryQnaItem.Query().
		Where(
			categoryqnaitem.ID(int(params.ID)),
			categoryqnaitem.LanguageCategoryID(int(params.CategoryId)),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	updated, err := row.Update().SetQuestion(req.Question).SetAnswer(req.Answer).Save(ctx)
	if err != nil {
		return nil, err
	}
	item := s.conv.ToCategoryQnaItem(updated)
	return &item, nil
}

func (s *Server) AdminDeleteCategoryQnaItem(ctx context.Context, params api.AdminDeleteCategoryQnaItemParams) error {
	row, err := s.db.CategoryQnaItem.Query().
		Where(
			categoryqnaitem.ID(int(params.ID)),
			categoryqnaitem.LanguageCategoryID(int(params.CategoryId)),
		).
		Only(ctx)
	if err != nil {
		return err
	}
	return s.db.CategoryQnaItem.DeleteOne(row).Exec(ctx)
}

// --- Course landing page QnA items -----------------------------------------

func (s *Server) AdminListLandingPageQnaItems(ctx context.Context, params api.AdminListLandingPageQnaItemsParams) ([]api.QnaItem, error) {
	return listAll(ctx,
		s.db.LandingPageQnaItem.Query().
			Where(landingpageqnaitem.LanguageLandingPageID(int(params.LandingPageId))).
			Order(ent.Asc(landingpageqnaitem.FieldID)).
			All,
		s.conv.ToLandingPageQnaItems,
	)
}

func (s *Server) AdminCreateLandingPageQnaItem(ctx context.Context, req *api.QnaItemInput, params api.AdminCreateLandingPageQnaItemParams) (*api.QnaItem, error) {
	row, err := s.db.LandingPageQnaItem.Create().
		SetLanguageLandingPageID(int(params.LandingPageId)).
		SetQuestion(req.Question).
		SetAnswer(req.Answer).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	item := s.conv.ToLandingPageQnaItem(row)
	return &item, nil
}

func (s *Server) AdminUpdateLandingPageQnaItem(ctx context.Context, req *api.QnaItemInput, params api.AdminUpdateLandingPageQnaItemParams) (*api.QnaItem, error) {
	row, err := s.db.LandingPageQnaItem.Query().
		Where(
			landingpageqnaitem.ID(int(params.ID)),
			landingpageqnaitem.LanguageLandingPageID(int(params.LandingPageId)),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	updated, err := row.Update().SetQuestion(req.Question).SetAnswer(req.Answer).Save(ctx)
	if err != nil {
		return nil, err
	}
	item := s.conv.ToLandingPageQnaItem(updated)
	return &item, nil
}

func (s *Server) AdminDeleteLandingPageQnaItem(ctx context.Context, params api.AdminDeleteLandingPageQnaItemParams) error {
	row, err := s.db.LandingPageQnaItem.Query().
		Where(
			landingpageqnaitem.ID(int(params.ID)),
			landingpageqnaitem.LanguageLandingPageID(int(params.LandingPageId)),
		).
		Only(ctx)
	if err != nil {
		return err
	}
	return s.db.LandingPageQnaItem.DeleteOne(row).Exec(ctx)
}
