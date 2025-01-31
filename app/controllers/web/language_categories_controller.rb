# frozen_string_literal: true

class Web::LanguageCategoriesController < Web::ApplicationController
  def index
    categories = Language::Category.all

    seo_tags = {
      title: t(".header")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      categories: Language::CategoryResource.new(categories)
    }
  end

  def show
    category = Language::Category.find_by! slug: params[:id]
    # @language_members_by_language = current_user.language_members.index_by(&:language_id)
    courses = category.language_version_infos.current
      .includes([ language: :current_version ])
      .with_locale.merge Language.web.ordered
    #
    # infos = Language::Version::Info.where(locale: I18n.locale, language: @languages)
    # infos_by_language = infos.index_by { |item| item.language.id }
    # item_builders = @languages.map { |l| CourseSchema.to_builder(l, infos_by_language.fetch(l.id)) }
    #
    # @builder = ItemListSchema.to_builder(item_builders)
    #
    # @blog_posts = @category.blog_posts.published.limit(3)
    #
    # @switching_locales.each do |locale,|
    #   @switching_locales[locale] = full_url_for(locale: AppHost.locale_for_url(locale))
    # end

    seo_tags = {
      title: t(".header", name: category)
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      courseCategory: Language::CategoryResource.new(category),
      categoryCourses: LanguageResource.new(courses)
    }
  end
end
