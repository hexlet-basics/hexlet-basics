# frozen_string_literal: true

class Api::LessonsController < Api::ApplicationController
  def check
    lesson = Language::Lesson.find(params[:id])
    lesson_version = lesson.versions.find(params[:data][:attributes][:version_id])
    code = params[:data][:attributes][:code]

    if lesson.outdated?(lesson_version)
      return render json: {
        message: t('.outdated_or_deleted_lesson')
      }, status: :gone
    end

    language_version = lesson_version.language_version
    lesson_exercise_data = LessonTester.run(lesson_version, language_version, code, current_user)

    if lesson_exercise_data[:passed]
      lesson_member = lesson.members.find_by!(user: current_user)
      lesson_member.finish!
    end

    render json: {
      attributes: lesson_exercise_data
    }, status: :ok
  end
end
