// Package config loads runtime configuration from the environment.
package config

import (
	"log/slog"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/samber/oops"
)

// AppEnv is the environment a process runs in. It, and not a Sentry label,
// decides environment policy (which secrets are acceptable, whether mail must
// be deliverable).
type AppEnv string

// The accepted values of APP_ENV.
const (
	Development AppEnv = "development"
	Production  AppEnv = "production"
	Test        AppEnv = "test"
)

// Common is what both processes read: where the data lives, where the site is,
// and how errors are reported. Deploy-critical settings have no code default,
// so a missing variable stops the start instead of running on a dev value.
type Common struct {
	// AppEnv defaults to development so a developer's checkout runs with no
	// setup; a deployment sets APP_ENV=production.
	AppEnv      AppEnv `env:"APP_ENV"      envDefault:"development" validate:"oneof=development production test"`
	DatabaseURL string `env:"DATABASE_URL" validate:"required"`
	// BlobBucketURL selects the gocloud.dev/blob backend (ADR-0005): local disk
	// (fileblob) in dev, s3blob in prod, chosen by URL scheme.
	BlobBucketURL string `env:"BLOB_BUCKET_URL" validate:"required"`
	// SiteURL is the origin of the public site that emailed links open. Unlike
	// AppHost (canonical URLs, always the production host) it points at the
	// local frontend in development, so a developer can follow the link.
	SiteURL string `env:"SITE_URL" validate:"required"`
	// PublicURL is the public origin the API is reached at, used to build
	// absolute asset URLs it serves itself (the `/api/storage/{key}` blob read
	// path). In production it is the site's own origin, since the ingress routes
	// `/api/*` to Go (ADR-0015); in development it reaches Go directly.
	PublicURL string `env:"PUBLIC_URL" validate:"required"`
	// CourseRepoBaseURL is the GitHub org base under which each course's exercises
	// repo lives: the loader clones `<base>/exercises-<slug>.git`.
	CourseRepoBaseURL string `env:"COURSE_REPO_BASE_URL" envDefault:"https://github.com/hexlet-basics"`
	// EmailTokenSecret signs the Magic Link and Password Reset tokens: the server
	// verifies them and the worker mints the emails that carry them. It is kept
	// apart from GuestCookieSecret so rotating one never invalidates the other.
	EmailTokenSecret Secret `env:"EMAIL_TOKEN_SECRET" validate:"required,notdevsecret"`
	// OpenAI powers the AI features (lesson-review summaries, the in-lesson
	// assistant). An empty token disables them. The model must be a reasoning
	// model: the assistant's stream sets a reasoning effort.
	OpenAIAccessToken Secret `env:"OPENAI_ACCESS_TOKEN"`
	OpenAIModel       string `env:"OPENAI_MODEL"        envDefault:"gpt-5.5"`
	// SentryDSN enables exception delivery. An empty DSN intentionally creates a
	// disabled client so local development and tests stay off-process.
	SentryDSN string `env:"SENTRY_DSN"`
	// LogLevelName is the threshold of the application log (ADR-0026). Only the
	// four names are accepted: slog.Level would also parse "INFO+2" and numbers,
	// which no operator needs.
	LogLevelName   string `env:"LOG_LEVEL"                     envDefault:"info" validate:"oneof=debug info warn error"`
	ReleaseVersion string `env:"HEXLET_BASICS_RELEASE_VERSION"`
}

// Server is the configuration of the HTTP process: Common plus what only a
// request needs.
type Server struct {
	Common

	Addr string `env:"ADDR" envDefault:":3001"`
	// AppHost is the canonical public site host the Yandex feed builds its
	// absolute URLs on, mirroring legacy AppHost.canonical. HTTPS is assumed.
	AppHost string `env:"APP_HOST" validate:"required"`
	// BookBlobKey is the bucket key of the free book's PDF. It is uploaded once,
	// out of band, at cutover rather than shipped in the image (~28 MB).
	BookBlobKey string `env:"BOOK_BLOB_KEY" envDefault:"book.pdf"`
	// GitHubWebhookSecret is the shared secret GitHub signs webhook deliveries
	// with (HMAC-SHA256). An empty secret disables the webhook (fail closed).
	GitHubWebhookSecret Secret `env:"GITHUB_WEBHOOK_SECRET"`
	// GuestCookieSecret signs the cookie that carries a visitor's lesson progress
	// until they sign up. Sessions need no secret (ADR-0017). The committed
	// .env.example supplies a development value; deployments provide their own.
	GuestCookieSecret Secret `env:"GUEST_COOKIE_SECRET" validate:"required,notdevsecret"`
	// ExerciseRunner bounds what a learner's submission may do to the host
	// (ADR-0013). Every limit is settable so a course that needs more room is
	// fixable without a release.
	ExerciseRunner ExerciseRunnerConfig `envPrefix:"CHECK_"`
	// OAuth switches on the OAuth authorization server for MCP clients (ADR-0022).
	OAuth OAuthConfig `envPrefix:"OAUTH_" validate:"-"`
	// MCP switches on the MCP endpoint at `/mcp` (ADR-0021).
	MCP MCPConfig `envPrefix:"MCP_" validate:"-"`
}

