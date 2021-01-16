# frozen_string_literal: true

require 'test_helper'

class Web::UsersControllerTest < ActionDispatch::IntegrationTest
  test 'create' do
    user_params = FactoryBot.attributes_for(:user)
    post users_url, params: { user_sign_up_form: user_params }
    assert_response :redirect

    assert { signed_in? }
  end
end
