// Package coursereview holds the rules for which course reviews a visitor may read.
package coursereview

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/predicate"
	"hexletbasics/ent/review"
	"hexletbasics/internal/domain"
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

// Directory is the read side of the reviews a visitor may see.
type Directory struct {
	db *ent.Client
}

// NewDirectory builds the directory over an ent client.
func NewDirectory(db *ent.Client) *Directory {
	return &Directory{db: db}
}

// Public returns one page of the published reviews of the locale, newest
// first, and how many there are in all.
func (d *Directory) Public(ctx context.Context, locale string, page domain.Page) ([]*ent.Review, int, error) {
	return domain.Paged(ctx, Public(d.db, locale), page)
}
