# frozen_string_literal: true

require 'test_helper'

class Web::Admin::Management::UsersControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  test 'index' do
    get admin_management_users_url(subdomain: subdomain)

    assert_response :success
  end
end
