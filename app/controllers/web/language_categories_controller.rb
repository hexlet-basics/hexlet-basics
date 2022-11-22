# frozen_string_literal: true

class Web::LanguageCategoriesController < Web::ApplicationController
  def index; end

  def show
    @category = Language::Category.find_by! slug: params[:id]
    @language_members_by_language = current_user.language_members.index_by(&:language_id)
    @languages = @category.languages
  end
end
