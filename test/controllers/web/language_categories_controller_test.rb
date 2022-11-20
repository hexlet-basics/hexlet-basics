require 'test_helper'

class LanguageCategoriesControllerTest < ActionDispatch::IntegrationTest
  test 'should get index' do
    get language_categories_url
    assert_response :success
  end

  test 'should get show' do
    get language_categories_url
    assert_response :success
  end
end