// Worker is the configuration of the asynchronous process: Common plus what
// only a job or a subscriber needs. It holds no cookie secret and no public
// host: the worker serves no visitor.
type Worker struct {
	Common

	// HealthAddr is where the health-only listener serves /healthz and /readyz
	// (ADR-0023); the worker has no other HTTP surface.
	HealthAddr string `env:"HEALTH_ADDR" envDefault:":8081" validate:"required"`
	// GitHubToken authenticates clones (private repos / higher rate limits). The
	// exercises repos are public, so it is optional; when set it is used as the
	// HTTP basic-auth password (x-access-token).
	GitHubToken Secret `env:"GITHUB_TOKEN"`
	// amoCRM receives lead_created events through a Watermill -> River handler.
	AmoCRMBaseURL   string `env:"AMOCRM_BASE_URL"`
	AmoCRMAuthToken Secret `env:"AMOCRM_AUTH_TOKEN"`
	YMCounter       string `env:"YM_COUNTER"`
	// AuditRetentionMonths is how long Audit Records are kept (ADR-0020). The
	// audit table is partitioned by month and retention drops whole partitions,
	// so a record lives at least this many months and up to one month longer.
	AuditRetentionMonths int `env:"AUDIT_RETENTION_MONTHS" envDefault:"24" validate:"min=1"`
	// Mail configures account email (ADR-0006).
	Mail MailConfig `envPrefix:"MAIL_"`
}

// MCPConfig configures the MCP endpoint. It is inert while Enabled is false,
// which is the default.
type MCPConfig struct {
	// Enabled mounts `/mcp` and its protected resource metadata. Every request
	// to the endpoint needs an OAuth token, so it requires OAUTH_ENABLED: an
	// endpoint no client could sign in to would only answer 401.
	Enabled bool `env:"ENABLED" envDefault:"false"`
	// MaxPageSize caps the records one tool call returns: a larger page size is
	// clamped to it, and the result says so. An agent that asks for everything
	// would otherwise fill its own context with one call.
	MaxPageSize int `env:"MAX_PAGE_SIZE" envDefault:"50" validate:"required,min=1"`
	// MaxResponseBytes caps the structured content of one result. A result over
	// it is cut at a record boundary and carries a pagination hint, never
	// invalid JSON.
	MaxResponseBytes int `env:"MAX_RESPONSE_BYTES" envDefault:"102400" validate:"required,min=1024"`
	// ToolTimeout is how long one tool call may run before its context is
	// cancelled and the agent gets an error.
	ToolTimeout time.Duration `env:"TOOL_TIMEOUT" envDefault:"30s" validate:"required,min=1s"`
	// RateWindow is the period every rate limit below counts over.
	RateWindow time.Duration `env:"RATE_WINDOW" envDefault:"1m" validate:"required,min=1s"`
	// The rate limits are calls per RateWindow, per class of tool, for one user
	// and for one Connected App (an OAuth client, shared by all its users).
	// Catalogue reads are public data; reads of personal data and writes are
	// stricter, and writes the strictest. They are in memory, per replica.
	RateUserCatalogue int `env:"RATE_USER_CATALOGUE" envDefault:"60" validate:"required,min=1"`
	RateUserPersonal  int `env:"RATE_USER_PERSONAL"  envDefault:"30" validate:"required,min=1"`
	RateUserWrite     int `env:"RATE_USER_WRITE"     envDefault:"10" validate:"required,min=1"`
	// The Connected App limits are higher than the per-user ones, since an app
	// speaks for many users; they stop one misbehaving app, not one user.
	RateClientCatalogue int `env:"RATE_CLIENT_CATALOGUE" envDefault:"600" validate:"required,min=1"`
	RateClientPersonal  int `env:"RATE_CLIENT_PERSONAL"  envDefault:"300" validate:"required,min=1"`
	RateClientWrite     int `env:"RATE_CLIENT_WRITE"     envDefault:"100" validate:"required,min=1"`
}

