# frozen_string_literal: true

class Api::LessonsController < Api::ApplicationController
  # include Gon::ControllerHelpers
  # include EventConcern
  # before_action :require_api_auth!

  def check
    lesson = Language::Lesson.find(params[:id])
    language = lesson.language
    lesson_version = language.current_lesson_versions.find(params[:version_id])
    code = params[:data][:attributes][:code]

    language_version = lesson_version.language_version
    lesson_exercise_data = LessonTester.new.run(lesson_version, language_version, code, current_user)

    if lesson_exercise_data[:passed] && !current_user.guest?
      lesson_member = lesson.members.find_by!(user: current_user)
      lesson_member.finish!

      # lesson_finished_event_options = {
      #   user: current_user.serializable_data,
      #   language: language.serializable_data,
      #   lesson: lesson.serializable_data,
      #   finished_at: lesson_member.updated_at
      # }

      # js_event :lesson_finished, lesson_finished_event_options

      language_member = language.members.find_or_create_by!(user: current_user)
      if language_member.may_finish?
        language_member.finish!

        # language_finished_event_options = {
        #   user: current_user.serializable_data,
        #   language: language.serializable_data,
        #   language_member: language_member.serializable_data,
        #   finished_at: lesson_member.updated_at
        # }

        # js_event :language_finished, language_finished_event_options
      end
    end

    render json: {
      attributes: lesson_exercise_data
    }, status: :ok
  end
end
