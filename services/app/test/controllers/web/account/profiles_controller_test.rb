# frozen_string_literal: true

require 'test_helper'

class Web::Account::ProfilesControllerTest < ActionDispatch::IntegrationTest
  test 'edit' do
    get edit_account_profile_url

    assert_response :success
  end
end
