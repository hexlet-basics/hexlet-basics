require "test_helper"

class Ai::Lessons::MessagesControllerTest < ActionDispatch::IntegrationTest
  setup do
    @language = languages(:ruby)
    @lesson = @language.lessons.first
  end

  test "create" do
    sign_in_as(:one)

    post ai_lesson_messages_path(@lesson), params: { message: "test" }
    assert_response :success
  end
end
