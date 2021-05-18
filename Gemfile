# frozen_string_literal: true

source 'https://rubygems.org'
git_source(:github) { |repo| "https://github.com/#{repo}.git" }

gem 'configus'
gem 'puma', '~> 5.3'
gem 'rails', '~> 6.1'
gem 'webpacker', github: 'rails/webpacker'

gem 'bootsnap', require: false
gem 'http_accept_language'
gem 'jbuilder', '~> 2.10'
gem 'kaminari'
gem 'omniauth-github'
gem 'omniauth-rails_csrf_protection'
gem 'ransack'
gem 'rollbar'
# gem 'sorbet-runtime'
gem 'term-ansicolor'
gem 'valid_email2'

group :development, :test do
  gem 'brakeman'
  gem 'byebug', platforms: %i[mri mingw x64_mingw]
  gem 'dotenv-rails'
  gem 'factory_bot_rails'
  gem 'faker'
  gem 'slim_lint', require: false
  # gem 'sorbet'
  gem 'sqlite3', '~> 1.4'
end

group :development do
  gem 'i18n-debug'
  gem 'listen', '~> 3.3'
  gem 'rubocop'
  gem 'rubocop-performance'
  gem 'rubocop-rails'
  gem 'rubocop-rake'
  gem 'rubocop-rspec'
  gem 'solargraph', github: 'castwide/solargraph'
  gem 'spring'
  gem 'web-console', '>= 4.1.0'
end

group :test do
  gem 'capybara', '>= 2.15'
  gem 'minitest-power_assert'
  gem 'selenium-webdriver'
  gem 'webdrivers'
end

group :production do
  gem 'pg'
end

gem 'aasm'
gem 'active_form_model'
gem 'bcrypt'
gem 'dry-auto_inject'
gem 'dry-container'
gem 'enumerize'
gem 'gon'
gem 'js-routes'
gem 'rails-i18n'
gem 'redcarpet', '~> 3.5'
gem 'reform-rails'
gem 'sidekiq'
gem 'simple_form'
gem 'slim-rails'
gem 'tzinfo-data', platforms: %i[mingw mswin x64_mingw jruby]

# gem 'sorbet-rails', '~> 0.7.2'
