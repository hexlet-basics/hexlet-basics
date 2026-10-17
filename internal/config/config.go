// Package config loads runtime configuration from the environment.
package config

import (
	"net/url"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/samber/oops"
)

const (
	developmentGuestCookieSecret = "dev-insecure-guest-cookie-secret-change-me"
	developmentEmailTokenSecret  = "dev-insecure-email-token-secret-change-me"
)

// Config holds the runtime configuration for the server. Fields are populated
// from environment variables (12-factor); defaults keep local dev zero-config.
type Config struct {
	Addr        string `env:"ADDR" envDefault:":3001"`
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@127.0.0.1:54330/code_basics_development"`
	// BlobBucketURL selects the gocloud.dev/blob backend (ADR-0005): local disk
	// (fileblob) in dev, s3blob in prod, chosen by URL scheme. The dev default
	// writes under ./storage (created on first open) so uploads survive restarts.
	BlobBucketURL string `env:"BLOB_BUCKET_URL" envDefault:"file://./storage?create_dir=true"`
	// BookBlobKey is the bucket key of the free book's PDF (legacy
	// public/book.pdf). It is uploaded once, out of band, at cutover rather than
	// shipped in the image: at ~28 MB it has no place in the repo or the build.
	BookBlobKey string `env:"BOOK_BLOB_KEY" envDefault:"book.pdf"`
	// AppHost is the canonical public site host the Yandex feed builds its
	// absolute URLs on, mirroring legacy AppHost.canonical. HTTPS is assumed in
	// prod; the default matches the legacy default.
	AppHost string `env:"APP_HOST" envDefault:"code-basics.com"`
	// SiteURL is the origin of the public site that emailed links open. Unlike
	// AppHost (canonical URLs, always the production host) it points at the
	// local frontend in development, so a developer can follow the link.
	SiteURL string `env:"SITE_URL" envDefault:"https://code-basics.com"`
	// PublicURL is the public origin this server is reached at, used to build
	// absolute asset URLs it serves itself (the `/api/storage/{key}` blob read
	// path). In production it is the site's own origin, since the ingress routes
	// `/api/*` to Go (ADR-0015); the dev default reaches Go directly. Separate
	// from AppHost, which is always the canonical production host.
	PublicURL string `env:"PUBLIC_URL" envDefault:"http://localhost:3001"`
	// CourseRepoBaseURL is the GitHub org base under which each course's exercises
	// repo lives: the loader clones `<base>/exercises-<slug>.git`. The default
	// matches the legacy hard-coded convention (apiconv.repositoryURL).
	CourseRepoBaseURL string `env:"COURSE_REPO_BASE_URL" envDefault:"https://github.com/hexlet-basics"`
	// GitHubToken authenticates clones (private repos / higher rate limits). The
	// exercises repos are public, so it is optional; when set it is used as the
	// HTTP basic-auth password (x-access-token).
	GitHubToken string `env:"GITHUB_TOKEN"`
	// GitHubWebhookSecret is the shared secret GitHub signs webhook deliveries
	// with (HMAC-SHA256). The webhook route rejects any delivery whose signature
	// doesn't verify; an empty secret disables the webhook (fail closed).
	GitHubWebhookSecret string `env:"GITHUB_WEBHOOK_SECRET"`
	// amoCRM receives lead_created events through a Watermill -> River handler.
	AmoCRMBaseURL   string `env:"AMOCRM_BASE_URL"`
	AmoCRMAuthToken string `env:"AMOCRM_AUTH_TOKEN"`
	YMCounter       string `env:"YM_COUNTER"`
	// OpenAI powers the AI features (lesson-review summaries, the in-lesson
	// assistant). The env name matches the legacy Rails deployment so the
	// secret carries over unchanged; an empty token disables the AI workers.
	// The model must be a reasoning model: the assistant's stream sets a
	// reasoning effort, which older chat models refuse.
	OpenAIAccessToken string `env:"OPENAI_ACCESS_TOKEN"`
	OpenAIModel       string `env:"OPENAI_MODEL" envDefault:"gpt-5.5"`
	// SentryDSN enables exception delivery. An empty DSN intentionally creates a
	// disabled client so local development and tests stay off-process.
	SentryDSN         string `env:"SENTRY_DSN"`
	SentryEnvironment string `env:"SENTRY_ENVIRONMENT" envDefault:"development"`
	ReleaseVersion    string `env:"HEXLET_BASICS_RELEASE_VERSION"`
	// GuestCookieSecret signs the cookie that carries a visitor's lesson progress
	// until they sign up. Sessions need no secret (ADR-0017): the cookie is an
	// opaque token, not a signed claim. The committed .env.example supplies a
	// development value, while deployments must provide a non-empty secret
	// explicitly.
	GuestCookieSecret string `env:"GUEST_COOKIE_SECRET,required,notEmpty"`
	// EmailTokenSecret signs the Magic Link and Password Reset tokens. It is kept
	// apart from GuestCookieSecret so rotating one never invalidates the other: a
	// leaked email token secret must not force every learner to sign in again.
	EmailTokenSecret string `env:"EMAIL_TOKEN_SECRET,required,notEmpty"`
	// AuditRetentionMonths is how long Audit Records are kept (ADR-0017). The
	// audit table is partitioned by month and retention drops whole partitions,
	// so a record lives at least this many months and up to one month longer.
	AuditRetentionMonths int `env:"AUDIT_RETENTION_MONTHS" envDefault:"24"`
	// Mail configures account email (ADR-0006).
	Mail MailConfig `envPrefix:"MAIL_"`
	// ExerciseRunner bounds what a learner's submission may do to the host
	// (ADR-0013). Every limit is settable because the defaults tighten what the
	// legacy runner allowed, and a course that turns out to need more room must
	// be fixable without a release.
	ExerciseRunner ExerciseRunnerConfig `envPrefix:"CHECK_"`
	// OAuth switches on the OAuth authorization server for MCP clients (ADR-0019).
	OAuth OAuthConfig `envPrefix:"OAUTH_"`
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
	Issuer string `env:"ISSUER"`
	// EncryptionKey is the secret the access tokens are encrypted under and the
	// ID-token signing key is derived from. The database stores only token ids,
	// so this key is what stands between a database leak and a usable token:
	// rotating it invalidates every access token. At least 32 characters.
	EncryptionKey string `env:"ENCRYPTION_KEY"`
}

