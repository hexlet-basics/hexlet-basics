# typed: strict

class Api::LessonsController < Api::ApplicationController
  allow_unauthenticated_access

  sig { returns(T.untyped) }
  def check
    lesson = Language::Lesson.find(params[:id])
    struct = ApplicationParamsStruct.from_params!(LessonCheckStruct, params.require(:data))

    lesson_version = lesson.language.current_lesson_versions.find(T.must(struct.version_id))
    language_version = lesson_version.language_version

    check = CourseProgressService.record_check(
      user: current_user,
      lesson:,
      lesson_version:,
      language_version:,
      code: struct.code.to_s,
      locale: I18n.locale
    )

    lesson_has_been_finished = check.events.any? { it.is_a?(LessonFinishedEvent) }
    language_has_been_finished = check.events.any? { it.is_a?(CourseFinishedEvent) }

    if check.exercise_data[:passed] && current_user.nil?
      session[:finished_as_guest] ||= {}
      session[:finished_as_guest][lesson.id] = true
    end

    response_data = OpenStruct.new({
      **check.exercise_data,
      lesson_has_been_finished:,
      language_has_been_finished:
    })

    render json: LessonCheckingResponseResource.new(response_data)
  end
end
