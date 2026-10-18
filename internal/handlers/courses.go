// Package handlers implements the generated api.Handler: it translates the
// contract to and from the use-case modules and owns the HTTP-only concerns
// (security, error statuses, middleware).
package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/account"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/auth"
	"hexletbasics/internal/blog"
	"hexletbasics/internal/book"
	"hexletbasics/internal/catalog"
	"hexletbasics/internal/coursepage"
	"hexletbasics/internal/coursereview"
	"hexletbasics/internal/curriculum"
	"hexletbasics/internal/dashboard"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/events"
	"hexletbasics/internal/feed"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/landingpage"
	"hexletbasics/internal/leadcapture"
	"hexletbasics/internal/lessonreview"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/permission"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/qnaitem"
	"hexletbasics/internal/relatedcourse"
	"hexletbasics/internal/search"
	"hexletbasics/internal/selfprotection"
	"hexletbasics/internal/sitebanner"
	"hexletbasics/internal/sitemap"
	"hexletbasics/internal/versionbuild"
)

// Server implements the generated ogen api.Handler backed by ent.
//
// It embeds api.UnimplementedHandler so newly-added contract operations compile
// as "not implemented" until their handler lands (contract-first, ADR-0001);
// methods defined on Server override the embedded stubs.
//
// settings supplies the public origins used to build absolute URLs in read models
// (self-served asset URLs via PublicURL, emailed and redirect links via
// SiteURL) — there is no *http.Request at the ogen handler boundary to derive
// them from. Canonical page links are the frontend's: it knows the page's own
// path.
type Server struct {
	api.UnimplementedHandler

	db       *ent.Client
	conv     apiconv.Converter
	settings Settings
	starter  *versionbuild.Starter
	// jobs schedules the background work handlers trigger: lesson reviews, the
	// related-courses pick, account emails.
	jobs *jobs.Enqueuer
	// relatedCoursesSet replaces a post's related courses in one transaction.
	relatedCoursesSet *relatedcourse.Replacer
	// progress owns what a learner has finished; handlers never derive it.
	progress *progress.Progress
	// assets owns upload policy (MIME allowlist, size cap) and persistence;
	// the upload operation only translates its outcome to the contract.
	assets *assetstore.Store
	// sessions begins and ends the signed-in visitor's session (ADR-0017);
	// security is the seam that recognises a request by it.
	sessions *auth.Authenticator
	security *Security
	// credentials checks and changes passwords; registrar creates accounts and
	// remover removes them.
	credentials *account.Credentials
	registrar   *account.Registrar
	remover     *account.Remover
	// userEditor and userDirectory are the administrators' write and read side of
	// user records.
	userEditor    *account.Editor
	userDirectory *account.Directory
	events        events.StandalonePublisher
	// emailTokens verifies the Magic Link and Password Reset links.
	emailTokens *emailtoken.Tokens
	// guests is the visitor-progress cookie.
	guests *guestCookies
	// leads stores a lead and raises LeadCreated in one transaction.
	leads *leadcapture.Recorder
	// books stores a book request and raises BookRequested in one transaction.
	books  *book.Recorder
	i18n   *localization.Translator
	errors *APIErrorHandler
	// yandexFeed builds the Yandex course catalogue behind the feed routes.
	yandexFeed *feed.Yandex
	// catalog reads the listed landing pages the course list and category pages share.
	catalog *landingpage.Catalog
	// coursePages reads the course page and the lesson page.
	coursePages *coursepage.Pages
	// assistant is the in-lesson chat: quota, prompt, stream and storage.
	assistant *assistant.Chat
	// blogPublic reads the posts a visitor may open.
	blogPublic *blog.Public
	// blogReader assembles posts (cover, likes, courses, reading time) for the
	// admin and the public lists alike; blogLikes records a like.
	blogReader *blog.Reader
	blogLikes  *blog.Liker
	// postRemover deletes a post with its dependents in one transaction.
	postRemover *blog.Remover
	// permissions writes a role's grants.
	permissions       *permission.Permissions
	selfProtection    *selfprotection.Guard
	bannerStates      *sitebanner.Editor
	reviewStates      *coursereview.Editor
	landingPageStates *landingpage.Editor
	courseReadiness   *catalog.Readiness
	// dashboard assembles the learner's courses with cards and progress.
	dashboard *dashboard.Dashboard
	// sitemap reads the rows behind the sitemap page.
	sitemap *sitemap.Sitemap
	// searcher finds records of the public catalogue by text.
	searcher          *search.Catalogue
	lessonReviews     *lessonreview.Scheduler
	relatedSuggester  *relatedcourse.Scheduler
	landingPageQna    *qnaitem.LandingPage
	categoryQna       *qnaitem.Category
	postEditor        *blog.Editor
	reviewDirectory   *coursereview.Directory
	courseDirectory   *catalog.CourseDirectory
	leadDirectory     *leadcapture.Directory
	messageDirectory  *assistant.Directory
	categoryDirectory *catalog.CategoryDirectory
	bannerDirectory   *sitebanner.Directory
	lessonDirectory   *curriculum.Directory
	// oauth serves the consent screen of the OAuth server; nil while it is off.
	oauth *oauthserver.Server
}

