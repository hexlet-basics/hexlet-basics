// Package reviews holds the rules for which course reviews a visitor may read.
package reviews

import (
	"hexletbasics/ent"
	"hexletbasics/ent/predicate"
	"hexletbasics/ent/review"
)

// StatePublished is the `state` of a review open to readers; the ent field is a
// plain string, so the value is named here, apart from the contract's enum.
const StatePublished = "published"

// Published is legacy `Review.published_state.with_locale`. Review locales are
// ru/en only, so another locale gets no rows, not an error.
func Published(locale string) predicate.Review {
	return review.And(review.StateEQ(StatePublished), review.LocaleEQ(locale))
}

// WithEdges loads the course (with its current version, which the API model
// needs) and the user, both embedded in full by the Review API model.
func WithEdges(q *ent.ReviewQuery) *ent.ReviewQuery {
	return q.
		WithCourse(func(cq *ent.CourseQuery) { cq.WithCurrentVersion() }).
		WithUser()
}

// Public is the published reviews of locale, newest (highest id) first.
// `pinned` is deliberately not an ordering key: legacy ordered by id alone.
func Public(db *ent.Client, locale string) *ent.ReviewQuery {
	return WithEdges(db.Review.Query()).
		Where(Published(locale)).
		Order(ent.Desc(review.FieldID))
}
