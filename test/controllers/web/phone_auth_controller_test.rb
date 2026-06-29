# typed: true
# frozen_string_literal: true

require "test_helper"

class Web::PhoneAuthControllerTest < ActionDispatch::IntegrationTest
  setup do
    # Test env uses :null_store; OTP needs a real cache to round-trip.
    @original_cache = Rails.cache
    Rails.cache = ActiveSupport::Cache::MemoryStore.new
    @phone = "+79161234567"
  end

  teardown do
    Rails.cache = @original_cache
  end

  def cached_code
    Rails.cache.read("phone_auth:code:#{@phone}")&.fetch(:code)
  end

  def test_new
    get new_phone_auth_url
    assert_response :success
  end

  def test_request_code_stores_code_and_redirects_to_verify
    post phone_auth_url, params: { data: { phone: @phone } }

    assert_redirected_to verify_phone_auth_url(phone: @phone)
    assert { cached_code.present? }
  end

  def test_request_code_rejects_invalid_phone
    post phone_auth_url, params: { data: { phone: "not-a-phone" } }
    assert_redirected_to new_phone_auth_url
    assert { cached_code.nil? }
  end

  def test_full_login_creates_user_and_signs_in
    post phone_auth_url, params: { data: { phone: @phone } }
    code = cached_code

    assert_difference -> { User.count }, 1 do
      post confirm_phone_auth_url, params: { data: { phone: @phone, code: } }
    end

    assert_redirected_to root_url
    assert { authenticated? }

    user = User.find_by(phone: @phone)
    assert { user.present? }
    assert { T.must(user).phone_verified_at.present? }
  end

  def test_existing_user_signs_in_without_duplicate
    User.create!(phone: @phone)

    post phone_auth_url, params: { data: { phone: @phone } }
    code = cached_code

    assert_no_difference -> { User.count } do
      post confirm_phone_auth_url, params: { data: { phone: @phone, code: } }
    end
    assert { authenticated? }
  end

  def test_wrong_code_does_not_sign_in
    post phone_auth_url, params: { data: { phone: @phone } }

    post confirm_phone_auth_url, params: { data: { phone: @phone, code: "0000-wrong" } }
    assert_redirected_to verify_phone_auth_url(phone: @phone)
    assert { !authenticated? }
  end
end
