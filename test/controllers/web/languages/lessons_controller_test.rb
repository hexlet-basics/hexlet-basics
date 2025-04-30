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
end