// OAuthConfig configures the in-process OAuth authorization server. Everything
// in it is inert while Enabled is false, which is the default: the code ships
// before the feature is exposed.
type OAuthConfig struct {
	// Enabled mounts `/oauth/*` and the well-known metadata documents.
	Enabled bool `env:"ENABLED" envDefault:"false"`
	// Issuer is the bare canonical origin (scheme and host, no path) the server
	// speaks as. It is never derived from the Host header: a forged header must
	// not change what a token claims, and the value is baked into every
	// registration and token, so it is expensive to change later.
	Issuer string `env:"ISSUER" validate:"required,origin"`
	// EncryptionKey is the secret the access tokens are encrypted under and the
	// ID-token signing key is derived from. The database stores only token ids,
	// so this key is what stands between a database leak and a usable token:
	// rotating it invalidates every access token. At least 32 characters.
	EncryptionKey Secret `env:"ENCRYPTION_KEY" validate:"required,min=32,notdevsecret"`
	// RegistrationLimit is how many client registrations (Dynamic Client
	// Registration is anonymous) one address may attempt per RegistrationWindow.
	RegistrationLimit int `env:"REGISTRATION_LIMIT" envDefault:"10" validate:"required,min=1"`
	// RegistrationWindow is the period RegistrationLimit counts over.
	RegistrationWindow time.Duration `env:"REGISTRATION_WINDOW" envDefault:"1h" validate:"required,min=1s"`
	// AccessTokenTTL is how long an access token lives. Short, because a bearer
	// token cannot be recalled once stolen except by revoking its row.
	AccessTokenTTL time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"1h" validate:"required,min=1m"`
	// RefreshTokenTTL is how long a refresh token (and so a connection that is not
	// used) lives. Each refresh issues a new token with a fresh full lifetime.
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"720h"`
}

// ExerciseRunnerConfig is the isolation budget one submission runs under. The
// zero values are never used: Load fills them from the hardened defaults.
type ExerciseRunnerConfig struct {
	// TimeoutSeconds is the budget the in-container `timeout` enforces, and is
	// therefore what classifies a submission as an infinite loop.
	TimeoutSeconds int `env:"TIMEOUT_SECONDS" envDefault:"6"`
	// GraceSeconds is added to it for the outer deadline, which fires only when
	// the container never reported at all.
	GraceSeconds int   `env:"GRACE_SECONDS" envDefault:"10"`
	MemoryBytes  int64 `env:"MEMORY_BYTES"  envDefault:"536870912"`
	// SwapBytes equals MemoryBytes by default, which means no swap. Legacy set
	// -1 — unlimited swap, which is the absence of a limit rather than one.
	SwapBytes int64 `env:"SWAP_BYTES" envDefault:"536870912"`
	// PidsLimit stops a fork bomb, which neither a memory cap nor an isolated
	// network does. Legacy set none.
	PidsLimit int64 `env:"PIDS_LIMIT" envDefault:"256"`
	// NanoCPUs caps CPU in billionths of a core.
	NanoCPUs int64 `env:"NANO_CPUS" envDefault:"1000000000"`
	// ReadonlyRootfs and RunAsUser are the two limits that can break a currently
	// green exercise — compiled languages write build artifacts beside their
	// sources — so they stay off until a run of every course's reference
	// solutions says otherwise.
	ReadonlyRootfs bool `env:"READONLY_ROOTFS" envDefault:"false"`
	// RunAsUser is spelled out rather than named USER: the prefixed name is what
	// a deployment sets, but a bare USER would sit one refactor away from
	// picking up the ambient account name of whoever started the process.
	RunAsUser      string `env:"RUN_AS_USER"`
	MaxOutputBytes int    `env:"MAX_OUTPUT_BYTES" envDefault:"65536"`
	// Concurrency bounds simultaneous runs: the real ceiling is the Docker
	// daemon's capacity, not this process's.
	Concurrency int `env:"CONCURRENCY" envDefault:"8"`
	// ImageTag is the moving tag a course publishes its latest build under. It
	// is pulled once per course version and frozen under a pinned reference, so
	// a rebuild cannot change how an already-promoted version is graded. Legacy
	// used `release` in production and `latest` everywhere else.
	ImageTag string `env:"IMAGE_TAG" envDefault:"latest"`
}

