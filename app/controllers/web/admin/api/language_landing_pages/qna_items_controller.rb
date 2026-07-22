# typed: strict
# frozen_string_literal: true

class Web::Admin::Api::LanguageLandingPages::QnaItemsController < Web::Admin::Api::ApplicationController
  sig { void }
  def index
    resource = Language::LandingPageQnaItemResource.new(language_landing_page.qna_items.order(:id))

    respond_with resource
  end

  sig { void }
  def create
    landing_page = language_landing_page
    qna_item = landing_page.qna_items.build(qna_item_params)
    qna_item.save
    resource = Language::LandingPageQnaItemResource.new(qna_item)

    respond_with resource, location: admin_api_language_landing_page_qna_item_url(landing_page, qna_item)
  end

  sig { void }
  def update
    item = qna_item
    item.update(qna_item_params)
    resource = Language::LandingPageQnaItemResource.new(item)

    respond_with resource
  end

  sig { void }
  def destroy
    item = qna_item
    item.destroy
    resource = Language::LandingPageQnaItemResource.new(item)

    respond_with resource
  end

  private

  sig { returns(Language::LandingPage) }
  def language_landing_page
    Language::LandingPage.find(params.expect(:language_landing_page_id))
  end

  sig { returns(Language::LandingPage::QnaItem) }
  def qna_item
    language_landing_page.qna_items.find(params.expect(:id))
  end

  sig { returns(ActionController::Parameters) }
  def qna_item_params
    params.fetch(:data, {}).permit(:question, :answer)
  end
end
