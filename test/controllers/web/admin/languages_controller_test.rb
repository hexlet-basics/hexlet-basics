# typed: true
# frozen_string_literal: true

require "test_helper"

class Web::Admin::LanguagesControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  def test_index
    get admin_languages_url
    assert_response :success
  end

  def test_new
    get new_admin_language_url
    assert_response :success
  end

  def test_review
    language = languages(:php)
    infos_count = language.current_lesson_infos.count

    original_queue_adapter = ActiveJob::Base.queue_adapter
    ActiveJob::Base.queue_adapter = :test

    assert_enqueued_jobs infos_count, only: ReviewLessonJob do
      post review_admin_language_url(language.id)
    end

    assert_response :redirect
  ensure
    ActiveJob::Base.queue_adapter = original_queue_adapter
  end

  def test_create
    slug = "racket"

    params = { data: { slug: slug } }
    post admin_languages_url, params: params
    assert_response :redirect

    assert { Language.find_by(slug: slug) }
  end

  def test_create_with_frontend_shaped_payload
    slug = "ai"

    params = {
      data: {
        progress: "draft",
        learn_as: "first_language",
        cover: { name: "cover", record: { id: nil, slug: nil, progress: "draft" } },
        id: nil,
        slug: slug,
        hexlet_program_landing_page: "ai-for-developers",
        repository_url: "https://github.com/hexlet-basics/exercises-"
      }
    }
    post admin_languages_url, params: params
    assert_response :redirect

    language = Language.find_by!(slug: slug)
    assert { !language.cover.attached? }
  end

  def test_create_with_cover
    slug = "racket"

    params = { data: { slug: slug, cover: fixture_file_upload("course-cover.png", "image/png") } }
    post admin_languages_url, params: params
    assert_response :redirect

    language = Language.find_by!(slug: slug)
    assert { language.cover.attached? }
  end

  def test_update_with_frontend_shaped_payload
    language = languages(:php)

    params = {
      data: {
        progress: "in_development",
        cover: { name: "cover", record: { id: language.id, slug: language.slug } }
      }
    }
    patch admin_language_url(language), params: params
    assert_response :redirect

    language.reload
    assert { language.in_development_progress? }
  end

  def test_edit
    language = languages(:php)

    get edit_admin_language_url(language)
    assert_response :success
  end

  def test_update
    language = languages(:php)

    params = { data: { progress: "in_development" } }
    patch admin_language_url(language), params: params
    assert_response :redirect

    language.reload
    assert { language.in_development_progress? }
  end
end
