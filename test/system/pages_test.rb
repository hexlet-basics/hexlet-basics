require "application_system_test_case"

class HomeTest < ApplicationSystemTestCase
  test "/" do
    visit root_path(suffix: :ru)
    visit root_path(suffix: nil)
  end

  test "/users/new" do
    visit new_user_path(suffix: :ru)
    visit new_user_path(suffix: nil)
  end

  test "/remind_password/new" do
    visit remind_password_path(suffix: :ru)
    visit remind_password_path(suffix: nil)
  end

  test "/session/new" do
    visit new_session_path(suffix: :ru)
    visit new_session_path(suffix: nil)
  end

  # test "/map" do
  #   visit map_path
  # end

  test "/cases" do
    visit cases_path(suffix: :ru)
    visit cases_path(suffix: nil)
  end

  test "/cases/for_teachers" do
    visit for_teachers_cases_path(suffix: :ru)
    visit for_teachers_cases_path(suffix: nil)
  end

  test "/reviews" do
    visit reviews_path(suffix: :ru)
    visit reviews_path(suffix: nil)
  end

  test "/blog_posts" do
    visit blog_posts_path(suffix: :ru)
    visit blog_posts_path(suffix: nil)
  end

  test "/pages/:id" do
    visit page_path("about", suffix: :ru)
    visit page_path("about", suffix: nil)
  end

  test "/blog_posts/:id" do
    post = blog_posts("from-full")

    visit blog_post_path(post, suffix: :ru)
    visit blog_post_path(post, suffix: nil)
  end

  test "/language_categories" do
    visit language_categories_path(suffix: :ru)
    visit language_categories_path(suffix: nil)
  end

  test "/language_categories/:id" do
    category = language_categories("programming-ru")

    visit language_category_path(category, suffix: :ru)
    visit language_category_path(category, suffix: nil)
  end

  test "/languages/:id" do
    landing_page = language_landing_pages("javascript-ru")
    visit language_path(landing_page.slug, suffix: :ru)
    visit language_path(landing_page.slug, suffix: nil)
  end

  test "/languages/:language_id/lessons/:id" do
    lesson = language_lessons("javascript-module1-lesson2")

    visit language_lesson_path(lesson.language, lesson, suffix: :ru)
    visit language_lesson_path(lesson.language, lesson, suffix: nil)
  end
end
