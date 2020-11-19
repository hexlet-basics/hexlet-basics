# frozen_string_literal: true

require 'test_helper'

class Web::SessionsControllerTest < ActionDispatch::IntegrationTest
  test 'new' do
    get new_session_path
    assert_response :success
  end

  test 'create' do
    user = users(:one)

    post session_path, params: { sign_in: { email: user.email, password: 'password' } }
    assert_response :redirect

    assert { signed_in? }
  end

  test 'destroy' do
    sign_in_as(:one)

    delete session_path
    assert_response :redirect

    assert { !signed_in? }
  end
end
