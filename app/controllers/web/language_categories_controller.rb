# frozen_string_literal: true

class Web::LanguageCategoriesController < Web::ApplicationController
  def index; end

  def show
    @category = Language::Category.find_by! slug: params[:id]
    @language_members_by_language = current_user.language_members.index_by(&:language_id)
    @languages = @category.languages.web.ordered

    infos = Language::Version::Info.where(locale: I18n.locale, language: @languages)
    infos_by_language = infos.index_by { |item| item.language.id }
    item_builders = @languages.map { |l| CourseSchema.to_builder(l, infos_by_language.fetch(l.id)) }

    @builder = ItemListSchema.to_builder(item_builders)

    @blog_posts = @category.blog_posts.published.limit(3)
  end
end
