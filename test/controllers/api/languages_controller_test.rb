# frozen_string_literal: true

require 'test_helper'

class Api::LanguagesControllerTest < ActionDispatch::IntegrationTest
  test '#index' do
    get api_languages_url(format: :json)
    assert_response :success
  end

  test '#show' do
    language = languages(:php)

    get api_language_url(language, format: :json)
    assert_response :success
  end
end
