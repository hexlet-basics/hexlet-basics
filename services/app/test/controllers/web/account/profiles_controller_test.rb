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
end
