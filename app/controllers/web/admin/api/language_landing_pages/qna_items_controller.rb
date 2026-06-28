# typed: true
# frozen_string_literal: true

class Web::Admin::Api::LanguageLandingPages::QnaItemsController < Web::Admin::Api::ApplicationController
  before_action :set_language_landing_page
  before_action :set_qna_item, only: %i[update destroy]

  def index
    resource = Language::LandingPageQnaItemResource.new(@language_landing_page.qna_items.order(:id))

    respond_with resource
  end

  def create
    qna_item = @language_landing_page.qna_items.build(qna_item_params)
    qna_item.save
    resource = Language::LandingPageQnaItemResource.new(qna_item)

    respond_with resource, location: admin_api_language_landing_page_qna_item_url(@language_landing_page, qna_item)
  end

  def update
    @qna_item.update(qna_item_params)
    resource = Language::LandingPageQnaItemResource.new(@qna_item)

    respond_with resource
  end

  def destroy
    @qna_item.destroy
    resource = Language::LandingPageQnaItemResource.new(@qna_item)

    respond_with resource
  end

  private

  def set_language_landing_page
    @language_landing_page = Language::LandingPage.find(params[:language_landing_page_id])
  end

  def set_qna_item
    @qna_item = @language_landing_page.qna_items.find(params[:id])
  end

  def qna_item_params
    params.fetch(:data, {}).permit(:question, :answer)
  end
end
