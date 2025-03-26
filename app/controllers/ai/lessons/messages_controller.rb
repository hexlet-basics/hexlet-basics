class Ai::Lessons::MessagesController < Ai::ApplicationController
  include ActionController::Live
  # include Import["openai_api"]

  def index
    lesson = Language::Lesson.find(params[:lesson_id])
    # language = lesson_info.language
    lesson_member = lesson.members.find_by!(user: current_user)

    messages = []
    if lesson_member.openai_thread_id?
      openai_api = OpenAI::Client.new
      result = openai_api.messages.list(thread_id: lesson_member.openai_thread_id)
      # raise result["data"].inspect
      messages = result["data"].map do
        {
          id: it["id"],
          role: it["role"],
          created_at: it["created_at"],
          content: it["content"].map { it["text"]["value"] }.join
        }
      end.reverse
    end

    render json: messages
  end

  def create
    response.headers["Content-Type"] = "text/event-stream"
    response.headers["Cache-Control"] = "no-cache"

    # Disable nginx buffering
    response.headers["X-Accel-Buffering"] = "no"
    response.headers["X-Content-Type-Options"] = "nosniff"
    response.headers["Transfer-Encoding"] = "chunked"
    response.headers["x-vercel-ai-data-stream"] = "v1"

    lesson = Language::Lesson.find(params[:lesson_id])
    lesson_info = lesson.infos.find_by!(locale: I18n.locale)
    language = lesson.language
    lesson_member = lesson.members.find_by!(user: current_user)
    # TODO: fix dry-inject
    openai_api = OpenAI::Client.new log_errors: true

    # TODO: if guest

    unless lesson_member.openai_thread_id
      thread = openai_api.threads.create
      lesson_member.openai_thread_id = thread["id"]
      lesson_member.save!
    end
    # raise lesson_member.inspect

    # create message for the thread
    _created_message = openai_api.messages.create(
      thread_id: lesson_member.openai_thread_id,
      parameters: {
        role: "user",
        content: params[:message]
      }
    )

    response.stream.write("4:#{{
      id: SecureRandom.uuid,
      role: "assistant",
      content: [
        {
          type: "text",
          text: { value: "" }
        }
      ]
    }.to_json}\n")

    instructions = "
      Ты помогаешь изучать #{language.slug} на основе загруженного курса в файлах.
      Этот тред посвящен уроку #{lesson_info.name}.lesson
      Курс состоит из теории, практики и тестов, которые выполняются прямо в браузере.
      Ты не показываешь решение практики, пользователь должен решить практику самостоятельно.
      Направляй, давай объяснения, помогай разобраться, предлагай шаги для решения, выдвигай гипотезы.
      Отвечай на языке: #{I18n.t(I18n.locale, scope: 'common.languages')}
      "

    _run_response = openai_api.runs.create(
      thread_id: lesson_member.openai_thread_id,
      parameters: {
        assistant_id: language.openai_assistant_id,
        instructions:,
        stream: proc do |chunk|
          # https://sdk.vercel.ai/docs/ai-sdk-ui/stream-protocol
          if chunk["object"] == "thread.message.delta"
            content = chunk.dig("delta", "content")
            text = content.map { it.dig("text", "value") }.join
            response.stream.write("0:#{text.to_json}\n")
          end
        end
      }
    )

    # response.stream.write("d:#{{
    #   finishReason: "stop",
    #   usage: {
    #     promptTokens: 10,
    #     completionTokens: 20
    #   }
    # }.to_json}\n")

  ensure
    response.stream.close
  end
end
