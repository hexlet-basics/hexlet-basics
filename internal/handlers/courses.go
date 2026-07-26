package handlers

import (
	"context"
	"sort"

	"hexletbasics/ent"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
)

// Server implements the generated ogen api.Handler backed by ent.
//
// It embeds api.UnimplementedHandler so newly-added contract operations compile
// as "not implemented" until their handler lands (contract-first, ADR-0001);
// methods defined on Server override the embedded stubs.
type Server struct {
	api.UnimplementedHandler
	db   *ent.Client
	conv apiconv.Converter
}

func NewServer(db *ent.Client) *Server {
	return &Server{db: db, conv: &apiconv.ConverterImpl{}}
}

// ListCourses returns the published course catalog.
//
// URL stays `/languages` for backward-compat; the domain concept is Course.
// Mirrors the legacy scope: listed landing pages joined to their Course,
// ordered by the Course display order (NULLS LAST), then Course id.
func (s *Server) ListCourses(ctx context.Context) ([]api.CourseCatalogItem, error) {
	pages, err := s.db.LandingPage.Query().
		Where(landingpage.Listed(true)).
		WithCourse().
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Order by the Course's integer `order` (NULLS LAST), then Course id. The
	// legacy query had no tie-breaker among a Course's several listed landing
	// pages, leaving that order undefined; we add landing page id last so the
	// catalog is deterministic.
	sort.SliceStable(pages, func(i, j int) bool {
		ci, cj := pages[i].Edges.Course, pages[j].Edges.Course
		oi, oj := ci.Order, cj.Order
		switch {
		case oi != nil && oj != nil && *oi != *oj:
			return *oi < *oj
		case oi != nil && oj == nil:
			return true
		case oi == nil && oj != nil:
			return false
		}
		if ci.ID != cj.ID {
			return ci.ID < cj.ID
		}
		return pages[i].ID < pages[j].ID
	})

	return s.conv.ToCatalogItems(pages), nil
}
