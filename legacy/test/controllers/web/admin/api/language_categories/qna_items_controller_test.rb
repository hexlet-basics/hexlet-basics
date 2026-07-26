# typed: true

require "test_helper"

class Web::Admin::Api::LanguageCategories::QnaItemsControllerTest < ActionDispatch::IntegrationTest
  extend T::Sig

  class SetupContext < T::Struct
    const :category, Language::Category
    const :qna_item, Language::Category::QnaItem
  end

  sig { returns(SetupContext) }
  def context
    SetupContext.new(
      category: language_categories("programming-ru"),
      qna_item: language_category_qna_items("programming-ru-1"),
    )
  end

  setup do
    sign_in_as(:admin)
  end

  def test_index
    get admin_api_language_category_qna_items_url(context.category), as: :json

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
    assert_equal "New question", category.qna_items.order(:id).last!.question
  end

  def test_update
    ctx = context

    patch admin_api_language_category_qna_item_url(ctx.category, ctx.qna_item),
      params: { data: { question: "Updated question", answer: "Updated answer" } },
      as: :json

    assert_response :success
    assert_equal "Updated question", ctx.qna_item.reload.question
    assert_equal "Updated answer", ctx.qna_item.answer
  end

  def test_destroy
    ctx = context

    assert_difference -> { ctx.category.qna_items.count }, -1 do
      delete admin_api_language_category_qna_item_url(ctx.category, ctx.qna_item), as: :json
    end

    assert_response :success
  end
end
