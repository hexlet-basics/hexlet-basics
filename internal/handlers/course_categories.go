package handlers

import (
	"context"

	"hexletbasics/internal/api"
)

// The admin operations are generated (ent/admincrud.go, from the AdminCRUD and
// AdminInput schema annotations). URL stays `/admin/course_categories` for
// backward-compat; the domain concept is a course category.
//
// Uniqueness (name/header/slug-per-locale) is enforced by DB unique indexes,
// not a pre-check: a violation returns ent's constraint error, which the central
// ErrorHandler maps to 409. A missing id is ent's not-found error (404).

func (s *Server) AdminListCourseCategories(ctx context.Context, params api.AdminListCourseCategoriesParams) (api.AdminListCourseCategoriesRes, error) {
	return s.db.CourseCategory.AdminList(ctx, params.Page, params.PerPage, s.conv.ToCourseCategories)
}

func (s *Server) AdminGetCourseCategory(ctx context.Context, params api.AdminGetCourseCategoryParams) (api.AdminGetCourseCategoryRes, error) {
	return s.db.CourseCategory.AdminGet(ctx, params.ID, s.conv.ToCourseCategory)
}

func (s *Server) AdminCreateCourseCategory(ctx context.Context, req *api.CourseCategoryInput) (api.AdminCreateCourseCategoryRes, error) {
	return s.db.CourseCategory.AdminCreate(ctx, req, s.conv.ToCourseCategory)
}

func (s *Server) AdminUpdateCourseCategory(ctx context.Context, req *api.CourseCategoryInput, params api.AdminUpdateCourseCategoryParams) (api.AdminUpdateCourseCategoryRes, error) {
	return s.db.CourseCategory.AdminUpdate(ctx, params.ID, req, s.conv.ToCourseCategory)
}

func (s *Server) AdminDeleteCourseCategory(ctx context.Context, params api.AdminDeleteCourseCategoryParams) (api.AdminDeleteCourseCategoryRes, error) {
	if err := s.db.CourseCategory.AdminDelete(ctx, params.ID); err != nil {
		return nil, err
	}
	return &api.AdminDeleteCourseCategoryNoContent{}, nil
}
