# frozen_string_literal: true

require "test_helper"

class LanguageCategoriesControllerTest < ActionDispatch::IntegrationTest
  def test_should_get_index
    get language_categories_url
    assert_response :success
  end

  def test_should_get_show
    category = language_categories("frontend-ru")
    get language_category_url(id: category.slug)
    assert_response :success
  end
end
