class Assistants::RunJob < ApplicationJob
  # queue_as :assistant

  def perform(lesson_member_id:, message:)
    lesson_member = Language::Lesson::Member.find(lesson_member_id)
    lesson = lesson_member.lesson
    lesson_info = lesson.infos.find_by!(locale: I18n.locale)
    language = lesson.language

    # TODO: fix dry-inject
    openai_api = OpenAI::Client.new do |f|
      f.response :logger, Rails.logger, bodies: true
      if configus.hexlet_proxy.url.present?
        f.proxy = { uri: configus.hexlet_proxy.url }
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
        content: message
      }
    )

    instructions = "
      Ты помогаешь изучать #{language.slug} на основе загруженного курса в файлах.
      Этот тред посвящен уроку #{lesson_info.name}.
      Курс состоит из теории, практики и тестов, которые выполняются прямо в браузере.
      Ты не показываешь решение практики, пользователь должен решить практику самостоятельно.
      Направляй, давай объяснения, помогай разобраться, предлагай шаги для решения, выдвигай гипотезы.
      Отвечай на языке: #{I18n.t(I18n.locale, scope: 'common.languages')}. Отвечай коротко.
      "

    chunk_index = 0

    _run_response = openai_api.runs.create(
      thread_id: lesson_member.openai_thread_id,
      parameters: {
        assistant_id: language.openai_assistant_id,
        instructions:,
        stream: proc do |chunk|
          # https://sdk.vercel.ai/docs/ai-sdk-ui/stream-protocol
          if chunk["object"] == "thread.message.delta"
            content = chunk.dig("delta", "content") || []
            text = content.map { it.dig("text", "value") }.join
            next if text.blank?
            # response.stream.write("0:#{text.to_json}\n")
            AssistantChannel.broadcast_to(
              lesson_member,
              delta: text,
              message_id: created_message["id"],
              index: chunk_index
            )

            chunk_index += 1
          end
        end
      }
    )

    AssistantChannel.broadcast_to(
      lesson_member,
      delta: "[DONE]",
      message_id: created_message["id"]
    )
  end
end
