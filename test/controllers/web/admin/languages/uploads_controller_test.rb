# frozen_string_literal: true

require 'test_helper'

class Web::Admin::Languages::UploadsControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  test 'create' do
    language = languages(:one)
    AppContainer.stub(:docker_image_exercise_loader, :success)

    post admin_language_uploads_path(language.id)
    assert_response :redirect

    assert { Language::Upload.find_by(language: language.id) }
  end
end