// Deps are the handler's collaborators. A struct rather than positional
// parameters: the fields are the concrete module types (ADR-0025), the only
// mocks being the external boundaries, and a test that never reaches a
// collaborator (no lead is submitted through the
// auth router) simply leaves its field zero instead of passing a labelled nil.
// Plain fields, no DI tags, so the handlers package stays injector-agnostic.
type Deps struct {
	DB                *ent.Client
	Settings          Settings
	Starter           *versionbuild.Starter
	Jobs              *jobs.Enqueuer
	RelatedCoursesSet *relatedcourse.Replacer
	Progress          *progress.Progress
	Assets            *assetstore.Store
	Sessions          *auth.Authenticator
	Registrar         *account.Registrar
	Remover           *account.Remover
	UserEditor        *account.Editor
	UserDirectory     *account.Directory
	Events            events.StandalonePublisher
	Leads             *leadcapture.Recorder
	Books             *book.Recorder
	I18n              *localization.Translator
	Errors            *APIErrorHandler
	YandexFeed        *feed.Yandex
	Assistant         *assistant.Chat
	PostRemover       *blog.Remover
	Permissions       *permission.Permissions
	SelfProtection    *selfprotection.Guard
	BannerStates      *sitebanner.Editor
	ReviewStates      *coursereview.Editor
	LandingPageStates *landingpage.Editor
	CourseReadiness   *catalog.Readiness
	Dashboard         *dashboard.Dashboard
	Sitemap           *sitemap.Sitemap
	Search            *search.Catalogue
	LessonReviews     *lessonreview.Scheduler
	RelatedSuggester  *relatedcourse.Scheduler
	LandingPageQna    *qnaitem.LandingPage
	CategoryQna       *qnaitem.Category
	PostEditor        *blog.Editor
	ReviewDirectory   *coursereview.Directory
	CourseDirectory   *catalog.CourseDirectory
	LeadDirectory     *leadcapture.Directory
	MessageDirectory  *assistant.Directory
	CategoryDirectory *catalog.CategoryDirectory
	BannerDirectory   *sitebanner.Directory
	LessonDirectory   *curriculum.Directory
	OAuth             *oauthserver.Server
}

// NewServer wires the handler to its dependencies.
func NewServer(deps Deps) *Server {
	secure := deps.Settings.SecureCookies
	guests := newGuestCookies(deps.Settings.GuestCookieSecret, secure)

	return &Server{
		db:                deps.DB,
		conv:              &apiconv.ConverterImpl{},
		settings:          deps.Settings,
		starter:           deps.Starter,
		jobs:              deps.Jobs,
		relatedCoursesSet: deps.RelatedCoursesSet,
		progress:          deps.Progress,
		assets:            deps.Assets,
		sessions:          deps.Sessions,
		security:          newSecurity(deps.DB, deps.Sessions, permission.NewPolicy(deps.DB), deps.I18n, deps.Errors, guests),
		credentials:       account.NewCredentials(deps.DB),
		registrar:         deps.Registrar,
		remover:           deps.Remover,
		userEditor:        deps.UserEditor,
		userDirectory:     deps.UserDirectory,
		events:            deps.Events,
		emailTokens:       emailtoken.New(deps.Settings.EmailTokenSecret),
		guests:            guests,
		i18n:              deps.I18n,
		leads:             deps.Leads,
		books:             deps.Books,
		errors:            deps.Errors,
		yandexFeed:        deps.YandexFeed,
		catalog:           landingpage.NewCatalog(deps.DB),
		coursePages:       coursepage.New(deps.DB, landingpage.NewCatalog(deps.DB), deps.Progress, deps.Settings.CourseRepoBaseURL),
		assistant:         deps.Assistant,
		postRemover:       deps.PostRemover,
		blogPublic:        blog.NewPublic(deps.DB),
		blogReader:        blog.NewReader(deps.DB, deps.Settings.PublicURL),
		blogLikes:         blog.NewLiker(deps.DB),
		permissions:       deps.Permissions,
		selfProtection:    deps.SelfProtection,
		bannerStates:      deps.BannerStates,
		reviewStates:      deps.ReviewStates,
		landingPageStates: deps.LandingPageStates,
		courseReadiness:   deps.CourseReadiness,
		dashboard:         deps.Dashboard,
		sitemap:           deps.Sitemap,
		searcher:          deps.Search,
		lessonReviews:     deps.LessonReviews,
		relatedSuggester:  deps.RelatedSuggester,
		landingPageQna:    deps.LandingPageQna,
		categoryQna:       deps.CategoryQna,
		postEditor:        deps.PostEditor,
		reviewDirectory:   deps.ReviewDirectory,
		courseDirectory:   deps.CourseDirectory,
		leadDirectory:     deps.LeadDirectory,
		messageDirectory:  deps.MessageDirectory,
		categoryDirectory: deps.CategoryDirectory,
		bannerDirectory:   deps.BannerDirectory,
		lessonDirectory:   deps.LessonDirectory,
		oauth:             deps.OAuth,
	}
}

// Security returns the seam that recognises a request by its session: ogen's
// generated security handler and the pipeline that wraps the generated server.
func (s *Server) Security() *Security {
	return s.security
}

// ListCourses returns the published course catalog.
//
// URL stays `/languages` for backward-compat; the domain concept is Course.
// Mirrors the legacy scope (`LandingPage.web.where(listed: true)`): published,
// listed landing pages of the request locale whose Course is completed, in
// catalog order (see landingpage.Catalog).
func (s *Server) ListCourses(ctx context.Context) ([]api.CourseCatalogItem, error) {
	pages, err := s.catalog.Pages(ctx, s.i18n.Locale(ctx))
	if err != nil {
		return nil, err
	}

	return s.conv.ToCatalogItems(pages), nil
}
