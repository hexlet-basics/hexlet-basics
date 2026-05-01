require "test_helper"

class Web::Admin::Api::LanguageCategories::QnaItemsControllerTest < ActionDispatch::IntegrationTest
  setup do
    sign_in_as(:admin)
  end

  def test_index
    category = language_categories("programming-ru")

    get admin_api_language_category_qna_items_url(category), as: :json

    assert_response :success
  end

  def test_create
    category = language_categories("frontend-ru")

    assert_difference -> { category.qna_items.count }, 1 do
      post admin_api_language_category_qna_items_url(category),
        params: { data: { question: "New question", answer: "New answer" } },
        as: :json
    end

    assert_response :created
    assert_equal "New question", category.qna_items.order(:id).last.question
  end

  def test_update
    category = language_categories("programming-ru")
    qna_item = language_category_qna_items("programming-ru-1")

    patch admin_api_language_category_qna_item_url(category, qna_item),
      params: { data: { question: "Updated question", answer: "Updated answer" } },
      as: :json

    assert_response :success
    assert_equal "Updated question", qna_item.reload.question
    assert_equal "Updated answer", qna_item.answer
  end

  def test_destroy
    category = language_categories("programming-ru")
    qna_item = language_category_qna_items("programming-ru-1")

    assert_difference -> { category.qna_items.count }, -1 do
      delete admin_api_language_category_qna_item_url(category, qna_item), as: :json
    end

    assert_response :success
  end
end
