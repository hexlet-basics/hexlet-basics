# typed: strict

class ReviewLessonJob < ApplicationJob
  extend T::Sig

  sig { params(lesson_info_id: Integer).void }
  def perform(lesson_info_id)
    info = Language::Lesson::Version::Info.find(lesson_info_id)
    review = Language::Lesson::Review.find_or_initialize_by language: info.language,
      lesson: info.lesson,
      locale: info.locale

    messages = T.must(info.lesson).ai_messages.role_user.order(id: :desc).limit(100)

    raw_output = nil
    if !messages.empty?
      instructions = <<~PROMPT
      Проанализируй вопросы, которые задают студенты ассистенту по уроку. Вопросы будут переданы ниже.
      Суммируй основные претензиии и пожелания. Предложи как поменять урок.
      PROMPT

      questions = messages.map(&:content).join("\n\n")

      response = RubyLLM.chat
        .with_instructions(instructions)
        .ask([
          "Урок (теория и упражнение): #{info.theory}\n\n#{info.instructions}",
          "Вопросы пользователей: #{questions}"
        ].join("\n\n"))

      raw_output = response.content
    end

    Rails.logger.info "REVIEW SUMMARY #{raw_output}"

    review.summary = raw_output
    review.lesson_version = info.version
    review.lesson_version_info = info
    review.save!
  end
end
