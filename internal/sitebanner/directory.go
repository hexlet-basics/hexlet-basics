// Package sitebanner reads the banner a page shows.
package sitebanner

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"

	"hexletbasics/ent"
	"hexletbasics/ent/banner"
	"hexletbasics/internal/domain"
)

// Directory is the public read side of the banners.
type Directory struct {
	db *ent.Client
}

// NewDirectory builds the directory over an ent client.
func NewDirectory(db *ent.Client) *Directory {
	return &Directory{db: db}
}

// Active is the banner the pages that carry one show: the published banner of
// the locale whose window is open at now, as legacy's `Banner.active.first`. An
// open end is a null bound, so a banner without dates is always live. Several
// live ones resolve to the one that started last, then the newest created: NULLS
// LAST on the start, because Postgres sorts nulls first under a bare DESC and an
// undated banner would otherwise beat a scheduled one. None live is
// domain.ErrNotFound: the page just shows nothing.
func (d *Directory) Active(ctx context.Context, locale string, now time.Time) (*ent.Banner, error) {
	row, err := d.db.Banner.Query().
		Where(
			banner.StateEQ("published"),
			banner.LocaleEQ(locale),
			banner.Or(banner.StartsAtIsNil(), banner.StartsAtLTE(now)),
			banner.Or(banner.FinishesAtIsNil(), banner.FinishesAtGTE(now)),
		).
		Order(
			banner.ByStartsAt(sql.OrderDesc(), sql.OrderNullsLast()),
			banner.ByCreatedAt(sql.OrderDesc()),
			banner.ByID(sql.OrderDesc()),
		).
		First(ctx)

	return row, domain.FromStore(err)
}
