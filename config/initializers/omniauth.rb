# frozen_string_literal: true

Rails.application.config.middleware.use OmniAuth::Builder do
  provider :github, configus.github.app_id, configus.github.app_secret, scope: 'user'
end
