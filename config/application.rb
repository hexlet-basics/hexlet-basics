# frozen_string_literal: true

require_relative 'boot'

require 'rails/all'

# Require the gems listed in Gemfile, including any gems
# you've limited to :test, :development, or :production.
Bundler.require(*Rails.groups)

load File.expand_path('../app/lib/configus.rb', __dir__)

module HexletBasics
  class Application < Rails::Application
    # Use the responders controller from the responders gem
    config.app_generators.scaffold_controller :responders_controller

    # Initialize configuration defaults for originally generated Rails version.
    config.load_defaults '6.1'
    config.hexlet_basics = config_for(:hexlet_basics)
    config.autoload_paths << Rails.root.join('lib')
    # Settings in config/environments/* take precedence over those specified here.
    # Application configuration can go into files in config/initializers
    # -- all .rb files in that directory are automatically loaded after loading
    # the framework and any gems in your application.
    config.active_job.queue_adapter = :sidekiq

    config.exceptions_app = routes

    config.i18n.available_locales = %i[en ru]
    config.i18n.default_locale = :ru
    config.generators do |g|
      g.assets false
    end

    routes.default_url_options[:host] = "#{configus.protocol}://#{configus.host}"

    # https://edgeguides.rubyonrails.org/configuring.html#configuring-action-dispatch
    config.action_dispatch.default_headers = config.action_dispatch.default_headers.merge(
      # https://yandex.ru/support/metrica/behavior/click-map.html#iframe
      'X-FRAME-OPTIONS' => 'ALLOW-FROM http://webvisor.com'
    )

    Rails.autoloaders.main.ignore(Rails.root.join('app/packs/*'))
  end
end
