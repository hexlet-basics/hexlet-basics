# frozen_string_literal: true

require 'test_helper'

class Web::UsersControllerTest < ActionDispatch::IntegrationTest
  test 'create' do
    user_params = { email: 'test@test.com', password: 'password' }
    post registrations_path, params: { user: user_params }
    assert_response :redirect

    assert { signed_in? }
  end
end
