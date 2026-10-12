package handlers

import (
	"context"


	"hexletbasics/ent"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/books"
	"hexletbasics/internal/config"
	"hexletbasics/internal/events"
	"hexletbasics/internal/feeds"
	"hexletbasics/internal/landingpages"
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
	// catalog reads the listed landing pages the course list and category pages share.
	catalog *landingpages.Catalog
	// assistant is the in-lesson chat: quota, prompt, stream and storage.
	assistant *assistant.Chat
}

// Deps are the handler's collaborators. A struct rather than positional
// parameters: most are interfaces the tests fill with moq mocks, and
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
		catalog:    landingpages.NewCatalog(deps.DB),
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
// Mirrors the legacy scope (`LandingPage.web.where(listed: true)`): published,
// listed landing pages of the request locale whose Course is completed, in
// catalog order (see landingpages.Catalog).
func (s *Server) ListCourses(ctx context.Context) ([]api.CourseCatalogItem, error) {
	pages, err := s.catalog.Pages(ctx, s.i18n.Locale(ctx))
	if err != nil {
		return nil, err
	}

	return s.conv.ToCatalogItems(pages), nil
}