// OAuthMinKeyLength is the shortest accepted EncryptionKey.
const OAuthMinKeyLength = 32

// ExerciseRunnerConfig is the isolation budget one submission runs under. The
// zero values are never used: Load fills them from the hardened defaults.
type ExerciseRunnerConfig struct {
	// TimeoutSeconds is the budget the in-container `timeout` enforces, and is
	// therefore what classifies a submission as an infinite loop.
	TimeoutSeconds int `env:"TIMEOUT_SECONDS" envDefault:"6"`
	// GraceSeconds is added to it for the outer deadline, which fires only when
	// the container never reported at all.
	GraceSeconds int   `env:"GRACE_SECONDS" envDefault:"10"`
	MemoryBytes  int64 `env:"MEMORY_BYTES" envDefault:"536870912"`
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
	PostboxEndpoint        string `env:"POSTBOX_ENDPOINT" envDefault:"https://postbox.cloud.yandex.net"`
	PostboxRegion          string `env:"POSTBOX_REGION" envDefault:"ru-central1"`
	PostboxAccessKeyID     string `env:"POSTBOX_ACCESS_KEY_ID"`
	PostboxSecretAccessKey string `env:"POSTBOX_SECRET_ACCESS_KEY"`
}

// PostboxEnabled reports whether real delivery is configured.
func (c MailConfig) PostboxEnabled() bool { return c.PostboxAccessKeyID != "" }

// Load reads configuration from environment variables, applying defaults.
// A local .env is loaded first if present; it only fills variables that are not
// already set in the real environment (godotenv never overrides), so exported
// vars and 12-factor deploys keep precedence over the dev-convenience file.
// SecureCookies reports whether the site is served over https, which is when
// cookies are marked Secure.
func (c *Config) SecureCookies() bool { return strings.HasPrefix(c.PublicURL, "https://") }

func Load() (*Config, error) {
	// Ignore a missing file: .env is a local convenience, absent in prod/CI.
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, oops.Wrapf(err, "parse config from environment")
	}
	if err := validateOAuth(cfg.OAuth); err != nil {
		return nil, oops.Wrapf(err, "validate OAuth config")
	}
	if cfg.AuditRetentionMonths < 1 {
		// Zero would drop the month in progress; the job would refuse it on
		// every run, so the process refuses to start instead.
		return nil, oops.Errorf("AUDIT_RETENTION_MONTHS must be at least 1, got %d", cfg.AuditRetentionMonths)
	}
	if err := validateProduction(cfg); err != nil {
		return nil, oops.Wrapf(err, "validate production config")
	}
	return cfg, nil
}

// validateProduction keeps environment-specific safety policy separate from
// env decoding. The development secret is committed for local convenience and
// therefore must never be accepted as a session-signing key in production.
func validateProduction(cfg *Config) error {
	if cfg.SentryEnvironment != "production" {
		return nil
	}
	if cfg.GuestCookieSecret == developmentGuestCookieSecret {
		return oops.Errorf("GUEST_COOKIE_SECRET must differ from the public development value")
	}
	if cfg.EmailTokenSecret == developmentEmailTokenSecret {
		return oops.Errorf("EMAIL_TOKEN_SECRET must differ from the public development value")
	}
	if !cfg.Mail.PostboxEnabled() {
		return oops.Errorf("MAIL_POSTBOX_ACCESS_KEY_ID is required in production")
	}
	if cfg.OAuth.Enabled && !strings.HasPrefix(cfg.OAuth.Issuer, "https://") {
		return oops.Errorf("OAUTH_ISSUER must be an https origin in production")
	}
	return nil
}

// validateOAuth checks the settings the OAuth server cannot run without. It
// applies only when the flag is on, so the default configuration needs none of
// them.
func validateOAuth(cfg OAuthConfig) error {
	if !cfg.Enabled {
		return nil
	}
	issuer, err := url.Parse(cfg.Issuer)
	if err != nil || issuer.Host == "" || (issuer.Scheme != "https" && issuer.Scheme != "http") {
		return oops.Errorf("OAUTH_ISSUER must be an origin such as https://code-basics.com")
	}
	if issuer.Path != "" || issuer.RawQuery != "" || issuer.Fragment != "" {
		return oops.Errorf("OAUTH_ISSUER must be a bare origin without a path: discovery lives at the root")
	}
	if len(cfg.EncryptionKey) < OAuthMinKeyLength {
		return oops.Errorf("OAUTH_ENCRYPTION_KEY must be at least %d characters", OAuthMinKeyLength)
	}
	return nil
}
