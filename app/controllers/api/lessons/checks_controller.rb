# frozen_string_literal: true

class Api::Lessons::ChecksController < Api::Lessons::ApplicationController
  def create
    lesson_version = resource_lesson.versions.find(params[:data][:attributes][:version_id])
    code = params[:data][:attributes][:code]

    language = resource_lesson.language

    if language.current_lessons.exclude?(lesson_version)
      return render json: {
        message: "lesson is outdated or deleted"
      }, status: 410
    end

    language_version = lesson_version.language_version
    data = LessonTester.run(lesson_version, language_version, code, current_user)

    if data[:passed]
      lesson_member = resource_lesson.members.find_by!(user: current_user)
      lesson_member.finish!
    end

    render json: {
      attributes: data
    }, status: 200
  end
end
