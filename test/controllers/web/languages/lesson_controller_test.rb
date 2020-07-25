# frozen_string_literal: true

require 'test_helper'

class Web::Languages::LessonsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @language = Language.first
    @lesson = @language.lessons.first
  end

  test 'show' do
    get language_lesson_path(@language.slug, @lesson.slug)
    assert_response :success
  end

  test 'next_lesson' do
    get next_lesson_language_lesson_path(@language.slug, @lesson.slug)
    assert_response :redirect
  end

  test 'prev_lesson' do
    get prev_lesson_language_lesson_path(@language.slug, @lesson.slug)
    assert_response :redirect
  end
end
