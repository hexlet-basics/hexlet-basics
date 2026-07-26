# typed: true

require "test_helper"

class Web::Admin::LanguageLessonsControllerTest < ActionDispatch::IntegrationTest
  def test_review
    lesson = language_lessons(:one)
    infos_count = lesson.infos.count

    original_queue_adapter = ActiveJob::Base.queue_adapter
    ActiveJob::Base.queue_adapter = :test

    assert_enqueued_jobs infos_count, only: ReviewLessonJob do
      post review_admin_language_lesson_url(lesson.id)
    end

    assert_response :redirect
  ensure
    ActiveJob::Base.queue_adapter = original_queue_adapter
  end
end
