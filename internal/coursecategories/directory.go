// Package coursecategories reads the course categories a visitor browses.
package coursecategories

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/categoryqnaitem"
	"hexletbasics/ent/coursecategory"
	"hexletbasics/ent/coursecategoryitem"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/landingpages"
)

// Directory is the public read side of the categories (legacy
// `Web::LanguageCategoriesController`). Both reads are scoped to a locale
// (legacy `with_locale`), so a category of another language is not found.
type Directory struct {
	db      *ent.Client
	catalog *landingpages.Catalog
}

// NewDirectory builds the directory over an ent client; the landing pages a
// category shows come from the same catalog rule the course list uses.
func NewDirectory(db *ent.Client) *Directory {
	return &Directory{db: db, catalog: landingpages.NewCatalog(db)}
}

// View is a category page: the category, the landing pages grouped under it and
// its questions.
type View struct {
	Category     *ent.CourseCategory
	LandingPages []*ent.CourseLandingPage
	QnaItems     []*ent.CategoryQnaItem
}

// Index lists the categories of the locale. Legacy issued no ORDER BY and got
// the table's physical order; id ascending is the deterministic stand-in.
func (d *Directory) Index(ctx context.Context, locale string) ([]*ent.CourseCategory, error) {
	rows, err := d.db.CourseCategory.Query().
		Where(coursecategory.LocaleEQ(locale)).
		Order(coursecategory.ByID()).
		All(ctx)
	return rows, domain.FromStore(err)
}

// View is the page of one category by slug. A miss is domain.ErrNotFound.
func (d *Directory) View(ctx context.Context, locale, slug string) (*View, error) {
	category, err := d.db.CourseCategory.Query().
		Where(coursecategory.LocaleEQ(locale), coursecategory.SlugEQ(slug)).
		Only(ctx)
	if err != nil {
		return nil, domain.FromStore(err)
	}

	// Legacy `category.language_landing_pages.web.where(listed: true)
	// .merge(Language.ordered)`: the catalog narrowed to the pages grouped under
	// the category through `language_category_items`.
	pages, err := d.catalog.Pages(ctx, locale,
		courselandingpage.HasCategoryItemsWith(coursecategoryitem.CourseCategoryID(category.ID)))
	if err != nil {
		return nil, domain.FromStore(err)
	}

	// Legacy read the association unordered; id ascending is oldest first.
	qnaItems, err := d.db.CategoryQnaItem.Query().
		Where(categoryqnaitem.CourseCategoryID(category.ID)).
		Order(categoryqnaitem.ByID()).
		All(ctx)
	if err != nil {
		return nil, domain.FromStore(err)
	}

	return &View{Category: category, LandingPages: pages, QnaItems: qnaItems}, nil
}
