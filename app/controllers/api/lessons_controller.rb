# frozen_string_literal: true

class Api::LessonsController < Api::ApplicationController
  def check
    lesson = Language::Lesson.find(params[:id])
    lesson_version = lesson.versions.find(params[:version_id])
    code = params[:data][:attributes][:code]

    # FIXME: add jsonapi
    # if lesson.outdated?(lesson_version)
    #   return render json: {
    #     errors: [t('.outdated_or_deleted_lesson')]
    #   }, status: :forbidden
    # end

    language_version = lesson_version.language_version
    lesson_exercise_data = LessonTester.run(lesson_version, language_version, code, current_user)

    if lesson_exercise_data[:passed]
      lesson_member = lesson.members.find_by!(user: current_user)
      lesson_member.finish!
    end

    render json: {
      type: 'check',
      attributes: lesson_exercise_data
    }, status: :ok
  end
end
