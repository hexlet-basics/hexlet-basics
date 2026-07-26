package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/coursecategory"
	"hexletbasics/ent/predicate"
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

// AdminGetCourseCategory returns a single course category by id, or 404.
func (s *Server) AdminGetCourseCategory(ctx context.Context, params api.AdminGetCourseCategoryParams) (api.AdminGetCourseCategoryRes, error) {
	row, err := s.db.CourseCategory.Get(ctx, int(params.ID))
	if ent.IsNotFound(err) {
		return &api.NotFoundError{Message: "course category not found"}, nil
	}
	if err != nil {
		return nil, err
	}
	item := s.conv.ToCourseCategory(row)
	return &item, nil
}

// AdminCreateCourseCategory creates a course category, or returns 422 with
// field errors when a uniqueness constraint (name/header/slug) is violated.
func (s *Server) AdminCreateCourseCategory(ctx context.Context, req *api.CourseCategoryInput) (api.AdminCreateCourseCategoryRes, error) {
	if errs, err := s.courseCategoryConflicts(ctx, req, 0); err != nil {
		return nil, err
	} else if len(errs) > 0 {
		return &api.ValidationError{Errors: errs}, nil
	}

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

// AdminUpdateCourseCategory updates a course category, or returns 422 on a
// uniqueness conflict (ignoring the record being updated).
func (s *Server) AdminUpdateCourseCategory(ctx context.Context, req *api.CourseCategoryInput, params api.AdminUpdateCourseCategoryParams) (api.AdminUpdateCourseCategoryRes, error) {
	if errs, err := s.courseCategoryConflicts(ctx, req, int(params.ID)); err != nil {
		return nil, err
	} else if len(errs) > 0 {
		return &api.ValidationError{Errors: errs}, nil
	}

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

// AdminDeleteCourseCategory removes a course category by id.
func (s *Server) AdminDeleteCourseCategory(ctx context.Context, params api.AdminDeleteCourseCategoryParams) error {
	return s.db.CourseCategory.DeleteOneID(int(params.ID)).Exec(ctx)
}

// courseCategoryConflicts mirrors the legacy AR uniqueness validations
// (name, header, slug) — none are DB constraints, so they are checked here.
// excludeID > 0 skips the record being updated. Slug uniqueness is scoped to
// the (currently null) locale, matching Rails `uniqueness: { scope: :locale }`.
func (s *Server) courseCategoryConflicts(ctx context.Context, req *api.CourseCategoryInput, excludeID int) (api.ValidationErrorErrors, error) {
	errs := api.ValidationErrorErrors{}

	checks := []struct {
		field string
		pred  []predicate.CourseCategory
		msg   string
	}{
		{"name", []predicate.CourseCategory{coursecategory.Name(req.Name)}, "has already been taken"},
		{"header", []predicate.CourseCategory{coursecategory.Header(req.Header)}, "has already been taken"},
		{"slug", []predicate.CourseCategory{coursecategory.Slug(req.Slug), coursecategory.LocaleIsNil()}, "has already been taken"},
	}

	for _, c := range checks {
		preds := c.pred
		if excludeID > 0 {
			preds = append(preds, coursecategory.IDNEQ(excludeID))
		}
		exists, err := s.db.CourseCategory.Query().Where(preds...).Exist(ctx)
		if err != nil {
			return nil, err
		}
		if exists {
			errs[c.field] = append(errs[c.field], c.msg)
		}
	}

	return errs, nil
}

func descriptionPtr(req *api.CourseCategoryInput) *string {
	if req.Description.Null {
		return nil
	}
	return &req.Description.Value
}
