# frozen_string_literal: true

class Web::LanguageCategoriesController < Web::ApplicationController
  def index
    categories = Language::Category.with_locale

    seo_tags = {
      title: t(".header"),
      description: t(".meta.description"),
      canonical: language_categories_url
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      categories: Language::CategoryResource.new(categories)
    }
  end

  def show
    category = Language::Category.with_locale.find_by! slug: params[:id]
    landing_pages = category.language_landing_pages.web.where(listed: true).merge(Language.ordered)

    seo_tags = {
      title: t(".header", name: category),
      description: t(".meta.category_descriptions.#{category.slug}"),
      canonical: language_category_url(category.slug)
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      courseCategory: Language::CategoryResource.new(category),
      qnaItems: Language::CategoryQnaItemResource.new(category.qna_items),
      categoryLandingPages: Language::LandingPageForListsResource.new(landing_pages),
      lead: LeadCrudResource.new(LeadForm.new)
    }
  end
end
