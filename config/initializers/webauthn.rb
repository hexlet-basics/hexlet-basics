# typed: false
# frozen_string_literal: true

WebAuthn.configure do |config|
  config.allowed_origins = [ "#{configus.protocol}://#{configus.host}" ]
  config.rp_name = "CodeBasics"
  config.rp_id = configus.host
end
