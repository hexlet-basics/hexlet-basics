package di

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/contrib/otelconf"

	"hexletbasics/ent"
	"hexletbasics/internal/accounts"
	"hexletbasics/internal/api"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/banners"
	"hexletbasics/internal/blogposts"
	"hexletbasics/internal/books"
	"hexletbasics/internal/config"
	"hexletbasics/internal/coursecategories"
	"hexletbasics/internal/courselessons"
	"hexletbasics/internal/courses"
	"hexletbasics/internal/dashboard"
	"hexletbasics/internal/events"
	"hexletbasics/internal/exerciserunner"
	"hexletbasics/internal/feeds"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/landingpages"
	"hexletbasics/internal/leads"
	"hexletbasics/internal/lessonreviews"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/qnaitems"
	"hexletbasics/internal/relatedcourses"
	"hexletbasics/internal/reviews"
	"hexletbasics/internal/search"
	"hexletbasics/internal/selfprotection"
	"hexletbasics/internal/sessions"
	"hexletbasics/internal/sitemap"
	"hexletbasics/internal/staffroles"
	"hexletbasics/internal/store"
	"hexletbasics/internal/telemetry"
	"hexletbasics/internal/versionbuilds"
)

// HTTP server timeouts. ReadHeaderTimeout caps slow-header (Slowloris) clients;
// the rest bound a misbehaving peer without cutting normal API traffic.
const (
	readHeaderTimeout = 10 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 120 * time.Second

	routerServiceName      = "router"
	httpHandlerServiceName = "http-handler"
)

// serverPackage contains the synchronous graph: HTTP delivery, event
// publishing, and insert-only River access.
// handlerServices is what handlers.Server is built from, resolved by concrete
// type: handlers.Deps holds interfaces, which the container cannot tell apart
// by type alone.
type handlerServices struct {
	DB                *ent.Client                 `do:""`
	Config            *config.Config              `do:""`
	Starter           *versionbuilds.Starter      `do:""`
	Jobs              *jobs.Enqueuer              `do:""`
	RelatedCoursesSet *relatedcourses.Replacer    `do:""`
	Progress          *progress.Progress          `do:""`
	Assets            *assetstore.Store           `do:""`
	Sessions          *sessions.Manager           `do:""`
	Registrar         *accounts.Registrar         `do:""`
	Remover           *accounts.Remover           `do:""`
	UserEditor        *accounts.Editor            `do:""`
	UserDirectory     *accounts.Directory         `do:""`
	Events            *events.Publisher           `do:""`
	Leads             *leads.Recorder             `do:""`
	Books             *books.Recorder             `do:""`
	I18n              *localization.Translator    `do:""`
	Errors            *handlers.APIErrorHandler   `do:""`
	YandexFeed        *feeds.Yandex               `do:""`
	Assistant         *assistant.Chat             `do:""`
	PostRemover       *blogposts.Remover          `do:""`
	Permissions       *staffroles.Permissions     `do:""`
	SelfProtection    *selfprotection.Guard       `do:""`
	BannerStates      *banners.Editor             `do:""`
	ReviewStates      *reviews.Editor             `do:""`
	LandingPageStates *landingpages.Editor        `do:""`
	CourseReadiness   *courses.Readiness          `do:""`
	Dashboard         *dashboard.Dashboard        `do:""`
	Sitemap           *sitemap.Sitemap            `do:""`
	Search            *search.Catalogue           `do:""`
	LessonReviews     *lessonreviews.Scheduler    `do:""`
	RelatedSuggester  *relatedcourses.Scheduler   `do:""`
	LandingPageQna    *qnaitems.LandingPage       `do:""`
	CategoryQna       *qnaitems.Category          `do:""`
	PostEditor        *blogposts.Editor           `do:""`
	ReviewDirectory   *reviews.Directory          `do:""`
	CourseDirectory   *courses.Directory          `do:""`
	LeadDirectory     *leads.Directory            `do:""`
	MessageDirectory  *assistant.Directory        `do:""`
	CategoryDirectory *coursecategories.Directory `do:""`
	BannerDirectory   *banners.Directory          `do:""`
	LessonDirectory   *courselessons.Directory    `do:""`
}

