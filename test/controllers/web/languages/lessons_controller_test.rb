# frozen_string_literal: true

require "test_helper"

class Web::Languages::LessonsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @lesson = language_lessons("javascript-module1-lesson2")
    @language = @lesson.language
    @info = @lesson.infos.last
    @user = users(:full)
  end

  test "show" do
    get language_lesson_url(@language.slug, @lesson.slug)
    assert_response :success
  end

  # test "show amp" do
  #   get language_lesson_url(@language.slug, @lesson.slug, format: :amp)
  #   assert_response :success
  # end

  test "show (signed in)" do
    sign_in_as(:full)

    get language_lesson_url(@language.slug, @lesson.slug)
    assert_response :success
  end

  test "show first lesson (signed in)" do
    # TODO add fixtures
    sign_in_as(:full)
    get language_lesson_url(@language.slug, @lesson.slug)
    assert_response :success
  end

  test "show last lesson (signed in)" do
    # TODO add fixtures
    sign_in_as(:full)
    get language_lesson_url(@language.slug, @lesson.slug)
    assert_response :success
  end

  # test "next_lesson" do
  #   sign_in_as(:full)
  #   third_language_lesson = language_lessons("javascript-module2-lesson1")
  #
  #   get next_lesson_language_lesson_url(@language.slug, @lesson.slug)
  #   assert_response :redirect
  #
  #   assert_redirected_to language_lesson_url(@language.slug, third_language_lesson.slug)
  # end

  # test "prev_lesson" do
  #   first_language_lesson = language_lessons("javascript-module1-lesson1")
  #
  #   get prev_lesson_language_lesson_url(@language.slug, @lesson.slug)
  #   assert_response :redirect
  #
  #   assert_redirected_to language_lesson_url(@language.slug, first_language_lesson.slug)
  # end
end
