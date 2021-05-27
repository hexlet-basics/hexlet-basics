# frozen_string_literal: true

require 'test_helper'

class Web::PasswordsControllerTest < ActionDispatch::IntegrationTest
  test 'edit' do
    user = users(:full)

    get edit_password_url(reset_password_token: user.reset_password_token, subdomain: I18n.locale)
    assert_response :success
  end

  test 'update' do
    user = users(:full)
    before_password_digest = user.password_digest

    patch password_url(reset_password_token: user.reset_password_token, subdomain: I18n.locale), params: { user_password_form: { password: 'new_password' } }
    assert_response :redirect

    user.reload
    assert { before_password_digest != user.password_digest }
  end
end
