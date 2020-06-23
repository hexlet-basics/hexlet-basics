# frozen_string_literal: true

require 'test_helper'

class Web::Admin::UploadsControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  test 'index' do
    get admin_uploads_path
    assert_response :success
  end
end
