# frozen_string_literal: true

require 'test_helper'

class Web::SessionsControllerTest < ActionDispatch::IntegrationTest
  test 'new' do
    get new_session_url
    assert_response :success
  end

  test 'create' do
    user = users(:one)

    post session_url, params: { sign_in: { email: user.email, password: 'password' } }
    assert_response :redirect

    assert { signed_in? }
  end

  test 'destroy' do
    sign_in_as(:one)
    assert { signed_in? }

    delete session_url
    assert_response :redirect

    # assert { !signed_in? }
  end
end
