class SurveyHandler
  def call(event)
    # locale = event.metadata[:locale]
    return if I18n.locale != :ru
    user_id = event.metadata[:user_id]
    return unless user_id

    user = User.find(user_id)

    case event
    when UserSignedUpEvent
      surveys = Survey.next_for_user(user)
      surveys.map do |survey|
        Survey.find_or_request_answer_if_needed_by(survey, user)
      end
    when LessonFinishedEvent
      course_slug = event.data.fetch(:course_slug)
      locale = event.data.fetch(:locale)
      course = Language.find_by slug: course_slug
      course_member = course.members.find_by user: user

      surveys = Survey.next_for_user(user, course_member: course_member)
      surveys.map do |survey|
        Survey.find_or_request_answer_if_needed_by(survey, user)
      end

    when BookRequestedEvent
      surveys = Survey.next_for_user(user, ignore_filters: true)
      surveys.map do |survey|
        Survey.find_or_request_answer_if_needed_by(survey, user)
      end
    else
      # nothing to do
    end
  end
end
