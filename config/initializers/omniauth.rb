# frozen_string_literal: true

Rails.application.config.middleware.use OmniAuth::Builder do
  setup = lambda {
    lambda { |env|
      request = ActionDispatch::Request.new(env)
      locale = AppHost.locale_for_url(request.params[:locale])
      strategy = env['omniauth.strategy']
      strategy.options[:callback_path] = Rails.application.routes.url_helpers.callback_auth_path(strategy.name, locale: locale)
    }
  }

  provider :developer unless Rails.env.production?
  provider :github, configus.github.app_id, configus.github.app_secret, scope: 'user:email', setup: setup.call
  provider :facebook, configus.facebook.app_id, configus.facebook.app_secret, scope: 'email', setup: setup.call
end
