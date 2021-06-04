# frozen_string_literal: true

require 'test_helper'

class Web::RemindPasswordsControllerTest < ActionDispatch::IntegrationTest
  test 'new' do
    get new_remind_password_url(subdomain: I18n.locale)
    assert_response :success
  end

  test 'create' do
    user = users(:full)
    before_token = user.reset_password_token

    post remind_password_url(subdomain: I18n.locale), params: { remind_password: { email: user.email } }
    assert_response :redirect

    user.reload
    assert { before_token != user.reset_password_token }
  end
end
