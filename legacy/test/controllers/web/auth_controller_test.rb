# typed: true
# frozen_string_literal: true

require "test_helper"

class Web::AuthControllerTest < ActionDispatch::IntegrationTest
  def test_check_github_auth
    skip("не адаптирован после переезда на Inertia, см. 6fbf064d")
    post T.unsafe(self).auth_request_url(:github)

    assert_response :redirect
  end

  def test_create
    skip("не адаптирован после переезда на Inertia, см. 6fbf064d")
    auth_hash = generate(:github_auth_hash)
    T.unsafe(OmniAuth).config.mock_auth[:github] = T.unsafe(OmniAuth::AuthHash::InfoHash).new(auth_hash)

    get T.unsafe(self).callback_auth_url(:github)

    assert_response :redirect

    user = User.find_by!(email: auth_hash[:info][:email].downcase)

    assert { user }
    assert { authenticated? }
  end
end
