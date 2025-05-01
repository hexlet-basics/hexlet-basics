class SurveyHandler
  def call(event)
    # locale = event.metadata[:locale]
    return if I18n.locale != :ru
    user_id = event.metadata[:user_id]
    return unless user_id

    user = User.find(user_id)

    case event
    when UserSignedUpEvent
      Survey.find_or_request_answer_if_needed_by("goal", user)
      Survey.find_or_request_answer_if_needed_by("coding-experience", user)
    when LessonFinishedEvent
      course_slug = event.data.fetch(:course_slug)
      locale = event.data.fetch(:locale)
      course = Language.find_by slug: course_slug
      course_member = course.members.find_by user: user
      should_request_answer1 = course_member.lesson_members.size > 2
      if should_request_answer1
        Survey.find_or_request_answer_if_needed_by("career-change-reason", user)
        Survey.find_or_request_answer_if_needed_by("career-change-study-plan", user)
      end

      should_request_answer2 = course_member.lesson_members.size > 5
      if should_request_answer2
        Survey.find_or_request_answer_if_needed_by("career-change-barrier", user)
        Survey.find_or_request_answer_if_needed_by("career-change-time-commitment", user)
      end

      should_request_answer3 = course_member.lesson_members.size > 8
      if should_request_answer3
        Survey.find_or_request_answer_if_needed_by("career-change-priority", user)
        Survey.find_or_request_answer_if_needed_by("career-change-preferred-intro-format", user)
        Survey.find_or_request_answer_if_needed_by("career-change-contact-method", user)
      end
    else
      # nothing to do
    end
  end
end
