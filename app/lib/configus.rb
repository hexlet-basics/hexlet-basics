# frozen_string_literal: true

Configus.build Rails.env do
  env :production do
    hexlet_basics_release_version ENV.fetch("HEXLET_BASICS_RELEASE_VERSION", nil)

    protocol :https
    host "code-basics.com"
    https_host "https://code-basics.com"

    github do
      app_id ENV.fetch("GITHUB_CLIENT_ID", nil)
      app_secret ENV.fetch("GITHUB_CLIENT_SECRET", nil)
    end

    facebook do
      app_id ENV.fetch("FACEBOOK_CLIENT_ID", nil)
      app_secret ENV.fetch("FACEBOOK_CLIENT_SECRET", nil)
    end

    disqus do
      ru "hexlet-basics"
      en "hexlet-basics-en"
    end

    gtm_id ENV.fetch("GOOGLE_TAG_MANAGER_KEY", nil)

    google do
      client do
        id ENV.fetch("GOOGLE_CLIENT_ID", nil)
        secret ENV.fetch("GOOGLE_CLIENT_SECRET", nil)
      end
    end

    mailer do
      from "support@hexlet.io"

      smtp do
        username ENV.fetch("SPARKPOST_SMTP_USERNAME", nil)
        password ENV.fetch("SPARKPOST_SMTP_PASSWORD", nil)
      end
    end

    sitemap do
      bucket do
        name ENV.fetch("SITEMAPS_S3_BUCKET", nil)
        credentials do
          access_key_id ENV.fetch("SITEMAPS_S3_ACCESS_KEY_ID", nil)
          secret_access_key ENV.fetch("SITEMAPS_S3_ACCESS_SECRET_KEY", nil)
          region ENV.fetch("SITEMAPS_S3_REGION", nil)
          endpoint ENV.fetch("SITEMAPS_S3_ENDPOINT", nil)
        end
      end
    end

    storage do
      region ENV.fetch("STORAGE_S3_REGION", nil)
      endpoint ENV.fetch("STORAGE_S3_ENDPOINT", nil)
      access_key_id ENV.fetch("STORAGE_S3_ACCESS_KEY_ID", nil)
      secret_access_key ENV.fetch("STORAGE_S3_ACCESS_SECRET_KEY", nil)
      bucket ENV.fetch("STORAGE_S3_BUCKET", nil)
    end

    csp_report_uri ENV.fetch("CSP_REPORT_URI", "")
  end

  env :staging, parent: :production do
  end

  env :development, parent: :production do
    protocol :https
    host "code-basics.test"
    gtm_id "test-id"
    disqus do
      ru "code-basics-test"
      en "code-basics-test"
    end
  end

  env :test, parent: :development do
    github do
      app_id
      app_secret
    end
  end
end
