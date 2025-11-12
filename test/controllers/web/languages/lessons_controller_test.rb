# frozen_string_literal: true

require "test_helper"

class Web::Languages::LessonsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @lesson = language_lessons("javascript-module1-lesson2")
    @language = @lesson.language
    @info = @lesson.infos.last
    @user = users(:full)
  end

  def test_show
    get language_lesson_url(@language.slug, @lesson.slug)
    assert_response :success
  end

  def test_show_signed_in
    sign_in_as(:full)

    get language_lesson_url(@language.slug, @lesson.slug)
    assert_response :success
  end

  def test_show_first_lesson_signed_in
    # TODO add fixtures
    sign_in_as(:full)
    get language_lesson_url(@language.slug, @lesson.slug)
    assert_response :success
  end

  def test_show_last_lesson_signed_in
    # TODO add fixtures
    sign_in_as(:full)
    get language_lesson_url(@language.slug, @lesson.slug)
    assert_response :success
  end
end
