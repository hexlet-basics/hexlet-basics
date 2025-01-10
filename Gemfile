source "https://rubygems.org"

# Bundle edge Rails instead: gem "rails", github: "rails/rails", branch: "main"
gem "rails", "~> 8.0.1"
# Use sqlite3 as the database for Active Record
gem "sqlite3", ">= 2.1"
# Use the Puma web server [https://github.com/puma/puma]
gem "puma", ">= 5.0"
# Build JSON APIs with ease [https://github.com/rails/jbuilder]
# gem "jbuilder"

# Windows does not include zoneinfo files, so bundle the tzinfo-data gem
gem "tzinfo-data", platforms: %i[ windows jruby ]
gem "alba"
gem "typelizer"

# Use the database-backed adapters for Rails.cache, Active Job, and Action Cable
gem "solid_cache"
gem "solid_queue"
# gem "solid_cable"

# Reduces boot times through caching; required in config/boot.rb
gem "bootsnap", require: false

# Deploy this application anywhere as a Docker container [https://kamal-deploy.org]
gem "kamal", require: false

# Add HTTP asset caching/compression and X-Sendfile acceleration to Puma [https://github.com/basecamp/thruster/]
gem "thruster", require: false

gem "inertia_rails"
gem "vite_rails"
gem "foreman"

# gem 'aws-sdk-s3'

# gem "configus"
gem "geocoder"
gem "image_processing"
# gem 'responders'

# gem 'http_accept_language'
gem "pagy"
# gem 'omniauth-facebook'
# gem 'omniauth-github'
# gem 'omniauth-rails_csrf_protection'
gem "ransack"
# gem 'redis'
gem "term-ansicolor"
# gem 'valid_email2'

# gem 'bootstrap'
# gem 'jquery-rails'

# gem 'terser'

# gem 'jsbundling-rails'

# gem 'sentry-rails'
# gem 'sentry-ruby'

# gem 'stimulus-rails'
# gem 'turbo-rails'

gem "js-routes"
gem "i18n-js"

group :development, :test, :staging do
  gem "debug", platforms: %i[ mri windows ], require: "debug/prelude"

  # Static analysis for security vulnerabilities [https://brakemanscanner.org/]
  gem "brakeman", require: false

  # Omakase Ruby styling [https://github.com/rails/rubocop-rails-omakase/]
  gem "rubocop-rails-omakase", require: false

  # gem 'factory_bot_rails'
  gem "faker"

  # https://github.com/Shopify/tapioca/pull/2129/files
  gem "tapioca", require: false
  gem "bullet"
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
  gem "selenium-webdriver"
  # gem 'minitest-power_assert'
  # gem 'webdrivers'
  # gem 'whiny_validation'
end

group :production do
  gem "pg"
end

gem "aasm"
gem "active_form_model", github: "Hexlet/active_form_model"
gem "bcrypt"
gem "dry-auto_inject"
gem "dry-container"
gem "enumerize"
gem "gon"
gem "googleauth"
gem "rails-i18n"
gem "redcarpet"
gem "simple_form"

# gem 'sitemap_generator'
# gem 'slim-rails'

gem "counter_culture"
gem "country_select"
# gem 'sorbet-rails' # Temp disable sorbet for memory optimization

gem "browser"
#
# gem 'meta-tags'
#
# gem 'browserslist_useragent'

gem "responders", "~> 3.1"

gem "configus", "~> 0.0.5"
