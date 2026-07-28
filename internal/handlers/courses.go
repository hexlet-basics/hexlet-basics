package handlers

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/config"
)

// Server implements the generated ogen api.Handler backed by ent.
//
// It embeds api.UnimplementedHandler so newly-added contract operations compile
// as "not implemented" until their handler lands (contract-first, ADR-0001);
// methods defined on Server override the embedded stubs.
//
// cfg supplies the public hosts used to build absolute URLs in read models
// (canonical page URLs via AppHost, self-served asset URLs via PublicURL) —
// there is no *http.Request at the ogen handler boundary to derive them from.
type Server struct {
	api.UnimplementedHandler
	db       *ent.Client
	conv     apiconv.Converter
	cfg      *config.Config
	enqueuer JobEnqueuer
}

// NewServer wires the handler to its dependencies. enqueuer inserts background
// jobs (exercise-version builds); it may be nil for read-only wirings, in which
// case a job-enqueuing operation returns an error rather than panicking.
func NewServer(db *ent.Client, cfg *config.Config, enqueuer JobEnqueuer) *Server {
	return &Server{db: db, conv: &apiconv.ConverterImpl{}, cfg: cfg, enqueuer: enqueuer}
}

// ListCourses returns the published course catalog.
//
// URL stays `/languages` for backward-compat; the domain concept is Course.
// Mirrors the legacy scope: listed landing pages joined to their Course,
// ordered by the Course display order (NULLS LAST), then Course id.
func (s *Server) ListCourses(ctx context.Context) ([]api.CourseCatalogItem, error) {
	// Order by the Course's integer `order` (NULLS LAST), then Course id. The
	// legacy query had no tie-breaker among a Course's several listed landing
	// pages, leaving that order undefined; landing page id goes last so the
	// catalog is deterministic. All ordering happens in SQL.
	pages, err := s.db.LandingPage.Query().
		Where(landingpage.Listed(true)).
		WithCourse(func(q *ent.CourseQuery) {
			// current_version populates Course.currentVersion in apiconv; without
			// eager-loading it the edge is nil and the field serializes as null.
			q.WithCurrentVersion()
		}).
		Order(
			landingpage.ByCourseField(course.FieldOrder, sql.OrderNullsLast()),
			// language_id is the FK to Course, so it equals the Course id — order
			// by it directly (a main-table column) instead of joining the edge
			// again for its id.
			landingpage.ByLanguageID(),
			// Qualify the landing page id: the edge ordering joins `languages`,
			// so a bare `id` term is ambiguous.
			func(s *sql.Selector) { s.OrderBy(s.C(landingpage.FieldID)) },
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return s.conv.ToCatalogItems(pages), nil
}
