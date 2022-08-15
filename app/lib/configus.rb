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

    disqus do
      ru 'hexlet-basics'
      en 'hexlet-basics-en'
    end

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
    gtm_id 'test-id'
    disqus do
      ru 'code-basics-test'
      en 'code-basics-test'
    end
  end

  env :test, parent: :development do
    github do
      app_id
      app_secret
    end
  end
end
