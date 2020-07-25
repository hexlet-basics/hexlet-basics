# frozen_string_literal: true

Rails.application.config.middleware.use OmniAuth::Builder do
  provider :github, ENV.fetch('GITHUB_APP_ID'), ENV.fetch('GITHUB_APP_SECRET'), scope: 'user'
end
