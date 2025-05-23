# frozen_string_literal: true

require "test_helper"

class Web::Account::ProfilesControllerTest < ActionDispatch::IntegrationTest
  test "edit" do
    sign_in_as(:full)

    get edit_account_profile_url
    assert_response :success
  end

  test "update" do
    user = sign_in_as(:full)
    new_name = "new first name"

    patch account_profile_url(id: user.id), params: {
      user: {
        first_name: new_name
      }
    }
    assert_response :redirect

    user.reload

    assert { user.first_name == new_name }
  end

  test "destroy" do
    user = sign_in_as(:full)
    delete account_profile_url
    assert_response :redirect

    user.reload

    assert { user.removed? }
    assert { !signed_in? }
  end
end
