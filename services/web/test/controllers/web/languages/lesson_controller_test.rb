# frozen_string_literal: true

require 'test_helper'

class Web::Languages::LessonsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @lesson = language_lessons(:two)
    @language = @lesson.language
  end

  test 'show' do
    get language_lesson_path(@language.slug, @lesson.slug)
    assert_response :success
  end

  test 'next_lesson' do
    third_language_lesson = language_lessons(:three)

    get next_lesson_language_lesson_path(@language.slug, @lesson.slug)
    assert_response :redirect

    assert_redirected_to language_lesson_path(@language.slug, third_language_lesson.slug)
  end

  test 'prev_lesson' do
    first_language_lesson = language_lessons(:one)

    get prev_lesson_language_lesson_path(@language.slug, @lesson.slug)
    assert_response :redirect

    assert_redirected_to language_lesson_path(@language.slug, first_language_lesson.slug)
  end
end
