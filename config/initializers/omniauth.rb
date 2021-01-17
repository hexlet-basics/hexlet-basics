# frozen_string_literal: true

Rails.application.config.middleware.use OmniAuth::Builder do
  # FIXME: разобраться с
  # [2021-01-17T13:30:40.607782 #138] ERROR -- omniauth: (developer) Authentication failure! Forbidden: OmniAuth::AuthenticityError, Forbidden
  # [2021-01-17T13:30:40.607493 #138]  WARN -- omniauth: Attack prevented by OmniAuth::AuthenticityTokenProtection
  unless Rails.env.test?
    provider :developer unless Rails.env.production?
    provider :github, configus.github.app_id, configus.github.app_secret, scope: 'user'
  end
end
