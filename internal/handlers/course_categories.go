package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/coursecategory"
	"hexletbasics/internal/api"
	"hexletbasics/internal/inputconv"
)

// AdminListCourseCategories returns a page of course categories, newest first.
//
// URL stays `/admin/language_categories` for backward-compat; the domain
// concept is a course category.
func (s *Server) AdminListCourseCategories(ctx context.Context, params api.AdminListCourseCategoriesParams) (api.AdminListCourseCategoriesRes, error) {
	return listPage(ctx, params.Page, params.PerPage,
		func() *ent.CourseCategoryQuery {
			return s.db.CourseCategory.Query().Order(ent.Desc(coursecategory.FieldID))
		},
		s.conv.ToCourseCategories,
		func(items []api.CourseCategory, total, page, perPage int32) *api.CourseCategoryPage {
			return &api.CourseCategoryPage{Items: items, Total: total, Page: page, PerPage: perPage}
		},
	)
}

// AdminGetCourseCategory returns a single course category by id. A missing id
// returns ent's not-found error, which the central ErrorHandler maps to 404.
func (s *Server) AdminGetCourseCategory(ctx context.Context, params api.AdminGetCourseCategoryParams) (api.AdminGetCourseCategoryRes, error) {
	return getOne(ctx, int(params.ID), s.db.CourseCategory.Get, s.conv.ToCourseCategory)
}

// AdminCreateCourseCategory creates a course category.
//
// Uniqueness (name/header/slug-per-locale) is enforced by DB unique indexes,
// not a pre-check: a violation returns ent's constraint error, which the central
// ErrorHandler maps to 409. This is race-free, unlike the former Exist queries.
func (s *Server) AdminCreateCourseCategory(ctx context.Context, req *api.CourseCategoryInput) (api.AdminCreateCourseCategoryRes, error) {
	row, err := s.db.CourseCategory.Create().
		SetName(req.Name).
		SetHeader(req.Header).
		SetSlug(req.Slug).
		SetNillableDescription(inputconv.Ptr(req.Description)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToCourseCategory(row)
	return &item, nil
}

// AdminUpdateCourseCategory updates a course category. A missing id (404) and a
// uniqueness violation (409) both surface as ent errors handled centrally.
func (s *Server) AdminUpdateCourseCategory(ctx context.Context, req *api.CourseCategoryInput, params api.AdminUpdateCourseCategoryParams) (api.AdminUpdateCourseCategoryRes, error) {
	upd := s.db.CourseCategory.UpdateOneID(int(params.ID)).
		SetName(req.Name).
		SetHeader(req.Header).
		SetSlug(req.Slug)
	applyNil(req.Description.Null, req.Description.Value, upd.SetDescription, upd.ClearDescription)
	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToCourseCategory(row)
	return &item, nil
}

// AdminDeleteCourseCategory removes a course category by id. A missing id
// returns ent's not-found error (mapped to 404 centrally).
func (s *Server) AdminDeleteCourseCategory(ctx context.Context, params api.AdminDeleteCourseCategoryParams) (api.AdminDeleteCourseCategoryRes, error) {
	if err := s.db.CourseCategory.DeleteOneID(int(params.ID)).Exec(ctx); err != nil {
		return nil, err
	}
	return &api.AdminDeleteCourseCategoryNoContent{}, nil
}
