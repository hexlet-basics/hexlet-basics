# typed: true
# frozen_string_literal: true

require "test_helper"
require "webauthn/fake_client"

class Web::PasskeySessionsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @client = WebAuthn::FakeClient.new("#{configus.protocol}://#{configus.host}")
    register_passkey_for(:one)
  end

  def test_new_returns_options
    get new_passkey_session_url

    assert_response :success
    assert { JSON.parse(@response.body)["challenge"].present? }
  end

  def test_login_with_passkey
    get new_passkey_session_url
    assertion = @client.get(challenge: JSON.parse(@response.body)["challenge"])

    post passkey_session_url, params: { credential: assertion.to_json }

    assert_redirected_to root_url
    assert { authenticated? }
  end

  def test_unknown_credential_does_not_sign_in
    other = WebAuthn::FakeClient.new("#{configus.protocol}://#{configus.host}")
    get new_passkey_session_url
    # `other` never registered, so its credential id is unknown to us.
    assertion = begin
      other.get(challenge: JSON.parse(@response.body)["challenge"])
    rescue RuntimeError
      nil
    end

    if assertion
      post passkey_session_url, params: { credential: assertion.to_json }
      assert { !authenticated? }
    end
  end

  private

  def register_passkey_for(fixture)
    sign_in_as(fixture)
    get new_account_passkey_url
    credential = @client.create(challenge: JSON.parse(@response.body)["challenge"])
    post account_passkeys_url, params: { credential: credential.to_json }
    sign_out
  end
end
