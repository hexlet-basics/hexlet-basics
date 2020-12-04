# frozen_string_literal: true

Configus.build Rails.env do
  env :production do
    github do
      app_id ENV['GITHUB_APP_ID']
      app_secret ENV['GITHUB_APP_SECRET']
    end
  end

  env :development, parent: :production do
  end

  env :test, parent: :development do
    github do
      app_id
      app_secret
    end
  end
end
