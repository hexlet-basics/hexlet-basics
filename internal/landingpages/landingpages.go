// Package landingpages holds the public course catalog's shared scope: which
// landing pages a visitor may be shown, and the state values that decide it.
// The category page and the Yandex feed both list exactly this set, so the
// rule lives once here rather than as a copy of four predicates in each.
package landingpages

import (
	"hexletbasics/ent/course"
	"hexletbasics/ent/landingpage"
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

// Listed is legacy `Language::LandingPage.web.where(listed: true)`:
// published landing pages in locale whose Course is completed, and that are
// listed. Callers add their own narrowing (a category, `main`) and ordering.
func Listed(locale string) predicate.LandingPage {
	return landingpage.And(
		landingpage.StateEQ(StatePublished),
		landingpage.LocaleEQ(locale),
		landingpage.Listed(true),
		landingpage.HasCourseWith(course.ProgressEQ(CourseCompleted)),
	)
}
