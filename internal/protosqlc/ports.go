// Package protosqlc is a THROWAWAY prototype (ticket 890): the ent slice
// (landingpages.Catalog, accounts.Directory, progress enrollment upsert)
// ported to sqlc. Never merged.
package protosqlc

import (
	"context"
	"database/sql"
	"errors"

	"hexletbasics/internal/protosqlc/db"
)

type Catalog struct{ q *db.Queries }

func NewCatalog(q db.DBTX) *Catalog { return &Catalog{q: db.New(q)} }

// Pages is landingpages.Catalog.Pages: scope, eager Course+Version and the
// three-term order all live in queries/landingpages.sql.
func (c *Catalog) Pages(ctx context.Context, locale string, categoryID *int64) ([]*db.CatalogPagesRow, error) {
	arg := db.CatalogPagesParams{Locale: locale}
	if categoryID != nil {
		arg.CategoryID = sql.NullInt64{Int64: *categoryID, Valid: true}
	}
	return c.q.CatalogPages(ctx, arg)
}

// VersionsOf is the hand-written replacement for ent's WithCurrentVersion.
func (c *Catalog) VersionsOf(ctx context.Context, courses []db.Language) (map[int64]*db.LanguageVersion, error) {
	var ids []int64
	for _, course := range courses {
		if course.CurrentVersionID.Valid {
			ids = append(ids, course.CurrentVersionID.Int64)
		}
	}
	rows, err := c.q.VersionsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make(map[int64]*db.LanguageVersion, len(rows))
	for _, v := range rows {
		out[v.ID] = v
	}
	return out, nil
}

func (c *Catalog) MainPages(ctx context.Context, locale string, courseIDs ...int64) (map[int64]*db.CatalogMainPagesRow, error) {
	if len(courseIDs) == 0 {
		return map[int64]*db.CatalogMainPagesRow{}, nil
	}
	rows, err := c.q.CatalogMainPages(ctx, db.CatalogMainPagesParams{Locale: locale, CourseIds: courseIDs})
	if err != nil {
		return nil, err
	}
	out := make(map[int64]*db.CatalogMainPagesRow, len(rows))
	for _, r := range rows {
		if _, seen := out[r.LanguageLandingPage.LanguageID]; !seen {
			out[r.LanguageLandingPage.LanguageID] = r
		}
	}
	return out, nil
}

type Directory struct{ q *db.Queries }

func NewDirectory(q db.DBTX) *Directory { return &Directory{q: db.New(q)} }

// List is Directory.List: window + total (two statements, as domain.Paged does).
func (d *Directory) List(ctx context.Context, offset, limit int32) ([]*db.User, int64, error) {
	total, err := d.q.CountUsers(ctx)
	if err != nil {
		return nil, 0, err
	}
	rows, err := d.q.ListUsers(ctx, db.ListUsersParams{Lim: limit, Off: offset})
	return rows, total, err
}

func (d *Directory) Search(ctx context.Context, term string) ([]*db.User, error) {
	if term == "" {
		return nil, nil
	}
	return d.q.SearchUsers(ctx, term)
}

// Filter is the dynamic admin list: nil means "do not filter on this".
func (d *Directory) Filter(ctx context.Context, term, locale *string, admin *bool, offset, limit int32) ([]*db.User, error) {
	arg := db.FilterUsersParams{Lim: limit, Off: offset}
	if term != nil {
		arg.Term = sql.NullString{String: *term, Valid: true}
	}
	if locale != nil {
		arg.Locale = sql.NullString{String: *locale, Valid: true}
	}
	if admin != nil {
		arg.Admin = sql.NullBool{Bool: *admin, Valid: true}
	}
	return d.q.FilterUsers(ctx, arg)
}

// Enroll is the ON CONFLICT DO NOTHING insert; created=false on a lost race.
func Enroll(ctx context.Context, q db.DBTX, userID, courseID int64) (id int64, created bool, err error) {
	id, err = db.New(q).InsertEnrollmentIfAbsent(ctx, db.InsertEnrollmentIfAbsentParams{
		UserID: userID, CourseID: courseID, State: sql.NullString{String: "started", Valid: true},
	})
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	return id, err == nil, err
}
