# frozen_string_literal: true

require "test_helper"

class Web::GoogleAuthControllerTest < ActionDispatch::IntegrationTest
  test "check google auth" do
    skip

    post google_onetap_callback_path
    assert_redirected_to root_path
  end

  test "create user" do
    skip

    headers = { "Cookie" => "g_csrf_token=g_csrf_token_test;" }
    params = {
      g_csrf_token: "g_csrf_token_test"
    }

    email = "example@mail.com"

    open_session do |s|
      s.post s.google_onetap_callback_url, params: params, headers: headers
      s.assert_redirected_to root_path
    end

    user = User.find_by! email: email

    assert { user.active? }
  end
end
