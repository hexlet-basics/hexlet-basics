// Package landingpage holds the public course catalog's shared scope: which
// landing pages a visitor may be shown, and the state values that decide it.
// The category page and the Yandex feed both list exactly this set, so the
// rule lives once here rather than as a copy of four predicates in each.
package landingpage

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/ent/predicate"
)

// StatePublished is the `state` of a landing page open to visitors
// (legacy `Language::LandingPage` aasm state). The ent field is a plain string,
// so the values the queries compare against are named here once.
const StatePublished = "published"

// CourseCompleted is the `progress` of a Course whose content is
// finished (legacy `Language` enumerize progress); in-development and draft
// courses stay off the public lists.
const CourseCompleted = "completed"

// Completed is a Course whose content is finished (legacy `Language.completed`):
// in-development and draft courses stay off public lists and out of the admin's
// current lesson views.
func Completed() predicate.Course {
	return course.ReadinessEQ(CourseCompleted)
}

// Published is the part of the public scope that does not look at the Course:
// published landing pages in locale. Web narrows it; the sitemap, which lists
// pages of courses that are not completed yet, builds on it directly.
func Published(locale string) predicate.CourseLandingPage {
	return courselandingpage.And(
		courselandingpage.StateEQ(StatePublished),
		courselandingpage.LocaleEQ(locale),
	)
}

// Web is legacy `Language::LandingPage.web`: published landing pages in locale
// whose Course is completed. Every public read starts here; callers add their
// own narrowing (`listed`, `main`, a category) and ordering.
func Web(locale string) predicate.CourseLandingPage {
	return courselandingpage.And(
		Published(locale),
		courselandingpage.HasCourseWith(Completed()),
	)
}

// Listed is legacy `Language::LandingPage.web.where(listed: true)`: the Web
// scope narrowed to the pages flagged for the public lists.
func Listed(locale string) predicate.CourseLandingPage {
	return courselandingpage.And(Web(locale), courselandingpage.Listed(true))
}

// Catalog reads the public course catalog: the listed landing pages a visitor
// is shown, in the order the catalog displays them. It holds the one copy of
// that query, so the course list and a category page cannot drift apart on
// scope, ordering or eager loading.
type Catalog struct {
	db *ent.Client
}

// NewCatalog builds the catalog over an ent client.
func NewCatalog(db *ent.Client) *Catalog {
	return &Catalog{db: db}
}

// Pages returns the Listed landing pages of locale, each with its Course and
// that Course's current Version loaded. within narrows the set (a category) and
// never widens it: the catalog scope always applies.
//
// Order is the Course display order (NULLS LAST), then Course id, then landing
// page id. Legacy `Language.ordered` had no tie-breaker among a Course's several
// listed pages, leaving their order undefined; the last two terms make it
// deterministic. Ordering happens in SQL.
func (c *Catalog) Pages(ctx context.Context, locale string, within ...predicate.CourseLandingPage) ([]*ent.CourseLandingPage, error) {
	return c.db.CourseLandingPage.Query().
		Where(Listed(locale)).
		Where(within...).
		// The current Version populates Course.currentVersion in apiconv; without
		// eager-loading it the edge is nil and the field serializes as null.
		WithCourse(func(q *ent.CourseQuery) { q.WithCurrentVersion() }).
		Order(
			courselandingpage.ByCourseField(course.FieldOrder, sql.OrderNullsLast()),
			// language_id is the FK to Course, so it equals the Course id: order by
			// it directly instead of joining the edge again for its id.
			courselandingpage.ByCourseID(),
			// The edge ordering joins `languages`, so a bare `id` is ambiguous.
			func(s *sql.Selector) { s.OrderBy(s.C(courselandingpage.FieldID)) },
		).
		All(ctx)
}

// MainPages returns, for each of courseIDs that has one, its main landing page
// in locale under the Web scope: legacy `Language#main_landing_page` narrowed to
// what a visitor may open. A course with no such page is absent from the map.
// One query serves every course, whatever their number.
//
// Legacy's `has_one` picked arbitrarily among several main pages in one locale;
// the oldest is kept here, so the answer does not depend on row order. Each page
// arrives with its Course and current Version, as Pages does.
func (c *Catalog) MainPages(ctx context.Context, locale string, courseIDs ...int) (map[int]*ent.CourseLandingPage, error) {
	if len(courseIDs) == 0 {
		return map[int]*ent.CourseLandingPage{}, nil
	}

	pages, err := c.db.CourseLandingPage.Query().
		Where(
			Web(locale),
			courselandingpage.MainEQ(true),
			courselandingpage.CourseIDIn(courseIDs...),
		).
		WithCourse(func(q *ent.CourseQuery) { q.WithCurrentVersion() }).
		Order(courselandingpage.ByID()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	byCourse := make(map[int]*ent.CourseLandingPage, len(pages))
	for _, page := range pages {
		if _, seen := byCourse[page.CourseID]; !seen {
			byCourse[page.CourseID] = page
		}
	}

	return byCourse, nil
}
