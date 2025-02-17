# frozen_string_literal: true

require "test_helper"

class Api::LessonsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @language = languages(:ruby)
    @lesson = @language.lessons.first
  end

  test "check lesson finished if exercise is correct" do
    user = sign_in_as(:one)

    language_member = @language.members.find_or_create_by!(user: user)
    lesson_member = @lesson.members.find_or_create_by!(language: @lesson.language, user: user, language_member: language_member)
    code = file_fixture("exercise/correct.rb").read
    expected = {
      passed: true,
      result: "passed"
    }

    post check_api_lesson_url(@lesson), params: { version_id: @lesson.versions.first.id, data: { attributes: { code: code } } }
    assert_response :success

    body = response.parsed_body
    language_member.reload
    lesson_member.reload

    assert { body[:passed] == expected[:passed] }
    assert { body[:result] == expected[:result] }
    assert { lesson_member.finished? }
    assert { language_member.finished? }
  end

  test "lesson result failed with incorrect solution" do
    expected = {
      passed: false,
      result: "failed"
    }
    code = file_fixture("exercise/incorrect.rb").read

    post check_api_lesson_url(@lesson), params: { version_id: @lesson.versions.first.id, data: { attributes: { code: code } } }
    assert_response :success

    body = response.parsed_body
    assert { body[:passed] == expected[:passed] }
    assert { body[:result] == expected[:result] }
  end
end
