package handlers

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"hexletbasics/ent"
	"hexletbasics/ent/categoryqnaitem"
	"hexletbasics/ent/course"
	"hexletbasics/ent/coursecategory"
	"hexletbasics/ent/coursecategoryitem"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/api"
	"hexletbasics/internal/landingpages"
)

// The public course categories (legacy `Web::LanguageCategoriesController`).
// Both reads are scoped to the request locale (legacy `with_locale`), so a
// category of another language is simply not found.

// ListPublicCourseCategories is the category index. Legacy issued no ORDER BY
// and got the table's physical order; id ascending is the deterministic
// stand-in.
func (s *Server) ListPublicCourseCategories(ctx context.Context) ([]api.CourseCategory, error) {
	rows, err := s.db.CourseCategory.Query().
		Where(coursecategory.LocaleEQ(s.i18n.Locale(ctx))).
		Order(coursecategory.ByID()).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return s.conv.ToCourseCategories(rows), nil
}

// GetPublicCourseCategory is the category page: the category, its landing
// pages and its questions. A miss is ent not-found, which the central handler
// maps to 404.
func (s *Server) GetPublicCourseCategory(ctx context.Context, params api.GetPublicCourseCategoryParams) (api.GetPublicCourseCategoryRes, error) {
	locale := s.i18n.Locale(ctx)

	category, err := s.db.CourseCategory.Query().
		Where(coursecategory.LocaleEQ(locale), coursecategory.SlugEQ(params.Slug)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Legacy `category.language_landing_pages.web.where(listed: true)
	// .merge(Language.ordered)`: pages grouped under the category through
	// `language_category_items`, published, in the request locale, listed, of a
	// completed course, by course order. Legacy had no tie-breaker; course id
	// then landing page id make the order deterministic, as in ListCourses.
	pages, err := s.db.LandingPage.Query().
		Where(
			landingpage.HasCategoryItemsWith(coursecategoryitem.CourseCategoryID(category.ID)),
			landingpages.Listed(locale),
		).
		WithCourse(func(q *ent.CourseQuery) { q.WithCurrentVersion() }).
		Order(
			landingpage.ByCourseField(course.FieldOrder, sql.OrderNullsLast()),
			landingpage.ByCourseID(),
			// The edge ordering joins `languages`, so a bare `id` is ambiguous.
			func(s *sql.Selector) { s.OrderBy(s.C(landingpage.FieldID)) },
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Legacy read the association unordered; id ascending is oldest first.
	qnaItems, err := s.db.CategoryQnaItem.Query().
		Where(categoryqnaitem.CourseCategoryID(category.ID)).
		Order(categoryqnaitem.ByID()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return &api.CourseCategoryView{
		Category:     s.conv.ToCourseCategory(category),
		LandingPages: s.conv.ToCatalogItems(pages),
		QnaItems:     s.conv.ToCategoryQnaItems(qnaItems),
	}, nil
}
