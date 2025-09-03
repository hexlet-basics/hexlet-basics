class ReviewLessonJob < ApplicationJob
  def perform(lesson_info_id)
    info = Language::Lesson::Version::Info.find(lesson_info_id)
    review = Language::Lesson::Review.find_or_initialize_by language: info.language,
      lesson: info.lesson,
      locale: info.locale

    messages = info.lesson.messages.user_role.order(id: :desc).limit(100)

    raw_output = nil
    if !messages.empty?
      openai_api = DepsLocator.current.openai_api

      instructions = <<~PROMPT
      Проанализируй вопросы, которые задают студенты ассистенту по уроку. Вопросы будут переданы ниже.
      Суммируй основные претензиии и пожелания. Предложи как поменять урок.
      PROMPT

      chat_completion = openai_api.chat(
        parameters: {
          model: :"gpt-4.1",
          messages: [
            { role: "system", content: instructions },
            { role: "user", content: "Урок (теория и упражнение): #{info.theory}\n\n#{info.instructions}" },
            { role: "user", content: "Вопросы пользователей: #{messages.join("\n\n")}" }
          ]
        }
      )

      raw_output = chat_completion.dig("choices", 0, "message", "content")
    end

    Rails.logger.info "REVIEW SUMMARY #{raw_output}"

    review.summary = raw_output
    review.lesson_version = info.version
    review.lesson_version_info = info
    review.save!
  end
end
