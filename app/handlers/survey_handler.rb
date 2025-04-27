class SurveyHandler
  def call(event)
    # locale = event.metadata[:locale]
    return if I18n.locale != :ru
    user_id = event.metadata[:user_id]
    return unless user_id

    user = User.find(user_id)

    case event
    when UserSignedUpEvent
      _survey_answer1 = Survey.find_or_request_answer_by("goal", user)
      _survey_answer2 = Survey.find_or_request_answer_by("coding-experience", user)
      # raise survey_answer.inspect
    else
      # nothing to do
    end
  end
end
