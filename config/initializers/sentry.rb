# frozen_string_literal: true

# NOTE Config via env variable SENTRY_DSN
Sentry.init do |config|
  config.release = "hexle-basics-platform@#{configus.hexlet_basics_release_version}"

  config.breadcrumbs_logger = %i[active_support_logger http_logger]
  config.traces_sample_rate = 0.002
end
