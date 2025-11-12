require "application_system_test_case"

class HomeTest < ApplicationSystemTestCase
  def test_root_path
    visit root_path(suffix: :ru)
    visit root_path(suffix: nil)
  end

  def test_users_new
    visit new_user_path(suffix: :ru)
    visit new_user_path(suffix: nil)
  end

  def test_remind_password_new
    visit remind_password_path(suffix: :ru)
    visit remind_password_path(suffix: nil)
  end

  def test_session_new
    visit new_session_path(suffix: :ru)
    visit new_session_path(suffix: nil)
  end

  # test "/map" do
  #   visit map_path
  # end

  def test_cases
    visit cases_path(suffix: :ru)
    visit cases_path(suffix: nil)
  end

  def test_cases_for_teachers
    visit for_teachers_cases_path(suffix: :ru)
    visit for_teachers_cases_path(suffix: nil)
  end

  def test_reviews
    visit reviews_path(suffix: :ru)
    visit reviews_path(suffix: nil)
  end

  def test_blog_posts
    visit blog_posts_path(suffix: :ru)
    visit blog_posts_path(suffix: nil)
  end

  def test_pages_id
    visit page_path("about", suffix: :ru)
    visit page_path("about", suffix: nil)
  end

  def test_blog_posts_id
    post = blog_posts("from-full")

    visit blog_post_path(post, suffix: :ru)
    visit blog_post_path(post, suffix: nil)
  end

  def test_language_categories
    visit language_categories_path(suffix: :ru)
    visit language_categories_path(suffix: nil)
  end

  def test_language_categories_id
    category = language_categories("programming-ru")

    visit language_category_path(category, suffix: :ru)
    visit language_category_path(category, suffix: nil)
  end

  def test_languages_id
    landing_page = language_landing_pages("javascript-ru")
    visit language_path(landing_page.slug, suffix: :ru)
    visit language_path(landing_page.slug, suffix: nil)
  end

  def test_languages_language_id_lessons_id
    lesson = language_lessons("javascript-module1-lesson2")

    visit language_lesson_path(lesson.language, lesson, suffix: :ru)
    visit language_lesson_path(lesson.language, lesson, suffix: nil)
  end
end
