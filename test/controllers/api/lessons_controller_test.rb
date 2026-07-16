# typed: true
# frozen_string_literal: true

require "test_helper"

class Api::LessonsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @language = languages(:ruby)
    @lesson = @language.lessons.first
  end

  def test_check_lesson_finished_if_exercise_is_correct
    user = sign_in_as(:one)

    language_member = @language.members.find_or_create_by!(user: user)
    lesson_member = @lesson.members.find_or_create_by!(language: @lesson.language, user: user, language_member: language_member)
    code = file_fixture("exercise/correct.rb").read
    expected = {
      passed: true,
      result: "passed"
    }

    post check_api_lesson_url(@lesson), params: { data: { version_id: @lesson.versions.first.id, code: code } }
    assert_response :success

    body = response.parsed_body
    language_member.reload
    lesson_member.reload

    assert { body[:passed] == expected[:passed] }
    assert { body[:result] == expected[:result] }
    assert { lesson_member.finished? }
    assert { language_member.finished? }
  end

  def test_lesson_result_failed_with_incorrect_solution
    expected = {
      passed: false,
      result: "failed"
    }
    code = file_fixture("exercise/incorrect.rb").read

    post check_api_lesson_url(@lesson), params: { data: { version_id: @lesson.versions.first.id, code: code } }
    assert_response :success

    body = response.parsed_body
    assert { body[:passed] == expected[:passed] }
    assert { body[:result] == expected[:result] }
  end

  def test_check_returns_unprocessable_entity_without_version_id
    code = file_fixture("exercise/correct.rb").read

    post check_api_lesson_url(@lesson), params: { data: { code: code } }
    assert_response :unprocessable_entity
  end

  def test_check_finishes_multiple_lessons_without_creating_survey_state
    code = file_fixture("exercise/correct.rb").read

    user = sign_in_as(:ready_to_start_learning)
    initial_survey_scenarios_count = user.survey_scenarios.count
    language = languages(:javascript)
    language_member = language.members.create!(user: user)

    lesson1 = language.current_lessons.ordered.first!
    lesson_version1 = language.current_lesson_versions.find_by!(lesson: lesson1)
    lesson1.members.create!(language:, user:, language_member:)

    post check_api_lesson_url(lesson1), params: { data: { version_id: lesson_version1.id, code: } }
    assert_response :success
    assert { response.parsed_body["passed"] }

    lesson2 = lesson_version1.next_lesson
    lesson_version2 = language.current_lesson_versions.find_by!(lesson: lesson2)
    lesson2.members.create!(language:, user:, language_member:)

    post check_api_lesson_url(lesson2), params: { data: { version_id: lesson_version2.id, code: } }
    assert_response :success
    assert { response.parsed_body["passed"] }

    lesson3 = lesson_version2.next_lesson
    lesson_version3 = language.current_lesson_versions.find_by!(lesson: lesson3)
    lesson3.members.create!(language:, user:, language_member:)

    post check_api_lesson_url(lesson3), params: { data: { version_id: lesson_version3.id, code: } }
    assert_response :success
    assert { response.parsed_body["passed"] }

    assert { user.survey_scenarios.count == initial_survey_scenarios_count }
  end
end
