source "https://rubygems.org"
# ruby file: ".ruby-version"

# Bundle edge Rails instead: gem "rails", github: "rails/rails", branch: "main"
gem "rails"
gem "puma"

# Audits gems for known security defects (use config/bundler-audit.yml to ignore issues)
gem "bundler-audit", require: false

# Build JSON APIs with ease [https://github.com/rails/jbuilder]
# gem "jbuilder"

# Windows does not include zoneinfo files, so bundle the tzinfo-data gem
gem "tzinfo-data", platforms: %i[ windows jruby ]
gem "alba"
gem "typelizer"
gem "ruby_llm"

# Use the database-backed adapters for Rails.cache, Active Job, and Action Cable
gem "solid_cache"
gem "solid_queue"
gem "solid_cable"

# Reduces boot times through caching; required in config/boot.rb
gem "bootsnap", require: false

# Deploy this application anywhere as a Docker container [https://kamal-deploy.org]
# gem "kamal", require: false

# Add HTTP asset caching/compression and X-Sendfile acceleration to Puma [https://github.com/basecamp/thruster/]
gem "thruster", require: false

gem "inertia_rails"
gem "vite_rails"
gem "foreman"
gem "ostruct"

# gem 'aws-sdk-s3'

# gem "configus"
gem "geocoder"
gem "image_processing"
gem "ruby-vips", "~> 2.0", require: false
# gem 'responders'

# gem 'http_accept_language'
gem "pagy"
gem "omniauth-facebook"
gem "omniauth-github"
gem "omniauth-rails_csrf_protection"
gem "phonelib"
gem "ransack"
# gem 'redis'
gem "term-ansicolor"
gem "valid_email2"
gem "webauthn"

# gem 'bootstrap'
# gem 'jquery-rails'

# gem 'terser'

# gem 'jsbundling-rails'

gem "stackprof"
gem "sentry-ruby"
gem "sentry-rails"

# Condense Rails' multi-line request logs into one structured line per request
gem "lograge"

gem "amocrm"

# gem 'stimulus-rails'
# gem 'turbo-rails'

gem "js-routes"

gem "sorbet"
gem "sorbet-schema"
gem "sorbet-result"

group :development, :test, :staging do
  gem "debug", platforms: %i[ mri windows ], require: "debug/prelude"
  gem "dotenv"

  # Static analysis for security vulnerabilities [https://brakemanscanner.org/]
  gem "brakeman", require: false

  # Find missing/unused i18n keys across locale files
  gem "i18n-tasks", require: false

  # Omakase Ruby styling [https://github.com/rails/rubocop-rails-omakase/]
  gem "rubocop-rails-omakase", require: false
  gem "rubocop-sorbet", require: false
  gem "rubocop-minitest", require: false
  gem "rubocop-capybara", require: false
  gem "rubocop-factory_bot", require: false
  gem "rubocop-rake", require: false

  gem "factory_bot_rails"
  gem "faker"

  gem "database_consistency", require: false

  gem "tapioca"
  gem "boba", require: false
  # gem "bullet"
end

group :development, :staging do
  # gem 'annotate'
  gem "annotaterb"
  # gem 'i18n-debug'
  # gem 'listen'
  # gem 'derailed_benchmarks'
  gem "web-console"
  gem "ruby-lsp", require: false
  gem "ruby-lsp-rails", require: false
end

group :test do
  gem "capybara"
  gem "cuprite"
  gem "minitest-power_assert"
  gem "simplecov"
  # gem 'webdrivers'
  # gem 'whiny_validation'
end

group :production do
  gem "pg"
end

gem "aasm"
gem "acts-as-taggable-on"
gem "webmock"
gem "aws-sdk-s3"
gem "bcrypt"
# gem "enumerize"
gem "gon"
gem "googleauth"
gem "rails-i18n"
# gem "redcarpet"
gem "commonmarker"
# gem "simple_form"

# gem 'sitemap_generator'
# gem 'slim-rails'

gem "counter_culture"
# NOTE: было только в dev/test/staging — из-за этого N+1 маскировались локально, но
# жили в проде. Включён во всех окружениях, чтобы авто-eager-load ассоциаций работал и на проде.
gem "goldiloader"
gem "ahoy_matey"
# gem "country_select"
# gem 'sorbet-rails' # Temp disable sorbet for memory optimization

gem "strong_migrations"
gem "browser"
gem "browserslist"
#
gem "meta-tags"
#
# gem 'browserslist_useragent'

gem "responders"

gem "configus"
gem "rails_event_store"

gem "pundit"
gem "rack-attack"
# gem "faraday"
# gem "httpx"
