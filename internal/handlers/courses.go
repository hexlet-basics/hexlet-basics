package handlers

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/config"
	"hexletbasics/internal/events"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/progress"
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
	db      *ent.Client
	conv    apiconv.Converter
	cfg     *config.Config
	starter VersionBuildStarter
	reviews LessonReviewEnqueuer
	// progress owns sequential progression; handlers never evaluate the gate.
	progress progress.Tracker
	auth     *AuthHandler
	i18n     *localization.Translator
	errors   *APIErrorHandler
}

// NewServer wires the handler to its dependencies.
func NewServer(
	db *ent.Client,
	cfg *config.Config,
	starter VersionBuildStarter,
	reviews LessonReviewEnqueuer,
	tracker progress.Tracker,
	registrar accounts.UserRegistrar,
	eventPublisher events.StandalonePublisher,
	translator *localization.Translator,
	errorHandler *APIErrorHandler,
) *Server {
	return &Server{
		db:       db,
		conv:     &apiconv.ConverterImpl{},
		cfg:      cfg,
		starter:  starter,
		reviews:  reviews,
		progress: tracker,
		auth:     NewAuthHandler(db, cfg, translator, errorHandler, registrar, eventPublisher, tracker),
		i18n:     translator,
		errors:   errorHandler,
	}
}

// AuthHandler returns the shared go-pkgz/auth adapter used by both the ogen
// handlers and the outer HTTP router.
func (s *Server) AuthHandler() *AuthHandler {
	return s.auth
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
			landingpage.ByCourseID(),
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
