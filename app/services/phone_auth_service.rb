# typed: strict
# frozen_string_literal: true

# Passwordless phone login: send a one-time code by SMS, then verify it.
# Codes live in Rails.cache (SolidCache), never in the DB.
class PhoneAuthService < ApplicationService
  CODE_TTL = T.let(5.minutes, ActiveSupport::Duration)
  RESEND_INTERVAL = T.let(1.minute, ActiveSupport::Duration)
  MAX_ATTEMPTS = 5
  MAX_PER_IP_PER_HOUR = 10

  class << self
    extend T::Sig

    # Returns the normalized E.164 phone on success, or a reason symbol on failure
    # (:invalid_phone, :too_frequent, :ip_rate_limited).
    sig { params(raw_phone: String, ip: T.nilable(String)).returns(Typed::Result[String, Symbol]) }
    def request_code(raw_phone, ip: nil)
      phone = Phonelib.parse(raw_phone).e164.presence
      return fail_with(:invalid_phone) unless phone && Phonelib.valid?(phone)

      return fail_with(:too_frequent) if Rails.cache.exist?(resend_key(phone))
      return fail_with(:ip_rate_limited) if ip && ip_count(ip) >= MAX_PER_IP_PER_HOUR

      code = generate_code
      Rails.cache.write(code_key(phone), { code:, attempts: 0 }, expires_in: CODE_TTL)
      Rails.cache.write(resend_key(phone), true, expires_in: RESEND_INTERVAL)
      bump_ip(ip) if ip

      SendSmsJob.perform_later(phone:, text: code_text(code))

      success_with(phone)
    end

    # Returns the signed-in/created User on success, or a reason symbol on failure
    # (:invalid_phone, :expired, :too_many_attempts, :invalid_code).
    sig { params(raw_phone: String, code: String).returns(Typed::Result[User, Symbol]) }
    def verify_code(raw_phone, code)
      phone = Phonelib.parse(raw_phone).e164.presence
      return fail_with(:invalid_phone) unless phone

      entry = Rails.cache.read(code_key(phone))
      return fail_with(:expired) unless entry

      if entry[:attempts] >= MAX_ATTEMPTS
        Rails.cache.delete(code_key(phone))
        return fail_with(:too_many_attempts)
      end

      unless ActiveSupport::SecurityUtils.secure_compare(entry[:code], code.to_s.strip)
        Rails.cache.write(code_key(phone), entry.merge(attempts: entry[:attempts] + 1), expires_in: CODE_TTL)
        return fail_with(:invalid_code)
      end

      Rails.cache.delete(code_key(phone))
      Rails.cache.delete(resend_key(phone))

      user = User.find_or_initialize_by(phone:)
      user.phone_verified_at = Time.current
      user.save!

      success_with(user)
    end

    private

    sig { returns(String) }
    def generate_code
      format("%04d", SecureRandom.random_number(10_000))
    end

    sig { params(code: String).returns(String) }
    def code_text(code)
      I18n.t("sms.phone_auth.code", code:)
    end

    sig { params(phone: String).returns(String) }
    def code_key(phone) = "phone_auth:code:#{phone}"

    sig { params(phone: String).returns(String) }
    def resend_key(phone) = "phone_auth:resend:#{phone}"

    sig { params(ip: String).returns(String) }
    def ip_key(ip) = "phone_auth:ip:#{ip}"

    sig { params(ip: String).returns(Integer) }
    def ip_count(ip) = Rails.cache.read(ip_key(ip)).to_i

    sig { params(ip: String).void }
    def bump_ip(ip)
      Rails.cache.write(ip_key(ip), ip_count(ip) + 1, expires_in: 1.hour)
    end
  end
end
