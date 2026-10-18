package di

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/contrib/otelconf"

	"hexletbasics/ent"
	"hexletbasics/internal/account"
	"hexletbasics/internal/api"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/assistant"
	"hexletbasics/internal/auth"
	"hexletbasics/internal/blog"
	"hexletbasics/internal/book"
	"hexletbasics/internal/catalog"
	"hexletbasics/internal/config"
	"hexletbasics/internal/coursereview"
	"hexletbasics/internal/curriculum"
	"hexletbasics/internal/dashboard"
	"hexletbasics/internal/events"
	"hexletbasics/internal/exerciserunner"
	"hexletbasics/internal/feed"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/health"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/landingpage"
	"hexletbasics/internal/leadcapture"
	"hexletbasics/internal/lessonreview"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/permission"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/qnaitem"
	"hexletbasics/internal/relatedcourse"
	"hexletbasics/internal/search"
	"hexletbasics/internal/selfprotection"
	"hexletbasics/internal/sitebanner"
	"hexletbasics/internal/sitemap"
	"hexletbasics/internal/store"
	"hexletbasics/internal/telemetry"
	"hexletbasics/internal/versionbuild"
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
	// healthServerServiceName names the worker's health-only *http.Server.
	healthServerServiceName = "health-server"
)

// serverPackage contains the synchronous graph: HTTP delivery, event
// publishing, and insert-only River access.
// handlerServices is what handlers.Server is built from, resolved by concrete
// type, the same types handlers.Deps holds.
type handlerServices struct {
	DB                *ent.Client                `do:""`
	Config            *config.Server             `do:""`
	Starter           *versionbuild.Starter      `do:""`
	Jobs              *jobs.Enqueuer             `do:""`
	RelatedCoursesSet *relatedcourse.Replacer    `do:""`
	Progress          *progress.Progress         `do:""`
	Assets            *assetstore.Store          `do:""`
	Sessions          *auth.Authenticator        `do:""`
	Registrar         *account.Registrar         `do:""`
	Remover           *account.Remover           `do:""`
	UserEditor        *account.Editor            `do:""`
	UserDirectory     *account.Directory         `do:""`
	Events            *events.Publisher          `do:""`
	Leads             *leadcapture.Recorder      `do:""`
	Books             *book.Recorder             `do:""`
	I18n              *localization.Translator   `do:""`
	Errors            *handlers.APIErrorHandler  `do:""`
	YandexFeed        *feed.Yandex               `do:""`
	Assistant         *assistant.Chat            `do:""`
	PostRemover       *blog.Remover              `do:""`
	Permissions       *permission.Permissions    `do:""`
	SelfProtection    *selfprotection.Guard      `do:""`
	BannerStates      *sitebanner.Editor         `do:""`
	ReviewStates      *coursereview.Editor       `do:""`
	LandingPageStates *landingpage.Editor        `do:""`
	CourseReadiness   *catalog.Readiness         `do:""`
	Dashboard         *dashboard.Dashboard       `do:""`
	Sitemap           *sitemap.Sitemap           `do:""`
	Search            *search.Catalogue          `do:""`
	LessonReviews     *lessonreview.Scheduler    `do:""`
	RelatedSuggester  *relatedcourse.Scheduler   `do:""`
	LandingPageQna    *qnaitem.LandingPage       `do:""`
	CategoryQna       *qnaitem.Category          `do:""`
	PostEditor        *blog.Editor               `do:""`
	ReviewDirectory   *coursereview.Directory    `do:""`
	CourseDirectory   *catalog.CourseDirectory   `do:""`
	LeadDirectory     *leadcapture.Directory     `do:""`
	MessageDirectory  *assistant.Directory       `do:""`
	CategoryDirectory *catalog.CategoryDirectory `do:""`
	BannerDirectory   *sitebanner.Directory      `do:""`
	LessonDirectory   *curriculum.Directory      `do:""`
}

