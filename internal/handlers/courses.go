package handlers

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/apiconv"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/blogposts"
	"hexletbasics/internal/config"
	"hexletbasics/internal/coursepages"
	"hexletbasics/internal/emailtokens"
	"hexletbasics/internal/events"
	"hexletbasics/internal/feeds"
	"hexletbasics/internal/landingpages"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/sessions"
	"hexletbasics/internal/staffroles"
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
	sessions *sessions.Manager
	security *Security
	// credentials checks and changes passwords; registrar creates accounts and
	// remover removes them.
	credentials *accounts.Credentials
	registrar   UserRegistrar
	remover     AccountRemover
	// userEditor and userDirectory are the administrators' write and read side of
	// user records.
	userEditor    UserEditor
	userDirectory UserDirectory
	events        events.StandalonePublisher
	// emailTokens verifies the Magic Link and Password Reset links.
	emailTokens *emailtokens.Tokens
	// guests is the visitor-progress cookie.
	guests *guestCookies
	// leads stores a lead and raises LeadCreated in one transaction.
	leads LeadCreator
	// books stores a book request and raises BookRequested in one transaction.
	books  BookRequester
	i18n   *localization.Translator
	errors *APIErrorHandler
	// yandexFeed builds the Yandex course catalogue behind the feed routes.
	yandexFeed *feeds.Yandex
	// catalog reads the listed landing pages the course list and category pages share.
	catalog *landingpages.Catalog
	// coursePages reads the course page and the lesson page.
	coursePages *coursepages.Pages
	// assistant is the in-lesson chat: quota, prompt, stream and storage.
	assistant *assistant.Chat
	// blogPublic reads the posts a visitor may open.
	blogPublic *blogposts.Public
	// blogReader assembles posts (cover, likes, courses, reading time) for the
	// admin and the public lists alike; blogLikes records a like.
	blogReader *blogposts.Reader
	blogLikes  *blogposts.Liker
	// postRemover deletes a post with its dependents in one transaction.
	postRemover PostRemover
	// permissions writes a role's grants.
	permissions PermissionWriter
	// dashboard assembles the learner's courses with cards and progress.
	dashboard CourseLister
	// sitemap reads the rows behind the sitemap page.
	sitemap           SitemapReader
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
	Config            *config.Config
	Starter           VersionBuildStarter
	Jobs              JobEnqueuer
	RelatedCoursesSet RelatedCoursesReplacer
	Progress          ProgressTracker
	Assets            *assetstore.Store
	Sessions          *sessions.Manager
	Registrar         UserRegistrar
	Remover           AccountRemover
	UserEditor        UserEditor
	UserDirectory     UserDirectory
	Events            events.StandalonePublisher
	Leads             LeadCreator
	Books             BookRequester
	I18n              *localization.Translator
	Errors            *APIErrorHandler
	YandexFeed        *feeds.Yandex
	Assistant         *assistant.Chat
	PostRemover       PostRemover
	Permissions       PermissionWriter
	Dashboard         CourseLister
	Sitemap           SitemapReader
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
	secure := deps.Config.SecureCookies()
	guests := newGuestCookies(deps.Config.GuestCookieSecret, secure)
	return &Server{
		db:                deps.DB,
		conv:              &apiconv.ConverterImpl{},
		cfg:               deps.Config,
		starter:           deps.Starter,
		jobs:              deps.Jobs,
		relatedCoursesSet: deps.RelatedCoursesSet,
		progress:          deps.Progress,
		assets:            deps.Assets,
		sessions:          deps.Sessions,
		security:          newSecurity(deps.DB, deps.Sessions, staffroles.NewPolicy(deps.DB), deps.I18n, deps.Errors, guests),
		credentials:       accounts.NewCredentials(deps.DB),
		registrar:         deps.Registrar,
		remover:           deps.Remover,
		userEditor:        deps.UserEditor,
		userDirectory:     deps.UserDirectory,
		events:            deps.Events,
		emailTokens:       emailtokens.New(deps.Config.EmailTokenSecret),
		guests:            guests,
		i18n:              deps.I18n,
		leads:             deps.Leads,
		books:             deps.Books,
		errors:            deps.Errors,
		yandexFeed:        deps.YandexFeed,
		catalog:           landingpages.NewCatalog(deps.DB),
		coursePages:       coursepages.New(deps.DB, landingpages.NewCatalog(deps.DB), deps.Progress, deps.Config.CourseRepoBaseURL),
		assistant:         deps.Assistant,
		postRemover:       deps.PostRemover,
		blogPublic:        blogposts.NewPublic(deps.DB),
		blogReader:        blogposts.NewReader(deps.DB, deps.Config.PublicURL),
		blogLikes:         blogposts.NewLiker(deps.DB),
		permissions:       deps.Permissions,
		dashboard:         deps.Dashboard,
		sitemap:           deps.Sitemap,
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
// catalog order (see landingpages.Catalog).
func (s *Server) ListCourses(ctx context.Context) ([]api.CourseCatalogItem, error) {
	pages, err := s.catalog.Pages(ctx, s.i18n.Locale(ctx))
	if err != nil {
		return nil, err
	}

	return s.conv.ToCatalogItems(pages), nil
}
