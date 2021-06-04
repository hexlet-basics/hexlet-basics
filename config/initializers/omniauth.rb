# frozen_string_literal: true

Rails.application.config.middleware.use OmniAuth::Builder do
  provider :developer unless Rails.env.production?
  provider :github, configus.github.app_id, configus.github.app_secret, scope: 'user'
  provider :facebook, configus.facebook.app_id, configus.facebook.app_secret, scope: 'user'
end
