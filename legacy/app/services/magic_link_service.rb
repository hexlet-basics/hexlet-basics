# typed: strict
# frozen_string_literal: true

# Passwordless email sign-in: email the user a one-time login link.
# Anti-enumeration: callers always report success regardless of whether the
# email exists — this service silently no-ops for unknown or cooling-down emails.
class MagicLinkService < ApplicationService
  RESEND_INTERVAL = T.let(1.minute, ActiveSupport::Duration)

  class << self
    extend T::Sig

    sig { params(email: String, suffix: T.nilable(String)).void }
    def request!(email, suffix:)
      normalized = email.to_s.strip.downcase.presence
      return unless normalized
      return if Rails.cache.exist?(cooldown_key(normalized))

      user = User.find_by(email: normalized)
      return unless user

      Rails.cache.write(cooldown_key(normalized), true, expires_in: RESEND_INTERVAL)
      UserMailer.with(user:, suffix:).magic_link.deliver_later
    end

    private

    sig { params(email: String).returns(String) }
    def cooldown_key(email) = "magic_link:cooldown:#{email}"
  end
end
