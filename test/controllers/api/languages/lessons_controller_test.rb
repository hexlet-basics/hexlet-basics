# frozen_string_literal: true

require 'test_helper'

class Api::Languages::LessonsControllerTest < ActionDispatch::IntegrationTest
  test '#show' do
    language = languages(:php)
    first_lesson = language.lessons.first

    get api_language_lesson_url language, first_lesson, format: :json
    assert_response :success
  end
end
