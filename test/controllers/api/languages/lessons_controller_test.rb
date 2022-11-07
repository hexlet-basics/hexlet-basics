# frozen_string_literal: true

require 'test_helper'

class Api::Languages::LessonsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @language = languages :php
  end

  test '#index' do
    get api_language_lessons_url(@language, format: :json)
    assert_response :success
  end

  test '#show' do
    first_lesson = @language.lessons.first

    get api_language_lesson_url(@language, first_lesson, format: :json)
    assert_response :success
  end
end
