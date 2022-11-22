# frozen_string_literal: true

require 'test_helper'

class LanguageCategoriesControllerTest < ActionDispatch::IntegrationTest
  test 'should get index' do
    get language_categories_url
    assert_response :success
  end

  test 'should get show' do
    category = language_categories(:programming)
    get language_category_url(id: category.slug)
    assert_response :success
  end
end
