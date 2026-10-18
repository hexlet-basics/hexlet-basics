// Package sitemap reads the rows behind the HTML sitemap page.
package sitemap

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/course"
	"hexletbasics/ent/coursecategory"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courseversion"
	"hexletbasics/internal/blog"
	"hexletbasics/internal/landingpages"
)

// Locales are the locales the sitemap lists, in page order (legacy
// `[I18n.locale, *(available_locales - [I18n.locale, :es])]` with the page
// only ever rendered for ru). es is left out because its links did not work
// on the legacy site.
var Locales = []string{"ru", "en"}

// Page is what the sitemap lists.
type Page struct {
	LandingPages []*ent.CourseLandingPage
	BlogPosts    []*ent.BlogPost
	Categories   []*ent.CourseCategory
}

// Sitemap reads the sitemap rows.
type Sitemap struct {
	db *ent.Client
}

// New builds the production sitemap reader.
func New(db *ent.Client) *Sitemap {
	return &Sitemap{db: db}
}

// Read spans locales on purpose: the page links every locale's pages from one
// place. That the page exists for ru only is the page's rule, not the
// operation's — legacy guarded the page, and the rows do not depend on the
// request locale.
func (s *Sitemap) Read(ctx context.Context) (*Page, error) {
	// Legacy `.in_order_of(:locale, ordered_locales).order(id: :asc)`: one query
	// per locale, concatenated in page order, keeps that order with typed
	// predicates, and lets each query tie the lessons check to its own locale.
	var pages []*ent.CourseLandingPage
	for _, locale := range Locales {
		rows, err := s.db.CourseLandingPage.Query().
			Where(
				landingpages.Published(locale),
				courselandingpage.Listed(true),
				// Legacy sent every lesson of the courses' current versions and the
				// page hid a landing page whose course had none in its locale. Unlike
				// the catalog, there is no completed-course filter here.
				courselandingpage.HasCourseWith(course.HasCurrentVersionWith(
					courseversion.HasLessonTranslationsWith(courselessontranslation.LocaleEQ(locale)),
				)),
			).
			Order(courselandingpage.ByID()).
			// Legacy selected only the link columns: the table is wide and the
			// list long.
			Select(courselandingpage.FieldID, courselandingpage.FieldCourseID, courselandingpage.FieldSlug,
				courselandingpage.FieldHeader, courselandingpage.FieldLocale).
			All(ctx)
		if err != nil {
			return nil, err
		}
		pages = append(pages, rows...)
	}

	// Legacy `BlogPost.published_state.order(id: :desc)`, every locale; the page
	// renders only the sitemap locales' groups.
	posts, err := s.db.BlogPost.Query().
		Where(blogpost.StateEQ(blog.StatePublished)).
		Order(ent.Desc(blogpost.FieldID)).
		// Skip the post bodies; the page shows only titles.
		Select(blogpost.FieldID, blogpost.FieldName, blogpost.FieldSlug, blogpost.FieldLocale).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Legacy `Language::Category.all` had no ORDER BY; id ascending is the
	// deterministic stand-in, as on the category index.
	categories, err := s.db.CourseCategory.Query().
		Order(coursecategory.ByID()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return &Page{LandingPages: pages, BlogPosts: posts, Categories: categories}, nil
}
