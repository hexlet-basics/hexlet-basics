# typed: false

# Collapse Rails' multi-line request logs into a single structured line per
# request. Only enabled where logs are shipped to STDOUT (production/staging).
Rails.application.configure do
  config.lograge.enabled = Rails.env.production? || Rails.env.staging?

  # Keep request correlation with Sentry: `request_id` is otherwise carried by
  # `config.log_tags`, which lograge bypasses.
  config.lograge.custom_options = lambda do |event|
    {
      request_id: event.payload[:request_id],
      host: event.payload[:host],
      params: event.payload[:params]&.except("controller", "action", "format", "utf8"),
      user_id: event.payload[:user_id]
    }.compact
  end
end
