require "test_helper"

class Ai::Lessons::MessagesControllerTest < ActionDispatch::IntegrationTest
  setup do
    @language = languages(:elixir)
    @lesson = language_lessons("elixir-variables")
  end

  test "create" do
    sign_in_as(:full)

    VCR.use_cassette("ai-lessons-messages-create") do
      post ai_lesson_messages_path(@lesson), params: { message: "test" }
    end
    assert_response :success
  end
end
