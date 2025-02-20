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
    language_info = language_version.infos.find_by!(locale: I18n.locale)

    lesson_exercise_data = LessonTester.new.run(lesson_version, language_version, code, current_user)

    lesson_has_been_finished = false
    language_has_been_finished = false

    if lesson_exercise_data[:passed] && !current_user.guest?
      lesson_member = lesson.members.find_by!(user: current_user)

      if lesson_member.may_finish?
        lesson_member.finish!
        lesson_has_been_finished = true

        event_data = {
          lesson_slug: lesson.slug,
          course_slug: language.slug,
          locale: language_info.locale
        }

        event = LessonFinishedEvent.new(data: event_data)
        publish_event(event, current_user)
      end

      language_member = language.members.find_by!(user: current_user)

      if language_member.may_finish?
        language_member.finish!
        language_has_been_finished = true

        event_data = {
          slug: language.slug,
          locale: language_info.locale
        }
        event = CourseFinishedEvent.new(data: event_data)
        publish_event(event, current_user)
      end
    end

    event_data = {
      lesson_slug: lesson.slug,
      course_slug: language.slug,
      locale: language_info.locale,
      passed: lesson_exercise_data[:passed]
    }
    event = SolutionCheckedEvent.new(data: event_data)
    publish_event(event, current_user)

    response_data = OpenStruct.new({
      **lesson_exercise_data,
      lesson_has_been_finished:,
      language_has_been_finished:
    })

    render json: LessonCheckingResponseResource.new(response_data), status: :ok
  end
end
