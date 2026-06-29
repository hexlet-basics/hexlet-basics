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
end
