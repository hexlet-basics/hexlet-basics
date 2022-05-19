# frozen_string_literal: true

Configus.build Rails.env do
  env :production do
    hexlet_basics_release_version ENV.fetch('HEXLET_BASICS_RELEASE_VERSION', nil)

    protocol :https
    host 'code-basics.com'
    https_host 'https://ru.code-basics.com'

    github do
      app_id ENV.fetch('GITHUB_CLIENT_ID', nil)
      app_secret ENV.fetch('GITHUB_CLIENT_SECRET', nil)
    end

    facebook do
      app_id ENV.fetch('FACEBOOK_CLIENT_ID', nil)
      app_secret ENV.fetch('FACEBOOK_CLIENT_SECRET', nil)
    end

    disqus_id 'hexlet-basics'

    gtm_id ENV.fetch('GOOGLE_TAG_MANAGER_KEY', nil)

    mailer do
      from 'code-basics@hexlet.io'

      smtp do
        username ENV.fetch('SPARKPOST_SMTP_USERNAME', nil)
        password ENV.fetch('SPARKPOST_SMTP_PASSWORD', nil)
      end
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