// MailConfig selects where account email goes. Postbox is used only when its
// access key is set, so development and CI log email instead of sending it.
type MailConfig struct {
	// From is the sender address; the legacy app sent from the same one, whose
	// domain Postbox must have verified.
	From string `env:"FROM" envDefault:"support@hexlet.io"`
	// PostboxEndpoint and PostboxRegion point the SES v2 client at Yandex Cloud
	// Postbox rather than AWS.
	PostboxEndpoint        string `env:"POSTBOX_ENDPOINT"          envDefault:"https://postbox.cloud.yandex.net"`
	PostboxRegion          string `env:"POSTBOX_REGION"            envDefault:"ru-central1"`
	PostboxAccessKeyID     string `env:"POSTBOX_ACCESS_KEY_ID"`
	PostboxSecretAccessKey Secret `env:"POSTBOX_SECRET_ACCESS_KEY"`
}

// PostboxEnabled reports whether real delivery is configured.
func (c MailConfig) PostboxEnabled() bool { return c.PostboxAccessKeyID != "" }

// LogLevel is the slog threshold LOG_LEVEL names; validation has already
// restricted it to the four known names.
func (c *Common) LogLevel() slog.Level {
	var level slog.Level

	err := level.UnmarshalText([]byte(c.LogLevelName))
	if err != nil {
		return slog.LevelInfo
	}

	return level
}

// SecureCookies reports whether the site is served over https, which is when
// cookies are marked Secure.
func (c *Common) SecureCookies() bool { return strings.HasPrefix(c.PublicURL, "https://") }

// LoadServer reads the HTTP process's configuration from the environment and
// validates it. The environment is the only source: a .env file is the
// caller's business (LoadDotenv), so production and tests never read one.
// Every violated rule is reported at once, named by variable.
func LoadServer() (*Server, error) {
	cfg := new(Server)

	err := env.Parse(cfg)
	if err != nil {
		return nil, oops.Wrapf(err, "parse server config from environment")
	}

	err = check(cfg, cfg.AppEnv)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadWorker reads the asynchronous process's configuration. It parses only the
// worker's own settings, so a deployment need not give it a server-only secret.
func LoadWorker() (*Worker, error) {
	cfg := new(Worker)

	err := env.Parse(cfg)
	if err != nil {
		return nil, oops.Wrapf(err, "parse worker config from environment")
	}

	err = check(cfg, cfg.AppEnv)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// crossChecks holds the rules that span settings, which a tag on one field
// cannot say.
func (c *Server) crossChecks() []string {
	var problems []string
	if c.OAuth.Enabled {
		problems = append(problems, violations(c.OAuth, "OAUTH_", c.AppEnv)...)
	}

	if c.MCP.Enabled {
		problems = append(problems, violations(c.MCP, "MCP_", c.AppEnv)...)
	}

	if c.MCP.Enabled && !c.OAuth.Enabled {
		problems = append(problems, "MCP_ENABLED requires OAUTH_ENABLED: the MCP endpoint accepts only OAuth tokens")
	}
	// A refresh token that does not outlive its access token could never refresh.
	if c.OAuth.Enabled && c.OAuth.RefreshTokenTTL <= c.OAuth.AccessTokenTTL {
		problems = append(problems, "OAUTH_REFRESH_TOKEN_TTL must be longer than OAUTH_ACCESS_TOKEN_TTL")
	}

	if c.AppEnv == Production && c.OAuth.Enabled && !strings.HasPrefix(c.OAuth.Issuer, "https://") {
		problems = append(problems, "OAUTH_ISSUER must be an https origin in production")
	}

	return problems
}

// crossChecks: in production the worker must be able to deliver mail.
func (c *Worker) crossChecks() []string {
	if c.AppEnv == Production && !c.Mail.PostboxEnabled() {
		return []string{"MAIL_POSTBOX_ACCESS_KEY_ID is required in production"}
	}

	return nil
}
