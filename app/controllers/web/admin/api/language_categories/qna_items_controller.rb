# typed: true
# frozen_string_literal: true

class Web::Admin::Api::LanguageCategories::QnaItemsController < Web::Admin::Api::ApplicationController
  before_action :set_language_category
  before_action :set_qna_item, only: %i[update destroy]

  def index
    resource = Language::CategoryQnaItemResource.new(@language_category.qna_items.order(:id))

    respond_with resource
  end

  def create
    qna_item = @language_category.qna_items.build(qna_item_params)
    qna_item.save
    resource = Language::CategoryQnaItemResource.new(qna_item)

    respond_with resource, location: admin_api_language_category_qna_item_url(@language_category, qna_item)
  end

  def update
    @qna_item.update(qna_item_params)
    resource = Language::CategoryQnaItemResource.new(@qna_item)

    respond_with resource
  end

  def destroy
    @qna_item.destroy
    resource = Language::CategoryQnaItemResource.new(@qna_item)

    respond_with resource
  end

  private

  def set_language_category
    @language_category = Language::Category.find(params[:language_category_id])
  end

  def set_qna_item
    @qna_item = @language_category.qna_items.find(params[:id])
  end

  def qna_item_params
    params.fetch(:data, {}).permit(:question, :answer)
  end
end
