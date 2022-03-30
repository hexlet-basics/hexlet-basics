# frozen_string_literal: true

Sentry.init do |config|
  config.dsn = 'https://9aa813134b4a4cd2a5700ace36241a58@o1090356.ingest.sentry.io/6298540'
  config.breadcrumbs_logger = %i[active_support_logger http_logger]
  config.traces_sample_rate = 0.5
end
