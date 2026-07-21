# typed: strict
# frozen_string_literal: true

class Web::Admin::Api::LanguageCategories::QnaItemsController < Web::Admin::Api::ApplicationController
  sig { void }
  def index
    resource = Language::CategoryQnaItemResource.new(language_category.qna_items.order(:id))

    respond_with resource
  end

  sig { void }
  def create
    category = language_category
    qna_item = category.qna_items.build(qna_item_params)
    qna_item.save
    resource = Language::CategoryQnaItemResource.new(qna_item)

    respond_with resource, location: admin_api_language_category_qna_item_url(category, qna_item)
  end

  sig { void }
  def update
    item = qna_item
    item.update(qna_item_params)
    resource = Language::CategoryQnaItemResource.new(item)

    respond_with resource
  end

  sig { void }
  def destroy
    item = qna_item
    item.destroy
    resource = Language::CategoryQnaItemResource.new(item)

    respond_with resource
  end

  private

  sig { returns(Language::Category) }
  def language_category
    Language::Category.find(params[:language_category_id])
  end

  sig { returns(Language::Category::QnaItem) }
  def qna_item
    language_category.qna_items.find(params[:id])
  end

  sig { returns(ActionController::Parameters) }
  def qna_item_params
    params.fetch(:data, {}).permit(:question, :answer)
  end
end
