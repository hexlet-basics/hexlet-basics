# frozen_string_literal: true

require "test_helper"

class Web::AuthControllerTest < ActionDispatch::IntegrationTest
  test "check github auth" do
    skip
    post auth_request_url(:github)
    assert_response :redirect
  end

  test "create" do
    skip
    auth_hash = generate(:github_auth_hash)
    OmniAuth.config.mock_auth[:github] = OmniAuth::AuthHash::InfoHash.new(auth_hash)

    get callback_auth_url(:github)
    assert_response :redirect

    user = User.find_by!(email: auth_hash[:info][:email].downcase)

    assert { user }
    assert { signed_in? }
  end
end
