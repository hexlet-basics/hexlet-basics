class SurveyHandler
  def call(event)
    # locale = event.metadata[:locale]
    return if I18n.locale != :ru
    user_id = event.metadata[:user_id]
    return unless user_id

    user = User.find(user_id)

    # Rails.logger.info("!!!!!! #{event}")
    case event
    when UserSignedUpEvent
      Survey.find_or_request_answer_if_needed_by("goal", user)
      Survey.find_or_request_answer_if_needed_by("coding-experience", user)
      # raise survey_answer.inspect
    when LessonFinishedEvent
      Rails.logger.info("LessonFinishedEvent")
      course_slug = event.data.fetch(:course_slug)
      locale = event.data.fetch(:locale)
      course = Language.find_by slug: course_slug
      course_member = course.members.find_by user: user
      should_request_answer = course_member.lesson_members.size > 2
      if should_request_answer
        Survey.find_or_request_answer_if_needed_by("study-plan", user)
      end
      # Если завершил 3 урока, то запускаем опрос для ответа (хочу сменить профессию)

      # Есть ли у вас готовый план обучения?
      # Да, у меня есть план и я его придерживаюсь
      # У меня есть примерное представление, того что надо делать
      # Я выбрал направление, но у меня нет плана обучения
      # Я все еще выбираю то, чем заниматься и/или на каком языке писать

      # Насколько вам нужна подержка в процессе смены профессии?
      # Думаю справлюсь со всем сам
      # Научиться программировать могу сам, но мне нужна помощь с трудоустройством
      # Идеально если бы было к кому обратиться за помощью
      # Пока не знаю

      # Если завершил 7 уроков, то запускаем опрос для ответа (хочу сменить профессию)

      # Сколько у вас есть времени на обучение в неделю?
    else
      # nothing to do
    end
  end
end
