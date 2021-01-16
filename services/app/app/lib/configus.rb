# frozen_string_literal: true

Configus.build Rails.env do
  env :production do
    protocol :https
    host 'code-basics.com'

    github do
      app_id ENV['GITHUB_CLIENT_ID']
      app_secret ENV['GITHUB_CLIENT_SECRET']
    end

    disqus_id 'hexlet-basics'

    gtm_id ENV['GOOGLE_TAG_MANAGER_KEY']

    mailer do
      from 'code-basics@hexlet.io'
    end
  end

  env :development, parent: :production do
    protocol :https
    host 'code-basics.test'
    disqus_id 'code-basics-test'
    gtm_id 'test-id'
  end

  env :test, parent: :development do
    github do
      app_id
      app_secret
    end
  end
end
