# typed: true

require "test_helper"

class Assistants::RunJobTest < ActiveJob::TestCase
  include ActionCable::TestHelper

  setup do
    @user = users(:full)
    @lesson = language_lessons("elixir-variables")
    lesson_member = @lesson.members.find_by!(user: @user)
    @ai_chat = AiChat.create!(user: @user, language_lesson_member: lesson_member)

    sse = <<~SSE
      data: {"choices":[{"delta":{"role":"assistant","content":"Привет"}}]}

      data: {"choices":[{"delta":{"content":", мир"}}]}

      data: {"choices":[{"delta":{},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":2}}

      data: [DONE]

    SSE

    stub_request(:post, "https://api.openai.com/v1/chat/completions")
      .to_return(
        status: 200,
        body: sse,
        headers: { "Content-Type" => "text/event-stream" }
      )
  end

  def test_streams_and_persists_message_authored_by_chat_user
    assert_broadcasts(AssistantChannel.broadcasting_for(@ai_chat), 3) do
      Assistants::RunJob.perform_now(
        ai_chat_id: @ai_chat.id,
        message: "Как решить?",
        user_code: "x = 1",
        output: "ok",
        locale: "ru"
      )
    end

    user_message = @ai_chat.ai_messages.role_user.last

    assert_equal @user.id, user_message.user_id

    assistant_message = @ai_chat.ai_messages.where(role: "assistant").last

    assert_equal "Привет, мир", assistant_message.content
  end
end
