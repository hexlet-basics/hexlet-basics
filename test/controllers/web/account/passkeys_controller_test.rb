# typed: true
# frozen_string_literal: true

require "test_helper"
require "webauthn/fake_client"

class Web::Account::PasskeysControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:one)
    @client = WebAuthn::FakeClient.new("#{configus.protocol}://#{configus.host}")
  end

  def test_register_passkey
    get new_account_passkey_url

    assert_response :success
    options = JSON.parse(@response.body)
    credential = @client.create(challenge: options["challenge"])

    assert_difference -> { @user.credentials.count }, 1 do
      post account_passkeys_url, params: { credential: credential.to_json }
    end
    assert_redirected_to edit_account_profile_url
    assert { @user.reload.webauthn_id? }
  end

  def test_destroy_passkey
    get new_account_passkey_url
    credential = @client.create(challenge: JSON.parse(@response.body)["challenge"])
    post account_passkeys_url, params: { credential: credential.to_json }
    passkey = @user.credentials.sole

    assert_difference -> { @user.credentials.count }, -1 do
      delete account_passkey_url(passkey)
    end
  end
end
