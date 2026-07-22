# typed: strict

class Api::LessonsController < Api::ApplicationController
  allow_unauthenticated_access

  sig { void }
  def check
    lesson = Language::Lesson.find(params.expect(:id))
    struct = ApplicationParamsStruct.from_params!(LessonCheckStruct, params.require(:data))

    version_id = T.must(struct.version_id)
    lesson_version = lesson.language.current_lesson_versions.find(version_id)
    language_version = lesson_version.language_version

    result = CourseProgressService.record_check(
      user: current_user,
      lesson:,
      lesson_version:,
      language_version:,
      code: struct.code.to_s,
      locale: I18n.locale.to_s
    )

    case result
    when Typed::Success
      check = result.payload
      if check.passed && current_user.nil?
        session[:finished_as_guest] ||= {}
        session[:finished_as_guest][lesson.id] = true
      end
      render json: LessonCheckingResponseResource.new(check)
    when Typed::Failure
      head :unprocessable_content
    end
  end
end
