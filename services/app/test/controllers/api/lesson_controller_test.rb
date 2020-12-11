# frozen_string_literal: true

require 'test_helper'

class Api::LessonsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @language = languages(:php)
    @lesson = @language.lessons.first
    @user = users(:one)
    @lesson_member = @lesson.members.find_or_create_by!(language: @lesson.language, user: @user)
  end

  test 'check if user logged in' do
    post check_api_lesson_path(@lesson), params: { version_id: @lesson.versions.first.id, data: { attributes: { code: 'code' } } }
    assert_response :forbidden
  end

  test 'check lesson finished' do
    sign_in_as(:one)
    post check_api_lesson_path(@lesson), params: { version_id: @lesson.versions.first.id, data: { attributes: { code: 'code' } } }

    assert_response :success
    @lesson_member.reload
    assert { @lesson_member.finished? }
  end

  test 'check language finished' do
    sign_in_as(:one)
    language_member = @language.members.find_or_create_by!(user: @user)
    post check_api_lesson_path(@lesson), params: { version_id: @lesson.versions.first.id, data: { attributes: { code: 'code' } } }
    assert_response :success

    language_member.reload
    assert { language_member.finished? }
  end
end
