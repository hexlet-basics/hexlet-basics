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

    VCR.use_cassette("ai-lessons-reviews-create") do
      post review_admin_language_url(language.id)
    end
    assert_response :redirect
  end

  def test_create
    slug = "racket"

    params = { language: { slug: slug } }
    post admin_languages_url, params: params
    assert_response :redirect

    assert { Language.find_by(slug: slug) }
  end

  def test_edit
    language = languages(:php)

    get edit_admin_language_url(language)
    assert_response :success
  end

  def test_update
    language = languages(:php)

    params = { language: { progress: "in_development" } }
    patch admin_language_url(language), params: params
    assert_response :redirect

    language.reload
    assert { language.in_development_progress? }
  end
end