var serverPackage = do.Package(
	// The process's own configuration root. Only this package (and the worker's)
	// reads the whole of it: everything else is handed the section or the values
	// it uses.
	do.Lazy[*config.Server](func(do.Injector) (*config.Server, error) {
		return config.LoadServer()
	}),
	do.Lazy[config.Common](func(i do.Injector) (config.Common, error) {
		cfg, err := do.Invoke[*config.Server](i)
		if err != nil {
			return config.Common{}, err
		}
		return cfg.Common, nil
	}),
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
	do.Lazy[*auth.Authenticator](func(i do.Injector) (*auth.Authenticator, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		cfg, err := do.Invoke[*config.Server](i)
		if err != nil {
			return nil, err
		}
		errorHandler, err := do.Invoke[*handlers.APIErrorHandler](i)
		if err != nil {
			return nil, err
		}
		return auth.New(
			db,
			auth.WithSecure(cfg.SecureCookies()),
			auth.WithErrorFunc(func(w http.ResponseWriter, r *http.Request, err error) {
				errorHandler.Write(r.Context(), w, r, err)
			}),
		), nil
	}),
	do.Lazy[*account.Registrar](func(i do.Injector) (*account.Registrar, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		publisher, err := do.Invoke[*events.Publisher](i)
		if err != nil {
			return nil, err
		}
		return account.NewRegistrar(db, publisher), nil
	}),
	do.Lazy[*leadcapture.Recorder](func(i do.Injector) (*leadcapture.Recorder, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		publisher, err := do.Invoke[*events.Publisher](i)
		if err != nil {
			return nil, err
		}
		return leadcapture.NewRecorder(db, publisher), nil
	}),
	do.Lazy[*book.Recorder](func(i do.Injector) (*book.Recorder, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		publisher, err := do.Invoke[*events.Publisher](i)
		if err != nil {
			return nil, err
		}
		return book.NewRecorder(db, publisher), nil
	}),
	do.Lazy[*blog.Remover](func(i do.Injector) (*blog.Remover, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		return blog.NewRemover(db), nil
	}),
	do.Lazy[*sitebanner.Editor](func(i do.Injector) (*sitebanner.Editor, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return sitebanner.NewEditor(db), nil
	}),
	do.Lazy[*coursereview.Editor](func(i do.Injector) (*coursereview.Editor, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return coursereview.NewEditor(db), nil
	}),
	do.Lazy[*landingpage.Editor](func(i do.Injector) (*landingpage.Editor, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return landingpage.NewEditor(db), nil
	}),
	do.Lazy[*catalog.Readiness](func(i do.Injector) (*catalog.Readiness, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return catalog.NewReadiness(db), nil
	}),
	do.Lazy[*selfprotection.Guard](func(i do.Injector) (*selfprotection.Guard, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return selfprotection.NewGuard(db), nil
	}),
	do.Lazy[*permission.Permissions](func(i do.Injector) (*permission.Permissions, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		tx, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		return permission.NewPermissions(db, tx), nil
	}),
	do.Lazy[*curriculum.Directory](func(i do.Injector) (*curriculum.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return curriculum.NewDirectory(db), nil
	}),
	do.Lazy[*sitebanner.Directory](func(i do.Injector) (*sitebanner.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return sitebanner.NewDirectory(db), nil
	}),
	do.Lazy[*catalog.CategoryDirectory](func(i do.Injector) (*catalog.CategoryDirectory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return catalog.NewCategoryDirectory(db), nil
	}),
	do.Lazy[*assistant.Directory](func(i do.Injector) (*assistant.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return assistant.NewDirectory(db), nil
	}),
	do.Lazy[*leadcapture.Directory](func(i do.Injector) (*leadcapture.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return leadcapture.NewDirectory(db), nil
	}),
	do.Lazy[*catalog.CourseDirectory](func(i do.Injector) (*catalog.CourseDirectory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return catalog.NewCourseDirectory(db), nil
	}),
	do.Lazy[*coursereview.Directory](func(i do.Injector) (*coursereview.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return coursereview.NewDirectory(db), nil
	}),
	do.Lazy[*blog.Editor](func(i do.Injector) (*blog.Editor, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return blog.NewEditor(db), nil
	}),
	do.Lazy[*qnaitem.Category](func(i do.Injector) (*qnaitem.Category, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return qnaitem.NewCategory(db), nil
	}),
	do.Lazy[*qnaitem.LandingPage](func(i do.Injector) (*qnaitem.LandingPage, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return qnaitem.NewLandingPage(db), nil
	}),
	do.Lazy[*lessonreview.Scheduler](func(i do.Injector) (*lessonreview.Scheduler, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		enqueuer, err := do.Invoke[*jobs.Enqueuer](i)
		if err != nil {
			return nil, err
		}
		return lessonreview.NewScheduler(db, enqueuer), nil
	}),
	do.Lazy[*relatedcourse.Scheduler](func(i do.Injector) (*relatedcourse.Scheduler, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		enqueuer, err := do.Invoke[*jobs.Enqueuer](i)
		if err != nil {
			return nil, err
		}
		return relatedcourse.NewScheduler(db, enqueuer), nil
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
		return dashboard.New(db, landingpage.NewCatalog(db), tracker, logger), nil
	}),
	do.Lazy[*account.Editor](func(i do.Injector) (*account.Editor, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return account.NewEditor(db), nil
	}),
	do.Lazy[*account.Directory](func(i do.Injector) (*account.Directory, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return account.NewDirectory(db), nil
	}),
	do.Lazy[*account.Remover](func(i do.Injector) (*account.Remover, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		return account.NewRemover(db), nil
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
	do.Lazy[*versionbuild.Starter](func(i do.Injector) (*versionbuild.Starter, error) {
		db, err := do.Invoke[*store.Store](i)
		if err != nil {
			return nil, err
		}
		enqueuer, err := do.Invoke[*jobs.Enqueuer](i)
		if err != nil {
			return nil, err
		}
		return versionbuild.NewStarter(db, enqueuer), nil
	}),
	// The Yandex feed builds its URLs on the canonical host, not the request's.
	do.Lazy[*feed.Yandex](func(i do.Injector) (*feed.Yandex, error) {
		db, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		cfg, err := do.Invoke[*config.Server](i)
		if err != nil {
			return nil, err
		}
		return feed.NewYandex(db, cfg.AppHost), nil
	}),
	do.Lazy[*handlers.Server](func(i do.Injector) (*handlers.Server, error) {
		services, err := do.InvokeStruct[handlerServices](i)
		if err != nil {
			return nil, err
		}
		// nil means the OAuth server is off: the consent operations answer as if
		// the request did not exist.
		var consents *oauthserver.Server
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
			Settings:          handlerSettings(services.Config),
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
		cfg, err := do.Invoke[*config.Server](i)
		if err != nil {
			return nil, err
		}
		translator, err := do.Invoke[*localization.Translator](i)
		if err != nil {
			return nil, err
		}
		llm := assistant.NewOpenAI(cfg.OpenAIAccessToken.Reveal(), cfg.OpenAIModel)
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
		cfg, err := do.Invoke[*config.Server](i)
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
		starter, err := do.Invoke[*versionbuild.Starter](i)
		if err != nil {
			return nil, err
		}
		cfg, err := do.Invoke[*config.Server](i)
		if err != nil {
			return nil, err
		}
		translator, err := do.Invoke[*localization.Translator](i)
		if err != nil {
			return nil, err
		}
		return handlers.NewGitHubWebhookHandler(starter, cfg.GitHubWebhookSecret.Reveal(), translator), nil
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
		cfg, err := do.Invoke[*config.Server](i)
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
		cfg, err := do.Invoke[*config.Server](i)
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
				EncryptionKey:        cfg.OAuth.EncryptionKey.Reveal(),
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
		cfg, err := do.Invoke[*config.Server](i)
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
			MyDashboard:      handler.MyDashboard,
			MyCourseProgress: handler.MyCourseProgress,
			HeldCapabilities: handler.Security().HeldCapabilities,
			API:              handlers.NewAdminDispatcher(apiServer, translator),
			Issuer:           cfg.OAuth.Issuer,
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
	// Readiness is the database only (ADR-0023): the injector asks every built
	// service that can check itself, and databaseCloser is the one that does.
	do.Lazy[*health.Probes](func(i do.Injector) (*health.Probes, error) {
		if _, err := do.Invoke[*databaseCloser](i); err != nil {
			return nil, err
		}
		return health.New(func(ctx context.Context) error {
			var errs []error
			for _, err := range i.HealthCheckWithContext(ctx) {
				errs = append(errs, err)
			}
			return errors.Join(errs...)
		}), nil
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
		probes, err := do.Invoke[*health.Probes](i)
		if err != nil {
			return nil, err
		}
		// The probes sit outside the contract, the locale middleware and
		// Sentry: a kubelet poll is neither a page view nor an error source.
		return probes.Wrap(telemetry.NewSentryHTTPHandler(sentryClient, localized)), nil
	}),
	// The process lifecycle coordinator starts and gracefully stops this
	// server. Keeping the provider on the vendor type avoids coupling DI to
	// supervision.
	do.Lazy[*http.Server](func(i do.Injector) (*http.Server, error) {
		cfg, err := do.Invoke[*config.Server](i)
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

// handlerSettings picks the values the handlers read out of the server
// configuration, so handlers never see the whole of it.
func handlerSettings(cfg *config.Server) handlers.Settings {
	return handlers.Settings{
		PublicURL:         cfg.PublicURL,
		SiteURL:           cfg.SiteURL,
		SecureCookies:     cfg.SecureCookies(),
		GuestCookieSecret: cfg.GuestCookieSecret.Reveal(),
		EmailTokenSecret:  cfg.EmailTokenSecret.Reveal(),
		CourseRepoBaseURL: cfg.CourseRepoBaseURL,
		BookBlobKey:       cfg.BookBlobKey,
	}
}
