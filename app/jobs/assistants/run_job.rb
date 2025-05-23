class Assistants::RunJob < ApplicationJob
  def perform(lesson_member_id:, message:, user_code:, output:)
    lesson_member = Language::Lesson::Member.find(lesson_member_id)
    lesson = lesson_member.lesson
    lesson_info = lesson.infos.find_by!(locale: I18n.locale)
    language = lesson.language

    unless language.openai_assistant_id
      throw RuntimeError.new "#{language} without openai_assistant_id"
    end

    # TODO: fix dry-inject
    openai_api = OpenAI::Client.new do |f|
      f.response :logger, Rails.logger, bodies: true
      if configus.hexlet_proxy.url.present?
        if configus.hexlet_proxy.url
          f.proxy = { uri: configus.hexlet_proxy.url }
        end
      end
    end

    # TODO: if guest

    unless lesson_member.openai_thread_id
      thread = openai_api.threads.create
      lesson_member.openai_thread_id = thread["id"]
      lesson_member.save!
    end

    # create message for the thread
    created_message = openai_api.messages.create(
      thread_id: lesson_member.openai_thread_id,
      parameters: {
        role: "user",
        content: [
          {
            type: "text",
            text: "Пользовательский код (решение практики): #{user_code}"
          },
          {
            type: "text",
            text: "Результат запуска пользовательского кода (используя тесты задания): #{output}"
          },
          {
            type: "text",
            text: "Вопрос пользователя: #{message}"
          }
        ]
      }
    )

    instructions = "
      Ты помогаешь изучать #{language.slug} на основе загруженного курса в файлах.
      Этот тред посвящен уроку #{lesson_info.name}.
      Курс состоит из теории, практики и тестов, которые выполняются прямо в браузере.
      Ты не показываешь решение практики, пользователь должен решить практику самостоятельно.
      Направляй, давай объяснения, помогай разобраться, предлагай шаги для решения, выдвигай гипотезы.
      Отвечай на языке: #{I18n.t(I18n.locale, scope: 'common.languages')}. Отвечай коротко.

      Если во время обсуждения обнаружится ошибка в теории или задании (описании, решении, тестах),
      то порекомендуй написать об этом в сообществе https://t.me/HexletLearningBot. Там сидит команда проекта.
      "

    m = lesson_member.messages.build
    m.role = "assistant"
    m.language_lesson = lesson
    m.user = lesson_member.user
    m.language = lesson.language
    deltas = []

    chunk_index = 0

    _run_response = openai_api.runs.create(
      thread_id: lesson_member.openai_thread_id,
      parameters: {
        assistant_id: language.openai_assistant_id,
        instructions:,
        stream: proc do |chunk|
          if chunk["object"] == "thread.message.delta"
            content = chunk.dig("delta", "content") || []
            texts = content.map { it.dig("text", "value") }
            next if texts.blank?
            # response.stream.write("0:#{text.to_json}\n")
            AssistantChannel.broadcast_to(
              lesson_member,
              delta: texts,
              message_id: created_message["id"],
              index: chunk_index
            )

            deltas << texts.join

            chunk_index += 1
          end
        end
      }
    )

    m.body = deltas.join
    m.save!

    AssistantChannel.broadcast_to(
      lesson_member,
      delta: [ "DONE" ],
      message_id: created_message["id"]
    )
  end
end
