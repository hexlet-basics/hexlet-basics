# typed: true

require "test_helper"

class Web::Admin::Api::LanguageLandingPages::QnaItemsControllerTest < ActionDispatch::IntegrationTest
  extend T::Sig

  class SetupContext < T::Struct
    const :landing_page, Language::LandingPage
    const :qna_item, Language::LandingPage::QnaItem
  end

  sig { returns(SetupContext) }
  def context
    SetupContext.new(
      landing_page: language_landing_pages("javascript-ru"),
      qna_item: language_landing_page_qna_items("javascript-ru-1"),
    )
  end

  setup do
    sign_in_as(:admin)
  end

  def test_index
    get admin_api_language_landing_page_qna_items_url(context.landing_page), as: :json

    assert_response :success
  end

  def test_create
    landing_page = language_landing_pages("php-ru")

    assert_difference -> { landing_page.qna_items.count }, 1 do
      post admin_api_language_landing_page_qna_items_url(landing_page),
        params: { data: { question: "New question", answer: "New answer" } },
        as: :json
    end

    assert_response :created
    assert_equal "New question", landing_page.qna_items.order(:id).last!.question
  end

  def test_update
    ctx = context

    patch admin_api_language_landing_page_qna_item_url(ctx.landing_page, ctx.qna_item),
      params: { data: { question: "Updated question", answer: "Updated answer" } },
      as: :json

    assert_response :success
    assert_equal "Updated question", ctx.qna_item.reload.question
    assert_equal "Updated answer", ctx.qna_item.answer
  end

  def test_destroy
    ctx = context

    assert_difference -> { ctx.landing_page.qna_items.count }, -1 do
      delete admin_api_language_landing_page_qna_item_url(ctx.landing_page, ctx.qna_item), as: :json
    end

    assert_response :success
  end
end
