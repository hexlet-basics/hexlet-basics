# frozen_string_literal: true

require 'test_helper'

class Web::GoogleAuthControllerTest < ActionDispatch::IntegrationTest
  test 'check google auth' do
    post google_onetap_callback_path
    assert_redirected_to root_path
  end
end
