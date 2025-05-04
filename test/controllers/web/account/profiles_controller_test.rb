# frozen_string_literal: true

require "test_helper"

class Web::Account::ProfilesControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = users :full
    sign_in_as(:full)
  end

  test "edit" do
    get edit_account_profile_url
    assert_response :success
  end

  test "update" do
    new_name = "new first name"

    patch account_profile_url(id: @user.id), params: {
      user: {
        first_name: new_name,
        contact_method: :telegram,
        contact_value: "@orgprog"
      }
    }
    assert_response :redirect

    @user.reload

    assert { @user.first_name == new_name }
    assert { @user.lead }
  end

  test "destroy" do
    delete account_profile_url
    assert_response :redirect

    @user.reload

    assert { @user.removed? }
    assert { !signed_in? }
  end
end
