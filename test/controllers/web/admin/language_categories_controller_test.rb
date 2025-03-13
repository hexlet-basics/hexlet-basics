require "test_helper"

class Web::Admin::LanguageCategoriesControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  test "index" do
    get admin_language_categories_url
    assert_response :success
  end

  test "new" do
    get new_admin_language_category_url
    assert_response :success
  end

  test "create" do
    slug = "racket"

    params = { language_category: { slug: slug } }
    post admin_language_categories_url, params: params
    assert_response :redirect

    # assert { Language.find_by(slug: slug) }
  end

  test "edit" do
    language_category = language_categories("php-ru")

    get edit_admin_language_category_url(language_category)
    assert_response :success
  end

  test "update" do
    language_category = language_categories("php-ru")

    params = { language_category: { name: "new description" } }
    patch admin_language_category_url(language_category), params: params
    assert_response :redirect

    # language.reload
    # assert { language_category.progress.in_development? }
  end
end
