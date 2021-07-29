# frozen_string_literal: true

require 'rollbar/rails'

Rollbar.configure do |config|
  config.enabled = Rails.env.production?
  config.js_enabled = false # Rails.env.production?

  config.js_options = {
    accessToken: ENV['ROLLBAR_CLIENT_ACCESS_TOKEN'],
    captureUncaught: true,
    captureUnhandledRejections: true,
    hostSafeList: [configus.host],
    payload: {
      client: {
        # source_map_enabled: true,
        code_version: ENV['CODE_BASICS_VERSION'],
        guess_uncaught_frames: true
      },
      environment: Rails.env
    }
  }
  # Without configuration, Rollbar is enabled in all environments.
  # To disable in specific environments, set config.enabled=false.

  config.access_token = ENV['ROLLBAR_ACCESS_TOKEN']

  config.person_username_method = 'username'

  # By default, Rollbar will try to call the `current_user` controller method
  # to fetch the logged-in user object, and then call that object's `id`,
  # `username`, and `email` methods to fetch those properties. To customize:
  # config.person_method = "my_current_user"
  # config.person_id_method = "my_id"
  # config.person_email_method = "my_email"

  # If you want to attach custom data to all exception and message reports,
  # provide a lambda like the following. It should return a hash.
  # config.custom_data_method = lambda { {:some_key => "some_value" } }

  # Add exception class names to the exception_level_filters hash to
  # change the level that exception is reported at. Note that if an exception
  # has already been reported and logged the level will need to be changed
  # via the rollbar interface.
  # Valid levels: 'critical', 'error', 'warning', 'info', 'debug', 'ignore'
  # 'ignore' will cause the exception to not be reported at all.
  config.exception_level_filters.merge!(
    'ActionController::RoutingError' => 'ignore',
    'ActionController::UnknownFormat' => 'ignore',
    'Pundit::NotAuthorizedError' => 'ignore',
    'Mime::Type::InvalidMimeType' => 'ignore'
  )

  #
  # You can also specify a callable, which will be called with the exception instance.
  # config.exception_level_filters.merge!('MyCriticalException' => lambda { |e| 'critical' })

  # Enable asynchronous reporting (uses girl_friday or Threading if girl_friday
  # is not installed)
  config.use_async = true
  # config.failover_handlers = [Rollbar::Delay::GirlFriday, Rollbar::Delay::Thread]
  # Supply your own async handler:
  # config.async_handler = Proc.new { |payload|
  #  Thread.new { Rollbar.process_payload_safely(payload) }
  # }

  # Enable asynchronous reporting (using sucker_punch)
  # config.use_sucker_punch

  # Enable delayed reporting (using Sidekiq)
  # config.use_sidekiq
  # You can supply custom Sidekiq options:
  # config.use_sidekiq 'queue' => 'my_queue'
end
