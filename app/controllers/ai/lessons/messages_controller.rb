class Ai::Lessons::MessagesController < Ai::ApplicationController
  include ActionController::Live
  # include Import["openai_api"]

  def create
    response.headers["Content-Type"] = "text/event-stream"
    sse = SSE.new(response.stream, retry: 300, event: "empty")

    lesson_info = Language::Lesson.find(params[:lesson_id])
    language = lesson_info.language
    lesson_member = lesson_info.members.find_by!(user: current_user)
    # TODO: fix dry-inject
    openai_api = OpenAI::Client.new

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

    _run_response = openai_api.runs.create(
      thread_id: lesson_member.openai_thread_id,
      parameters: {
        assistant_id: language.openai_assistant_id,
        stream: proc do |chunk|
            # if chunk["object"] == "thread.message.delta"
            # delta = chunk.dig("delta", "content", 0, "text", "value")
            sse.write chunk, event: chunk["object"]
          # sse.write chunk.dig("delta")
          # end
        end
      }
    )

  rescue => e
    # response.stream.write("event: error\ndata: #{e.message}\n\n")
    raise e
  ensure
    sse.close
  end
end
