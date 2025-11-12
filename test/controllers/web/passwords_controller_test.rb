# frozen_string_literal: true

require "test_helper"

class Web::PasswordsControllerTest < ActionDispatch::IntegrationTest
  def test_edit
    user = users(:full)

    get edit_password_url(reset_password_token: user.reset_password_token)
    assert_response :success
  end

  def test_update
    user = users(:full)
    before_password_digest = user.password_digest

    patch password_url(reset_password_token: user.reset_password_token), params: { user: { password: "new_password" } }
    assert_response :redirect

    user.reload

    assert { before_password_digest != user.password_digest }
  end
end
