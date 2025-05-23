require_relative "boot"

ENV["DISABLE_TYPELIZER"] = !ENV.fetch("ENABLE_TYPELIZER", "0").to_i.positive? ? "1" : "0"

# require "rails/all"
require "rails"
# Pick the frameworks you want:
require "active_model/railtie"
require "active_job/railtie"
require "active_record/railtie"
require "active_storage/engine"
require "action_controller/railtie"
require "action_mailer/railtie"
# require "action_mailbox/engine"
# require "action_text/engine"
require "action_view/railtie"
require "action_cable/engine"
# require "sprockets/railtie"
require "rails/test_unit/railtie"

# Require the gems listed in Gemfile, including any gems
# you've limited to :test, :development, or :production.
Bundler.require(*Rails.groups)

load File.expand_path("../app/lib/configus.rb", __dir__)
load File.expand_path("../app/lib/app_host.rb", __dir__)

module HexletBasics
  class Application < Rails::Application
    def self.assets_precompiling?
      Rake.application.top_level_tasks.include?("assets:precompile")
    end

    # Initialize configuration defaults for originally generated Rails version.
    config.load_defaults 8.0

    config.generators do |g|
      g.helper false
      g.stylesheets false
      g.assets false
      g.template_engine false
    end

    # Please, add to the `ignore` list any other `lib` subdirectories that do
    # not contain `.rb` files, or that should not be reloaded or eager loaded.
    # Common ones are `templates`, `generators`, or `middleware`, for example.

    config.autoload_lib(ignore: %w[assets tasks])
    config.exceptions_app = routes
    config.active_storage.service_urls_expire_in = 1.week

    config.active_record.default_column_serializer = YAML

    routes.default_url_options = { host: AppHost.canonical }

    # config.solid_queue.logger = ActiveSupport::Logger.new(STDOUT)
    # Configuration for the application, engines, and railties goes here.
    #
    # These settings can be overridden in specific environments using the files
    # in config/environments, which are processed later.
    #
    # config.time_zone = "Central Time (US & Canada)"
    # config.eager_load_paths << Rails.root.join("extras")
    # config.middleware.use ActiveStorage::SetCurrent

    config.i18n.available_locales = %i[en ru es]
    config.i18n.default_locale = :en
  end
end
