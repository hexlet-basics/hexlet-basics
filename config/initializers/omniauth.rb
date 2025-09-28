require "omni_auth/strategies/vk"

Rails.application.config.middleware.use OmniAuth::Builder do
  setup = lambda { |env|
    request = ActionDispatch::Request.new(env)
    locale = AppHost.locale_for_url(request.params[:locale])
    strategy = env["omniauth.strategy"]
    strategy.options[:callback_path] = Rails.application.routes.url_helpers.callback_auth_path(strategy.name, locale: locale)
  }
  # provider :developer unless Rails.env.production?
  # provider :github, configus.github.app_id, configus.github.app_secret, scope: 'user:email', setup: setup
  # provider :facebook, configus.facebook.app_id, configus.facebook.app_secret, scope: 'email', setup: setup
  provider :vk, configus.vk.app_id, nil, redirect_url: configus.vk.redirect_url, setup:
end
