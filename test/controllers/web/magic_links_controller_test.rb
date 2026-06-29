# typed: true
# frozen_string_literal: true

require "test_helper"

class Web::MagicLinksControllerTest < ActionDispatch::IntegrationTest
  setup do
    # Test env uses :null_store; the per-email cooldown needs a real cache.
    @original_cache = Rails.cache
    Rails.cache = ActiveSupport::Cache::MemoryStore.new
    @user = users(:one)
  end

  teardown do
    Rails.cache = @original_cache
  end

  def test_new
    get new_magic_link_url
    assert_response :success
  end

  def test_create_sends_link_for_existing_user
    assert_emails 1 do
      post magic_links_url, params: { data: { email: @user.email } }
    end
    assert_redirected_to root_url
  end

  def test_create_does_not_send_for_unknown_email_but_succeeds
    assert_no_emails do
      post magic_links_url, params: { data: { email: "nobody@example.com" } }
    end
    assert_redirected_to root_url
  end

  def test_show_with_valid_token_signs_in
    token = @user.generate_token_for(:magic_link)

    get magic_link_url(token)
    assert_redirected_to root_url
    assert { authenticated? }
  end

  def test_show_with_invalid_token_does_not_sign_in
    get magic_link_url("garbage-token")
    assert_redirected_to new_magic_link_url
    assert { !authenticated? }
  end
end
