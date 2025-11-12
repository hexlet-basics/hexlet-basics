# frozen_string_literal: true

require "test_helper"

class Web::RemindPasswordsControllerTest < ActionDispatch::IntegrationTest
  def test_new
    get new_remind_password_url
    assert_response :success
  end

  def test_create
    user = users(:full)
    before_token = user.reset_password_token

    post remind_password_url, params: { remind_password_form: { email: user.email } }
    assert_response :redirect

    user.reload
    assert { before_token != user.reset_password_token }
  end
end
