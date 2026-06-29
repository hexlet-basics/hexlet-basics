# typed: strict

class Api::LessonsController < Api::ApplicationController
  allow_unauthenticated_access

  sig { returns(T.untyped) }
  def check
    lesson = Language::Lesson.find(params[:id])
    language = lesson.language
    lesson_version = language.current_lesson_versions.find(params[:version_id])
    code = params[:data][:attributes][:code]

    language_version = lesson_version.language_version

    lesson_exercise_data = LessonTester.new.run(lesson_version, language_version, code, current_user)
    passed = lesson_exercise_data[:passed]

    check = CourseProgressService.record_check(
      user: current_user,
      lesson:,
      language:,
      passed:,
      locale: I18n.locale
    )

    lesson_has_been_finished = check.events.any? { it.is_a?(LessonFinishedEvent) }
    language_has_been_finished = check.events.any? { it.is_a?(CourseFinishedEvent) }

    if passed && current_user.nil?
      session[:finished_as_guest] ||= {}
      session[:finished_as_guest][lesson.id] = true
    end

    response_data = OpenStruct.new({
      **lesson_exercise_data,
      lesson_has_been_finished:,
      language_has_been_finished:
    })

    render json: LessonCheckingResponseResource.new(response_data)
  end
end
