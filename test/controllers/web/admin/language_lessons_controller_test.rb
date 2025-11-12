require "test_helper"

class Web::Admin::LanguageLessonsControllerTest < ActionDispatch::IntegrationTest
  def test_review
    lesson = language_lessons(:one)

    VCR.use_cassette("ai-lessons-reviews-create") do
      post review_admin_language_lesson_url(lesson.id)
    end
    assert_response :redirect
  end
end
