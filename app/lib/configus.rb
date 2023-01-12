# frozen_string_literal: true

Configus.build Rails.env do
  env :production do
    hexlet_basics_release_version ENV.fetch('HEXLET_BASICS_RELEASE_VERSION', nil)

    protocol :https
    host 'code-basics.com'
    https_host 'https://code-basics.com'

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

    google do
      client_id ENV.fetch('GOOGLE_CLIENT_ID', nil)
    end

    mailer do
      from 'code-basics@hexlet.io'

      smtp do
        username ENV.fetch('SPARKPOST_SMTP_USERNAME', nil)
        password ENV.fetch('SPARKPOST_SMTP_PASSWORD', nil)
      end
    end

    sitemap do
      bucket do
        name ENV.fetch('DO_SPACES_SITEMAP_BUCKET', nil)
        credentials do
          access_key_id ENV.fetch('DO_SPACES_ACCESS_ID', nil)
          secret_access_key ENV.fetch('DO_SPACES_SECRET_KEY', nil)
          region 'us-east-1'
          endpoint "https://#{ENV.fetch('DO_SPACES_REGION', 'nyc3')}.digitaloceanspaces.com"
        end
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
