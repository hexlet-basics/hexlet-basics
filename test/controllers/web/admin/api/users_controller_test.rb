# frozen_string_literal: true

require 'test_helper'

class Web::Admin::Api::UsersControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = sign_in_as(:admin)
  end

  test '#search' do
    get search_admin_api_users_url(format: :json)
    assert_response :success
  end
end
