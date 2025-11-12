# frozen_string_literal: true

require "test_helper"

class Web::Admin::Management::UsersControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  def test_index
    get admin_management_users_url

    assert_response :success
  end

  def test_edit
    user = users(:one)

    get edit_admin_management_user_url(user)

    assert_response :success
  end

  def test_update
    user = users(:one)

    attrs = {
      admin: true
    }

    patch admin_management_user_url(user), params: { user: attrs }

    assert_response :redirect
  end
end
