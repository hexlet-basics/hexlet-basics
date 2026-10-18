package handlers

import (
	"context"

	"hexletbasics/internal/api"
)

// The public course categories (legacy `Web::LanguageCategoriesController`).
// Both reads are scoped to the request locale, so a category of another
// language is simply not found (404 centrally).

// ListPublicCourseCategories is the category index.
func (s *Server) ListPublicCourseCategories(ctx context.Context) ([]api.CourseCategory, error) {
	rows, err := s.categoryDirectory.Index(ctx, s.i18n.Locale(ctx))
	if err != nil {
		return nil, err
	}

	return s.conv.ToCourseCategories(rows), nil
}

// GetPublicCourseCategory is the category page: the category, its landing
// pages and its questions.
func (s *Server) GetPublicCourseCategory(ctx context.Context, params api.GetPublicCourseCategoryParams) (*api.CourseCategoryView, error) {
	view, err := s.categoryDirectory.View(ctx, s.i18n.Locale(ctx), params.Slug)
	if err != nil {
		return nil, err
	}

	return &api.CourseCategoryView{
		Category:     s.conv.ToCourseCategory(view.Category),
		LandingPages: s.conv.ToCatalogItems(view.LandingPages),
		QnaItems:     s.conv.ToCategoryQnaItems(view.QnaItems),
	}, nil
}
