# frozen_string_literal: true

require 'test_helper'

class Web::Account::ProfilesControllerTest < ActionDispatch::IntegrationTest
  def setup
    sign_in_as(:full)
  end

  test 'edit' do
    get edit_account_profile_url

    assert_response :success
  end

  test 'destroy' do
    user = users :full
    delete account_profile_url

    assert_response :redirect

    user.reload

    assert { user.removed? }
    assert { !signed_in? }
  end
end
