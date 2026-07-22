# typed: true

require "test_helper"

class Ai::Lessons::MessagesControllerTest < ActionDispatch::IntegrationTest
  setup do
    @language = languages(:elixir)
    @lesson = language_lessons("elixir-variables")
  end

  def test_create
    sign_in_as(:full)

    original_queue_adapter = ActiveJob::Base.queue_adapter
    ActiveJob::Base.queue_adapter = :test

    assert_enqueued_jobs 1, only: Assistants::RunJob do
      post ai_lesson_messages_path(@lesson), params: { message: "test" }
    end

    assert_response :success
  ensure
    ActiveJob::Base.queue_adapter = original_queue_adapter
  end

  def test_index_returns_history
    user = sign_in_as(:full)
    lesson_member = @lesson.members.find_by!(user:)
    ai_chat = AiChat.create!(user:, language_lesson_member: lesson_member)
    ai_chat.ai_messages.create!(role: "user", content: "question", user:)
    ai_chat.ai_messages.create!(role: "assistant", content: "answer")

    get ai_lesson_messages_path(@lesson)

    assert_response :success
    body = response.parsed_body

    assert_equal 2, body.size
    assert_equal(%w[user assistant], body.pluck("role"))
    assert_equal(%w[question answer], body.pluck("content"))
  end
end
