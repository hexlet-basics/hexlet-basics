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
	"hexletbasics/internal/coursepage"
	"hexletbasics/internal/emailtoken"
	"hexletbasics/internal/events"
	"hexletbasics/internal/feed"
	"hexletbasics/internal/landingpage"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/permission"
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
	starter  VersionBuildStarter
	// jobs schedules the background work handlers trigger: lesson reviews, the
	// related-courses pick, account emails.
	jobs JobEnqueuer
	// relatedCoursesSet replaces a post's related courses in one transaction.
	relatedCoursesSet RelatedCoursesReplacer
	// progress owns what a learner has finished; handlers never derive it.
	progress ProgressTracker
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
	registrar   UserRegistrar
	remover     AccountRemover
	// userEditor and userDirectory are the administrators' write and read side of
	// user records.
	userEditor    UserEditor
	userDirectory UserDirectory
	events        events.StandalonePublisher
	// emailTokens verifies the Magic Link and Password Reset links.
	emailTokens *emailtoken.Tokens
	// guests is the visitor-progress cookie.
	guests *guestCookies
	// leads stores a lead and raises LeadCreated in one transaction.
	leads LeadCreator
	// books stores a book request and raises BookRequested in one transaction.
	books  BookRequester
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
	postRemover PostRemover
	// permissions writes a role's grants.
	permissions       PermissionWriter
	selfProtection    SelfProtection
	bannerStates      StateMover
	reviewStates      StateMover
	landingPageStates StateMover
	courseReadiness   ReadinessMarker
	// dashboard assembles the learner's courses with cards and progress.
	dashboard CourseLister
	// sitemap reads the rows behind the sitemap page.
	sitemap SitemapReader
	// searcher finds records of the public catalogue by text.
	searcher          CatalogueSearcher
	lessonReviews     LessonReviewScheduler
	relatedSuggester  RelatedCoursesScheduler
	landingPageQna    LandingPageQna
	categoryQna       CategoryQna
	postEditor        PostEditor
	reviewDirectory   ReviewDirectory
	courseDirectory   CourseDirectory
	leadDirectory     LeadDirectory
	messageDirectory  MessageDirectory
	categoryDirectory CategoryDirectory
	bannerDirectory   BannerDirectory
	lessonDirectory   LessonDirectory
	// oauth serves the consent screen of the OAuth server; nil while it is off.
	oauth OAuthConsents
}

// Deps are the handler's collaborators. A struct rather than positional
// parameters: most are interfaces the tests fill with moq mocks, and
// a test that never reaches a collaborator (no lead is submitted through the
// auth router) simply leaves its field zero instead of passing a labelled nil.
// Plain fields, no DI tags, so the handlers package stays injector-agnostic.
type Deps struct {
	DB                *ent.Client
	Settings          Settings
	Starter           VersionBuildStarter
	Jobs              JobEnqueuer
	RelatedCoursesSet RelatedCoursesReplacer
	Progress          ProgressTracker
	Assets            *assetstore.Store
	Sessions          *auth.Authenticator
	Registrar         UserRegistrar
	Remover           AccountRemover
	UserEditor        UserEditor
	UserDirectory     UserDirectory
	Events            events.StandalonePublisher
	Leads             LeadCreator
	Books             BookRequester
	I18n              *localization.Translator
	Errors            *APIErrorHandler
	YandexFeed        *feed.Yandex
	Assistant         *assistant.Chat
	PostRemover       PostRemover
	Permissions       PermissionWriter
	SelfProtection    SelfProtection
	BannerStates      StateMover
	ReviewStates      StateMover
	LandingPageStates StateMover
	CourseReadiness   ReadinessMarker
	Dashboard         CourseLister
	Sitemap           SitemapReader
	Search            CatalogueSearcher
	LessonReviews     LessonReviewScheduler
	RelatedSuggester  RelatedCoursesScheduler
	LandingPageQna    LandingPageQna
	CategoryQna       CategoryQna
	PostEditor        PostEditor
	ReviewDirectory   ReviewDirectory
	CourseDirectory   CourseDirectory
	LeadDirectory     LeadDirectory
	MessageDirectory  MessageDirectory
	CategoryDirectory CategoryDirectory
	BannerDirectory   BannerDirectory
	LessonDirectory   LessonDirectory
	OAuth             OAuthConsents
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
