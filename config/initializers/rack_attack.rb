# typed: false

# Throttle abusive traffic. Counters live in Rails.cache (solid_cache), so limits
# are shared across processes without a separate Redis. Routes carry an optional
# `/:suffix` locale prefix, so auth paths are matched with anchored regexes.
class Rack::Attack
  ### Allow-list ###

  # Never throttle the health check.
  safelist("allow/health-check") do |req|
    req.path == "/up"
  end

  ### Throttles ###

  # Broad backstop against scraping/floods: 300 requests per IP per 5 minutes.
  throttle("req/ip", limit: 300, period: 5.minutes) do |req|
    req.ip unless req.path.start_with?("/assets", "/vite", "/packs")
  end

  # Login attempts by IP: 10 per 20 seconds.
  throttle("logins/ip", limit: 10, period: 20.seconds) do |req|
    req.ip if req.post? && req.path.match?(%r{/(session|passkey_session)\z})
  end

  # Login attempts by submitted email, independent of IP (credential stuffing).
  throttle("logins/email", limit: 10, period: 1.minute) do |req|
    if req.post? && req.path.match?(%r{/session\z})
      email = req.params.dig("data", "email") || req.params["email"]
      email.to_s.downcase.presence
    end
  end

  # Sign-ups by IP: 5 per minute.
  throttle("signups/ip", limit: 5, period: 1.minute) do |req|
    req.ip if req.post? && req.path.match?(%r{/users\z}) && !req.path.include?("/admin/")
  end

  # Phone-code requests by IP: 5 per minute (complements PhoneAuthService limits).
  throttle("phone_auth/ip", limit: 5, period: 1.minute) do |req|
    req.ip if req.post? && req.path.match?(%r{/phone_auth\z})
  end

  # Password-reset requests by IP: 5 per minute.
  throttle("password_reset/ip", limit: 5, period: 1.minute) do |req|
    req.ip if req.post? && req.path.match?(%r{/remind_password\z})
  end

  ### Response ###

  self.throttled_responder = lambda do |request|
    retry_after = (request.env["rack.attack.match_data"] || {})[:period]
    [
      429,
      { "Content-Type" => "text/plain", "Retry-After" => retry_after.to_s },
      [ "Too many requests. Please retry later.\n" ]
    ]
  end
end

# Surface throttling in the logs so blocks are visible in production.
ActiveSupport::Notifications.subscribe("throttle.rack_attack") do |_name, _start, _finish, _id, payload|
  req = payload[:request]
  Rails.logger.warn("[rack-attack] throttled #{req.env['rack.attack.matched']} ip=#{req.ip} path=#{req.path}")
end