var serverPackage = do.Package(
	do.Lazy[*handlers.APIErrorHandler](func(i do.Injector) (*handlers.APIErrorHandler, error) {
		translator, err := do.Invoke[*localization.Translator](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return handlers.NewAPIErrorHandler(translator, logger), nil
	}),
	// Sessions live in the database; a store failure answers in the API's own
	// error shape through the same handler everything else uses.
	do.Lazy[*sessions.Manager](func(i do.Injector) (*sessions.Manager, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		errorHandler, err := do.Invoke[*handlers.APIErrorHandler](i)
		if err != nil {
			return nil, err
		}
		return sessions.New(
			db,
			sessions.WithSecure(cfg.SecureCookies()),
			sessions.WithErrorFunc(func(w http.ResponseWriter, r *http.Request, err error) {
				errorHandler.Write(r.Context(), w, r, err)
			}),
		), nil
	}),
	do.Lazy[*accounts.Registrar](func(i do.Injector) (*accounts.Registrar, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		publisher, err := do.Invoke[*events.Publisher](i)
		if err != nil {
			return nil, err
		}
		return accounts.NewRegistrar(db, publisher), nil
	}),
	do.Lazy[*leads.Recorder](func(i do.Injector) (*leads.Recorder, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		publisher, err := do.Invoke[*events.Publisher](i)
		if err != nil {
			return nil, err
		}
		return leads.NewRecorder(db, publisher), nil
	}),
	do.Lazy[*books.Recorder](func(i do.Injector) (*books.Recorder, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		publisher, err := do.Invoke[*events.Publisher](i)
		if err != nil {
			return nil, err
		}
		return books.NewRecorder(db, publisher), nil
	}),
	do.Lazy[*blogposts.Remover](func(i do.Injector) (*blogposts.Remover, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		return blogposts.NewRemover(db), nil
	}),
	do.Lazy[*banners.Editor](func(i do.Injector) (*banners.Editor, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return banners.NewEditor(db), nil
	}),
	do.Lazy[*reviews.Editor](func(i do.Injector) (*reviews.Editor, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return reviews.NewEditor(db), nil
	}),
	do.Lazy[*landingpages.Editor](func(i do.Injector) (*landingpages.Editor, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return landingpages.NewEditor(db), nil
	}),
	do.Lazy[*courses.Readiness](func(i do.Injector) (*courses.Readiness, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return courses.NewReadiness(db), nil
	}),
	do.Lazy[*selfprotection.Guard](func(i do.Injector) (*selfprotection.Guard, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return selfprotection.NewGuard(db), nil
	}),
	do.Lazy[*staffroles.Permissions](func(i do.Injector) (*staffroles.Permissions, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		tx, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		return staffroles.NewPermissions(db, tx), nil
	}),
	do.Lazy[*courselessons.Directory](func(i do.Injector) (*courselessons.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return courselessons.NewDirectory(db), nil
	}),
	do.Lazy[*banners.Directory](func(i do.Injector) (*banners.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return banners.NewDirectory(db), nil
	}),
	do.Lazy[*coursecategories.Directory](func(i do.Injector) (*coursecategories.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return coursecategories.NewDirectory(db), nil
	}),
	do.Lazy[*assistant.Directory](func(i do.Injector) (*assistant.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return assistant.NewDirectory(db), nil
	}),
	do.Lazy[*leads.Directory](func(i do.Injector) (*leads.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return leads.NewDirectory(db), nil
	}),
	do.Lazy[*courses.Directory](func(i do.Injector) (*courses.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return courses.NewDirectory(db), nil
	}),
	do.Lazy[*reviews.Directory](func(i do.Injector) (*reviews.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return reviews.NewDirectory(db), nil
	}),
	do.Lazy[*blogposts.Editor](func(i do.Injector) (*blogposts.Editor, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return blogposts.NewEditor(db), nil
	}),
	do.Lazy[*qnaitems.Category](func(i do.Injector) (*qnaitems.Category, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return qnaitems.NewCategory(db), nil
	}),
	do.Lazy[*qnaitems.LandingPage](func(i do.Injector) (*qnaitems.LandingPage, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return qnaitems.NewLandingPage(db), nil
	}),
	do.Lazy[*lessonreviews.Scheduler](func(i do.Injector) (*lessonreviews.Scheduler, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		enqueuer, err := do.Invoke[*jobs.Enqueuer](i)
		if err != nil {
			return nil, err
		}
		return lessonreviews.NewScheduler(db, enqueuer), nil
	}),
	do.Lazy[*relatedcourses.Scheduler](func(i do.Injector) (*relatedcourses.Scheduler, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		enqueuer, err := do.Invoke[*jobs.Enqueuer](i)
		if err != nil {
			return nil, err
		}
		return relatedcourses.NewScheduler(db, enqueuer), nil
	}),
	do.Lazy[*search.Catalogue](func(i do.Injector) (*search.Catalogue, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return search.NewCatalogue(db), nil
	}),
	do.Lazy[*sitemap.Sitemap](func(i do.Injector) (*sitemap.Sitemap, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return sitemap.New(db), nil
	}),
	do.Lazy[*dashboard.Dashboard](func(i do.Injector) (*dashboard.Dashboard, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		tracker, err := do.Invoke[*progress.Progress](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return dashboard.New(db, landingpages.NewCatalog(db), tracker, logger), nil
	}),
	do.Lazy[*accounts.Editor](func(i do.Injector) (*accounts.Editor, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return accounts.NewEditor(db), nil
	}),
	do.Lazy[*accounts.Directory](func(i do.Injector) (*accounts.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return accounts.NewDirectory(db), nil
	}),
	do.Lazy[*accounts.Remover](func(i do.Injector) (*accounts.Remover, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		return accounts.NewRemover(db), nil
	}),
	do.Lazy[*river.Client[*sql.Tx]](func(i do.Injector) (*river.Client[*sql.Tx], error) {
		db, err := do.Invoke[*sql.DB](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		otelSDK, err := do.Invoke[*otelconf.SDK](i)
		if err != nil {
			return nil, err
		}
		return jobs.NewInsertOnlyClient(
			db,
			logger,
			otelSDK.TracerProvider(),
			otelSDK.MeterProvider(),
		)
	}),
	do.Lazy[*jobs.Enqueuer](func(i do.Injector) (*jobs.Enqueuer, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		riverClient, err := do.Invoke[*river.Client[*sql.Tx]](i)
		if err != nil {
			return nil, err
		}
		return jobs.NewEnqueuer(db, riverClient), nil
	}),
	do.Lazy[*versionbuilds.Starter](func(i do.Injector) (*versionbuilds.Starter, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		enqueuer, err := do.Invoke[*jobs.Enqueuer](i)
		if err != nil {
			return nil, err
		}
		return versionbuilds.NewStarter(db, enqueuer), nil
	}),
	// The Yandex feed builds its URLs on the canonical host, not the request's.
	do.Lazy[*feeds.Yandex](func(i do.Injector) (*feeds.Yandex, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		return feeds.NewYandex(db, cfg.AppHost), nil
	}),
	do.Lazy[*handlers.Server](func(i do.Injector) (*handlers.Server, error) {
		services, err := do.InvokeStruct[handlerServices](i)
		if err != nil {
			return nil, err
		}
		// A nil *oauthserver.Server must not become a non-nil interface: the
		// consent operations treat a nil seam as "the OAuth server is off".
		var consents handlers.OAuthConsents
		if services.Config.OAuth.Enabled {
			oauthServer, err := do.Invoke[*oauthserver.Server](i)
			if err != nil {
				return nil, err
			}
			consents = oauthServer
		}
		return handlers.NewServer(handlers.Deps{
			OAuth:             consents,
			DB:                services.DB,
			Config:            services.Config,
			Starter:           services.Starter,
			Jobs:              services.Jobs,
			RelatedCoursesSet: services.RelatedCoursesSet,
			Progress:          services.Progress,
			Assets:            services.Assets,
			Sessions:          services.Sessions,
			Registrar:         services.Registrar,
			Remover:           services.Remover,
			UserEditor:        services.UserEditor,
			UserDirectory:     services.UserDirectory,
			Events:            services.Events,
			Leads:             services.Leads,
			Books:             services.Books,
			I18n:              services.I18n,
			Errors:            services.Errors,
			YandexFeed:        services.YandexFeed,
			Assistant:         services.Assistant,
			PostRemover:       services.PostRemover,
			Permissions:       services.Permissions,
			SelfProtection:    services.SelfProtection,
			BannerStates:      services.BannerStates,
			ReviewStates:      services.ReviewStates,
			LandingPageStates: services.LandingPageStates,
			CourseReadiness:   services.CourseReadiness,
			Dashboard:         services.Dashboard,
			Sitemap:           services.Sitemap,
			Search:            services.Search,
			LessonReviews:     services.LessonReviews,
			RelatedSuggester:  services.RelatedSuggester,
			LandingPageQna:    services.LandingPageQna,
			CategoryQna:       services.CategoryQna,
			PostEditor:        services.PostEditor,
			ReviewDirectory:   services.ReviewDirectory,
			CourseDirectory:   services.CourseDirectory,
			LeadDirectory:     services.LeadDirectory,
			MessageDirectory:  services.MessageDirectory,
			CategoryDirectory: services.CategoryDirectory,
			BannerDirectory:   services.BannerDirectory,
			LessonDirectory:   services.LessonDirectory,
		}), nil
	}),
	// The in-lesson assistant answers inside the request (a stream, no job), so
	// its model client lives in the HTTP process. Unlike the worker's LLM jobs
	// it is built even without credentials: a question then fails at the
	// provider before the first token and the learner gets an error status,
	// where a skipped job could simply wait in the queue.
	do.Lazy[*assistant.Chat](func(i do.Injector) (*assistant.Chat, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		txStore, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		tracker, err := do.Invoke[*progress.Progress](i)
		if err != nil {
			return nil, err
		}
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		translator, err := do.Invoke[*localization.Translator](i)
		if err != nil {
			return nil, err
		}
		llm := assistant.NewOpenAI(cfg.OpenAIAccessToken, cfg.OpenAIModel)
		return assistant.NewChat(db, txStore, tracker, llm, translator), nil
	}),
	// The progress module owns what a learner has finished. It writes and publishes
	// through the same transaction seam every other business module uses.
	do.Lazy[*progress.Progress](func(i do.Injector) (*progress.Progress, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		txStore, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		publisher, err := do.Invoke[*events.Publisher](i)
		if err != nil {
			return nil, err
		}
		runner, err := do.Invoke[*exerciserunner.Docker](i)
		if err != nil {
			return nil, err
		}
		return progress.New(db, txStore, publisher, runner), nil
	}),
	// The exercise runner is the only Docker consumer in the app (ADR-0013).
	// The client resolves the daemon from the environment and connects lazily,
	// so a deployment with no daemon reachable fails on the first check rather
	// than refusing to boot — which is what keeps every other surface serving.
	do.Lazy[*exerciserunner.Docker](func(i do.Injector) (*exerciserunner.Docker, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		return exerciserunner.New(exerciserunner.OptionsFrom(cfg.ExerciseRunner))
	}),
	do.Lazy[*api.Server](func(i do.Injector) (*api.Server, error) {
		handler, err := do.Invoke[*handlers.Server](i)
		if err != nil {
			return nil, err
		}
		errorHandler, err := do.Invoke[*handlers.APIErrorHandler](i)
		if err != nil {
			return nil, err
		}
		otelSDK, err := do.Invoke[*otelconf.SDK](i)
		if err != nil {
			return nil, err
		}
		translator, err := do.Invoke[*localization.Translator](i)
		if err != nil {
			return nil, err
		}
		txStore, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		// WithMiddleware is the unit of work of every admin write: the
		// transaction, and the Audit Record written inside it (ADR-0020).
		// WithErrorHandler installs the central ent-error -> HTTP-status
		// mapping (404/409), so handlers return raw ent errors instead of
		// typed DTOs.
		publisher, err := do.Invoke[*events.Publisher](i)
		if err != nil {
			return nil, err
		}
		return api.NewServer(
			handler,
			handler.Security(),
			api.WithErrorHandler(errorHandler.Write),
			// One call: a second WithMiddleware would replace the first. The
			// second middleware is the staff reads of personal data: the server
			// only publishes the fact, the worker writes the Audit Record
			// (ADR-0020, ADR-0010).
			api.WithMiddleware(
				handlers.NewAuditMiddleware(txStore, translator, errorHandler),
				handlers.NewPersonalDataReadMiddleware(publisher, translator, errorHandler),
			),
			api.WithTracerProvider(otelSDK.TracerProvider()),
			api.WithNotFound(handlers.NewNotFoundHandler(translator)),
			api.WithMethodNotAllowed(handlers.NewMethodNotAllowedHandler(translator)),
		)
	}),
	do.Lazy[*handlers.AttachmentHandler](func(i do.Injector) (*handlers.AttachmentHandler, error) {
		assets, err := do.Invoke[*assetstore.Store](i)
		if err != nil {
			return nil, err
		}
		errorHandler, err := do.Invoke[*handlers.APIErrorHandler](i)
		if err != nil {
			return nil, err
		}
		return handlers.NewAttachmentHandler(assets, errorHandler), nil
	}),
	do.Lazy[*handlers.GitHubWebhookHandler](func(i do.Injector) (*handlers.GitHubWebhookHandler, error) {
		starter, err := do.Invoke[*versionbuilds.Starter](i)
		if err != nil {
			return nil, err
		}
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		translator, err := do.Invoke[*localization.Translator](i)
		if err != nil {
			return nil, err
		}
		return handlers.NewGitHubWebhookHandler(starter, cfg.GitHubWebhookSecret, translator), nil
	}),
	// Both the raw router and the middleware-wrapped application handler have
	// the same http.Handler interface. Named bindings keep those two seams
	// explicit without adding wrapper types whose only purpose would be DI.
	do.LazyNamed[http.Handler](routerServiceName, func(i do.Injector) (http.Handler, error) {
		apiServer, err := do.Invoke[*api.Server](i)
		if err != nil {
			return nil, err
		}
		attachmentHandler, err := do.Invoke[*handlers.AttachmentHandler](i)
		if err != nil {
			return nil, err
		}
		webhookHandler, err := do.Invoke[*handlers.GitHubWebhookHandler](i)
		if err != nil {
			return nil, err
		}
		handler, err := do.Invoke[*handlers.Server](i)
		if err != nil {
			return nil, err
		}
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		// The OAuth server is mounted only behind its flag (ADR-0022).
		var oauthHandler http.Handler
		if cfg.OAuth.Enabled {
			oauthServer, err := do.Invoke[*oauthserver.Server](i)
			if err != nil {
				return nil, err
			}
			oauthHandler = oauthServer.Handler()
		}
		// The MCP endpoint is mounted only behind its flag (ADR-0021). A nil
		// *mcpserver.Server must not become a non-nil http.Handler.
		var mcpHandler http.Handler
		if cfg.MCP.Enabled {
			mcpEndpoint, err := do.Invoke[*mcpserver.Server](i)
			if err != nil {
				return nil, err
			}
			mcpHandler = mcpEndpoint.Handler()
		}
		return handlers.NewRouter(
			handler.Security().Pipeline(apiServer),
			attachmentHandler,
			webhookHandler,
			oauthHandler,
			mcpHandler,
		), nil
	}),
	// The OAuth authorization server (ADR-0022). It receives the application's
	// half of its seam, who is signed in, from the handlers package; the
	// package itself never imports handlers. Built only when the flag is on,
	// since an invalid issuer or a missing key would otherwise stop a
	// deployment that does not use it.
	do.Lazy[*oauthserver.Server](func(i do.Injector) (*oauthserver.Server, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		txStore, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		// Off, the service resolves to nil: the router and the consent seam
		// only ask for it behind the same flag.
		if !cfg.OAuth.Enabled {
			return nil, nil
		}
		return oauthserver.New(oauthserver.Deps{
			Config: oauthserver.Config{
				Issuer:               cfg.OAuth.Issuer,
				EncryptionKey:        cfg.OAuth.EncryptionKey,
				RegistrationLimit:    cfg.OAuth.RegistrationLimit,
				RegistrationWindow:   cfg.OAuth.RegistrationWindow,
				AccessTokenLifetime:  cfg.OAuth.AccessTokenTTL,
				RefreshTokenLifetime: cfg.OAuth.RefreshTokenTTL,
			},
			DB:            db,
			Transactor:    txStore,
			App:           handlers.OAuthSession{},
			ClientAddress: handlers.ClientIP,
		})
	}),
	// The MCP endpoint (ADR-0021): verifies tokens with the OAuth server and
	// serves tools through the same handler code as the HTTP operations. Built
	// only when its flag is on; config validation guarantees OAuth is on then.
	do.Lazy[*mcpserver.Server](func(i do.Injector) (*mcpserver.Server, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		if !cfg.MCP.Enabled {
			return nil, nil
		}
		oauthServer, err := do.Invoke[*oauthserver.Server](i)
		if err != nil {
			return nil, err
		}
		handler, err := do.Invoke[*handlers.Server](i)
		if err != nil {
			return nil, err
		}
		apiServer, err := do.Invoke[*api.Server](i)
		if err != nil {
			return nil, err
		}
		translator, err := do.Invoke[*localization.Translator](i)
		if err != nil {
			return nil, err
		}
		return mcpserver.New(mcpserver.Deps{
			Verifier:  oauthServer,
			Catalogue: handler,
			// The generated server with its security handler and audit
			// middleware: the one checkpoint of every admin tool (ADR-0021).
			Learner:      handler,
			Capabilities: handler.Security(),
			API:          handlers.NewAdminDispatcher(apiServer, translator),
			Issuer:       cfg.OAuth.Issuer,
			// The links of search results and fetched records point at the public
			// site, with the locale prefix of the call.
			SiteURL: cfg.SiteURL,
			Locale:  translator.Locale,
			Limits: mcpserver.Limits{
				MaxPageSize:      cfg.MCP.MaxPageSize,
				MaxResponseBytes: cfg.MCP.MaxResponseBytes,
				ToolTimeout:      cfg.MCP.ToolTimeout,
				RateWindow:       cfg.MCP.RateWindow,
				User: mcpserver.Rates{
					Catalogue: cfg.MCP.RateUserCatalogue,
					Personal:  cfg.MCP.RateUserPersonal,
					Write:     cfg.MCP.RateUserWrite,
				},
				Client: mcpserver.Rates{
					Catalogue: cfg.MCP.RateClientCatalogue,
					Personal:  cfg.MCP.RateClientPersonal,
					Write:     cfg.MCP.RateClientWrite,
				},
			},
		})
	}),
	do.LazyNamed[http.Handler](httpHandlerServiceName, func(i do.Injector) (http.Handler, error) {
		translator, err := do.Invoke[*localization.Translator](i)
		if err != nil {
			return nil, err
		}
		router, err := do.InvokeNamed[http.Handler](i, routerServiceName)
		if err != nil {
			return nil, err
		}
		sentryClient, err := do.Invoke[*sentry.Client](i)
		if err != nil {
			return nil, err
		}
		// No CORS: the browser reaches the API on the site's own origin under
		// `/api` (ADR-0015), through the ingress in production and the Vite
		// proxy in development.
		localized := translator.Middleware(router)
		return telemetry.NewSentryHTTPHandler(sentryClient, localized), nil
	}),
	// The process lifecycle coordinator starts and gracefully stops this
	// server. Keeping the provider on the vendor type avoids coupling DI to
	// supervision.
	do.Lazy[*http.Server](func(i do.Injector) (*http.Server, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}
		handler, err := do.InvokeNamed[http.Handler](i, httpHandlerServiceName)
		if err != nil {
			return nil, err
		}

		return &http.Server{
			Addr:              cfg.Addr,
			Handler:           handler,
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
		}, nil
	}),
)
