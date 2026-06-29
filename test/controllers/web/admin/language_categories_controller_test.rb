# typed: true

require "test_helper"

class Web::Admin::LanguageCategoriesControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  def test_index
    get admin_language_categories_url
    assert_response :success
  end

  def test_new
    get new_admin_language_category_url
    assert_response :success
  end

  def test_create
    slug = "web-development"

    params = { data: {
      slug: slug, name: slug, header: slug
    } }
    post admin_language_categories_url, params: params
    assert_response :redirect

    assert { Language::Category.find_by(slug: slug) }
  end

  def test_edit
    language_category = language_categories("frontend-ru")

    get edit_admin_language_category_url(language_category)
    assert_response :success
  end

  def test_update
    language_category = language_categories("frontend-ru")
    landing_page = language_landing_pages("php-ru")

    params = {
      data: {
        name: "new description",
        language_landing_page_ids: [ landing_page.id ]
      }
    }
    patch admin_language_category_url(language_category), params: params
    assert_response :redirect

    assert_equal [ landing_page.id ], language_category.items.reload.pluck(:language_landing_page_id)
  end
end
