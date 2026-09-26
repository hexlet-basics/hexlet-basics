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
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/books"
	"hexletbasics/internal/config"
	"hexletbasics/internal/events"
	"hexletbasics/internal/feeds"
	"hexletbasics/internal/leads"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/progress"
)

// Server implements the generated ogen api.Handler backed by ent.
//
// It embeds api.UnimplementedHandler so newly-added contract operations compile
// as "not implemented" until their handler lands (contract-first, ADR-0001);
// methods defined on Server override the embedded stubs.
//
// cfg supplies the public origins used to build absolute URLs in read models
// (self-served asset URLs via PublicURL, emailed and redirect links via
// SiteURL) — there is no *http.Request at the ogen handler boundary to derive
// them from. Canonical page links are the frontend's: it knows the page's own
// path.
type Server struct {
	api.UnimplementedHandler
	db      *ent.Client
	conv    apiconv.Converter
	cfg     *config.Config
	starter VersionBuildStarter
	reviews LessonReviewEnqueuer
	// relatedCourses schedules the AI related-courses pick for blog posts.
	relatedCourses RelatedCoursesSuggestionEnqueuer
	// relatedCoursesSet replaces a post's related courses in one transaction.
	relatedCoursesSet RelatedCoursesReplacer
	// progress owns sequential progression; handlers never evaluate the gate.
	progress progress.Tracker
	// assets owns upload policy (MIME allowlist, size cap) and persistence;
	// the upload operation only translates its outcome to the contract.
	assets *assetstore.Store
	auth   *AuthHandler
	// leads stores a lead and raises LeadCreated in one transaction.
	leads leads.Creator
	// books stores a book request and raises BookRequested in one transaction.
	books  books.Requester
	i18n   *localization.Translator
	errors *APIErrorHandler
	// yandexFeed builds the Yandex course catalogue behind the feed routes.
	yandexFeed *feeds.Yandex
	// assistant is the in-lesson chat: quota, prompt, stream and storage.
	assistant *assistant.Chat
}

// Deps are the handler's collaborators. A struct rather than positional
// parameters: most are interfaces the tests fill with recording adapters, and
// a test that never reaches a collaborator (no lead is submitted through the
// auth router) simply leaves its field zero instead of passing a labelled nil.
// Plain fields, no DI tags, so the handlers package stays injector-agnostic.
type Deps struct {
	DB                *ent.Client
	Config            *config.Config
	Starter           VersionBuildStarter
	Reviews           LessonReviewEnqueuer
	RelatedCourses    RelatedCoursesSuggestionEnqueuer
	RelatedCoursesSet RelatedCoursesReplacer
	Emails            AccountEmailEnqueuer
	Progress          progress.Tracker
	Assets            *assetstore.Store
	Registrar         accounts.UserRegistrar
	Remover           accounts.AccountRemover
	Events            events.StandalonePublisher
	Leads             leads.Creator
	Books             books.Requester
	I18n              *localization.Translator
	Errors            *APIErrorHandler
	YandexFeed        *feeds.Yandex
	Assistant         *assistant.Chat
}

// NewServer wires the handler to its dependencies.
func NewServer(deps Deps) *Server {
	return &Server{
		db:                deps.DB,
		conv:              &apiconv.ConverterImpl{},
		cfg:               deps.Config,
		starter:           deps.Starter,
		reviews:           deps.Reviews,
		relatedCourses:    deps.RelatedCourses,
		relatedCoursesSet: deps.RelatedCoursesSet,
		progress:          deps.Progress,
		assets:            deps.Assets,
		auth: NewAuthHandler(deps.DB, deps.Config, deps.I18n, deps.Errors, deps.Registrar,
			deps.Remover, deps.Events, deps.Progress, deps.Emails),
		i18n:       deps.I18n,
		leads:      deps.Leads,
		books:      deps.Books,
		errors:     deps.Errors,
		yandexFeed: deps.YandexFeed,
		assistant:  deps.Assistant,
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
