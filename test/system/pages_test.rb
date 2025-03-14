require "application_system_test_case"

class HomeTest < ApplicationSystemTestCase
  test "/" do
    visit root_path
  end

  test "/ru" do
    visit root_path(suffix: :ru)
  end

  test "/map" do
    visit map_path
  end

  test "/reviews" do
    visit reviews_path(suffix: :ru)
    visit reviews_path(suffix: nil)
  end

  test "/blog_posts" do
    visit blog_posts_path(suffix: :ru)
    visit blog_posts_path(suffix: nil)
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
