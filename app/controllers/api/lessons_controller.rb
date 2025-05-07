class Api::LessonsController < Api::ApplicationController
  def check
    lesson = Language::Lesson.find(params[:id])
    language = lesson.language
    lesson_version = language.current_lesson_versions.find(params[:version_id])
    code = params[:data][:attributes][:code]

    language_version = lesson_version.language_version

    lesson_exercise_data = LessonTester.new.run(lesson_version, language_version, code, current_user)
    passed = lesson_exercise_data[:passed]


    solution_checking_event_data = {
      lesson_slug: lesson.slug,
      course_slug: language.slug,
      locale: I18n.locale,
      passed:
    }

    event = SolutionCheckedEvent.new(data: solution_checking_event_data)
    publish_event(event, current_user)

    lesson_has_been_finished = false
    language_has_been_finished = false

    if passed && !current_user.guest?
      language_member = language.members.find_by!(user: current_user)
      lesson_member = lesson.members.find_by!(user: current_user)

      if lesson_member.may_finish?
        lesson_member.finish!
        lesson_has_been_finished = true

        lesson_finished_event_data = {
          occurrence_count: language_member.lesson_members.count,
          lesson_slug: lesson.slug,
          course_slug: language.slug,
          locale: I18n.locale
        }

        event = LessonFinishedEvent.new(data: lesson_finished_event_data)
        publish_event(event, current_user)
      end

      if language_member.may_finish?
        language_member.finish!
        language_has_been_finished = true

        course_finished_event_data = {
          occurrence_count: current_user.language_members.finished.count,
          slug: language.slug,
          locale: I18n.locale
        }
        event = CourseFinishedEvent.new(data: course_finished_event_data)
        publish_event(event, current_user)
      end
    end

    if passed && current_user.guest?
      session[:finished_as_guest] ||= {}
      session[:finished_as_guest][lesson.id] = true
    end

    response_data = OpenStruct.new({
      **lesson_exercise_data,
      lesson_has_been_finished:,
      language_has_been_finished:
    })

    render json: LessonCheckingResponseResource.new(response_data), status: :ok
  end
end
