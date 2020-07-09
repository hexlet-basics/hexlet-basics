# frozen_string_literal: true

require_relative 'boot'

require 'rails/all'

# Require the gems listed in Gemfile, including any gems
# you've limited to :test, :development, or :production.
Bundler.require(*Rails.groups)

module HexletBasics
  class Application < Rails::Application
    # Initialize configuration defaults for originally generated Rails version.
    config.load_defaults 6.0
    config.hexlet_basics = config_for(:hexlet_basics)
    config.autoload_paths << Rails.root.join('lib')
    # Settings in config/environments/* take precedence over those specified here.
    # Application configuration can go into files in config/initializers
    # -- all .rb files in that directory are automatically loaded after loading
    # the framework and any gems in your application.
    config.active_job.queue_adapter = :sidekiq

    config.i18n.available_locales = [:en, :ru]
    config.i18n.default_locale = :ru
    config.generators do |g|
      g.assets false
    end
  end
end
