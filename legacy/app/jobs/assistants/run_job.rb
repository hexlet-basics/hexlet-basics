# typed: false

class Assistants::RunJob < ApplicationJob
  def perform(ai_chat_id:, message:, user_code:, output:, locale:)
    I18n.with_locale(locale) do
      run(ai_chat_id:, message:, user_code:, output:)
    end
  end

  private

  def run(ai_chat_id:, message:, user_code:, output:)
    ai_chat = AiChat.find(ai_chat_id)
    lesson = ai_chat.language_lesson_member.lesson
    lesson_info = lesson.infos.find_by!(locale: I18n.locale)
    language = lesson.language

    community_url = I18n.t("common.community_url")
    answer_language = I18n.t(I18n.locale, scope: "common.languages")

    instructions = <<~PROMPT
      Ты помогаешь изучать #{language.slug} в рамках курса на Hexlet.
      Этот разговор посвящён уроку «#{lesson_info.name}».
      Курс состоит из теории, практики и тестов, которые выполняются прямо в браузере.
      Ты не показываешь решение практики, пользователь должен решить практику самостоятельно.
      Направляй, давай объяснения, помогай разобраться, предлагай шаги для решения, выдвигай гипотезы.
      Отвечай на языке: #{answer_language}. Отвечай коротко.

      Если во время обсуждения обнаружится ошибка в теории или задании (описании, решении, тестах),
      то порекомендуй написать об этом в сообществе #{community_url}. Там сидит команда проекта.

      Теория урока:
      #{lesson_info.theory}

      Задание урока:
      #{lesson_info.instructions}
    PROMPT

    ai_chat.with_instructions(instructions)

    prompt = <<~PROMPT
      Пользовательский код (решение практики): #{user_code}

      Результат запуска пользовательского кода (используя тесты задания): #{output}

      Вопрос пользователя: #{message}
    PROMPT

    message_id = nil
    index = 0

    ai_chat.ask(prompt) do |chunk|
      next if chunk.content.blank?

      message_id ||= ai_chat.ai_messages.where(role: "assistant").order(:id).last&.id

      AssistantChannel.broadcast_to(
        ai_chat,
        delta: [ chunk.content ],
        message_id:,
        index:
      )

      index += 1
    end

    AssistantChannel.broadcast_to(
      ai_chat,
      delta: [ "DONE" ],
      message_id:
    )
  end
end
