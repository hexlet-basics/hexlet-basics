package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/coursecategory"
	"hexletbasics/internal/api"
)

// AdminListCourseCategories returns a page of course categories, newest first.
//
// URL stays `/admin/language_categories` for backward-compat; the domain
// concept is a course category.
func (s *Server) AdminListCourseCategories(ctx context.Context, params api.AdminListCourseCategoriesParams) (*api.CourseCategoryPage, error) {
	page := newPagination(params.Page, params.PerPage)

	total, err := s.db.CourseCategory.Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.CourseCategory.Query().
		Order(ent.Desc(coursecategory.FieldID)).
		Offset(page.Offset()).
		Limit(page.Limit()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return &api.CourseCategoryPage{
		Items:   s.conv.ToCourseCategories(rows),
		Total:   int32(total),
		Page:    page.Page,
		PerPage: page.PerPage,
	}, nil
}

// AdminGetCourseCategory returns a single course category by id. A missing id
// returns ent's not-found error, which the central ErrorHandler maps to 404.
func (s *Server) AdminGetCourseCategory(ctx context.Context, params api.AdminGetCourseCategoryParams) (*api.CourseCategory, error) {
	row, err := s.db.CourseCategory.Get(ctx, int(params.ID))
	if err != nil {
		return nil, err
	}
	item := s.conv.ToCourseCategory(row)
	return &item, nil
}

// AdminCreateCourseCategory creates a course category.
//
// Uniqueness (name/header/slug-per-locale) is enforced by DB unique indexes,
// not a pre-check: a violation returns ent's constraint error, which the central
// ErrorHandler maps to 409. This is race-free, unlike the former Exist queries.
func (s *Server) AdminCreateCourseCategory(ctx context.Context, req *api.CourseCategoryInput) (*api.CourseCategory, error) {
	row, err := s.db.CourseCategory.Create().
		SetName(req.Name).
		SetHeader(req.Header).
		SetSlug(req.Slug).
		SetNillableDescription(descriptionPtr(req)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToCourseCategory(row)
	return &item, nil
}

// AdminUpdateCourseCategory updates a course category. A missing id (404) and a
// uniqueness violation (409) both surface as ent errors handled centrally.
func (s *Server) AdminUpdateCourseCategory(ctx context.Context, req *api.CourseCategoryInput, params api.AdminUpdateCourseCategoryParams) (*api.CourseCategory, error) {
	upd := s.db.CourseCategory.UpdateOneID(int(params.ID)).
		SetName(req.Name).
		SetHeader(req.Header).
		SetSlug(req.Slug)
	// A null description clears the column (matches the legacy assign_attributes
	// semantics); a value sets it. SetNillable would instead skip a null.
	if req.Description.Null {
		upd.ClearDescription()
	} else {
		upd.SetDescription(req.Description.Value)
	}
	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}

	item := s.conv.ToCourseCategory(row)
	return &item, nil
}

// AdminDeleteCourseCategory removes a course category by id. A missing id
// returns ent's not-found error (mapped to 404 centrally).
func (s *Server) AdminDeleteCourseCategory(ctx context.Context, params api.AdminDeleteCourseCategoryParams) error {
	return s.db.CourseCategory.DeleteOneID(int(params.ID)).Exec(ctx)
}

func descriptionPtr(req *api.CourseCategoryInput) *string {
	if req.Description.Null {
		return nil
	}
	return &req.Description.Value
}
