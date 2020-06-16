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

  test 'new' do
    get new_admin_upload_path
    assert_response :success
  end

  test 'create' do
    upload_params = FactoryBot.attributes_for(:upload)

    post admin_uploads_path, params: { upload: upload_params }
    assert_response :redirect

    assert { Upload.find_by(language_name: upload_params[:language_name]) }
  end
end
