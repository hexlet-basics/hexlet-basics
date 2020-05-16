# frozen_string_literal: true

require 'test_helper'

class Web::SessionsControllerTest < ActionDispatch::IntegrationTest
  test 'new' do
    get new_sessions_path
    assert_response :success
  end

  test 'create' do
    user = users(:one)
    password = '123123'
    user.update(password: password)

    post sessions_path, params: { user: { email: user.email, password: user.password } }
    assert_response :redirect

    assert { signed_in? }
  end

  test 'destroy' do
    user = sign_in_as(:one)

    delete sessions_path
    assert_response :redirect

    assert { !signed_in? }
  end
end
