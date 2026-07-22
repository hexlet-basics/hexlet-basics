# typed: true
# frozen_string_literal: true

require "application_system_test_case"

class HomeTest < ApplicationSystemTestCase
  def test_root_path
    visit root_path(suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
    visit root_path(suffix: nil)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_users_new
    visit new_user_path(suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
    visit new_user_path(suffix: nil)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_remind_password_new
    visit new_remind_password_path(suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
    visit new_remind_password_path(suffix: nil)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_session_new
    visit new_session_path(suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
    visit new_session_path(suffix: nil)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_cases
    visit cases_path(suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_cases_for_teachers
    visit for_teachers_cases_path(suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_reviews
    visit reviews_path(suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
    visit reviews_path(suffix: nil)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_blog_posts
    visit blog_posts_path(suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
    visit blog_posts_path(suffix: nil)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_pages_id
    visit page_path("about", suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
    visit page_path("about", suffix: nil)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_blog_posts_id
    post = blog_posts("from-full")

    visit blog_post_path(post.slug, suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_language_categories
    visit language_categories_path(suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
    visit language_categories_path(suffix: nil)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_language_categories_id
    visit language_category_path(language_categories("programming-ru").slug, suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
    visit language_category_path(language_categories("programming-en").slug, suffix: nil)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_languages_id
    visit language_path(language_landing_pages("javascript-ru").slug, suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
    visit language_path(language_landing_pages("php-en").slug, suffix: nil)

    assert_equal 200, page.status_code
    assert_selector "body"
  end

  def test_languages_language_id_lessons_id
    lesson = language_lessons("javascript-module1-lesson2")

    visit language_lesson_path(lesson.language.slug, lesson.slug, suffix: :ru)

    assert_equal 200, page.status_code
    assert_selector "body"
  end
end
