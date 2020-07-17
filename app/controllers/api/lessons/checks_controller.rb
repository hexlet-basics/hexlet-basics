# frozen_string_literal: true

class Api::Lessons::ChecksController < Api::Lessons::ApplicationController
  def create
    lesson_version = resource_lesson.versions.find(params[:version_id])
    code = params[:data][:attributes][:code]

    language_version = lesson_version.language_version
    check_data = CheckLesson.run(lesson_version, language_version, code, current_user)

    if check_data[:passed]
      lesson_member = resource_lesson.members.find_by(user: current_user, lesson_version: lesson_version)
      lesson_member.finish!
    end

    render json: {
      type: 'check',
      attributes: check_data
    }
  end
end
