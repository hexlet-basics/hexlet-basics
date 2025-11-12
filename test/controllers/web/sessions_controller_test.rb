# frozen_string_literal: true

require "test_helper"

class Web::SessionsControllerTest < ActionDispatch::IntegrationTest
  def test_new
    get new_session_url
    assert_response :success
  end

  def test_create
    user = users(:one)

    post session_url, params: { user: { email: user.email, password: "password" } }
    assert_response :redirect

    assert { signed_in? }
  end

  def test_destroy
    sign_in_as(:one)
    assert { signed_in? }

    delete session_url
    assert_response :redirect

    # assert { !signed_in? }
  end
end
