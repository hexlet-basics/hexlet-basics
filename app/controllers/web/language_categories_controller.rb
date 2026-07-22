# typed: strict
# frozen_string_literal: true

class Web::LanguageCategoriesController < Web::ApplicationController
  allow_unauthenticated_access

  sig { void }
  def index
    categories = Language::Category.with_locale

    seo_tags = {
      title: t(".header"),
      description: t(".meta.description"),
      canonical: language_categories_url,
      og: {
        title: t(".header"),
        description: t(".meta.description")
      },
      twitter: {
        card: "summary",
        site: t("links.hexlet_twitter_handle")
      }
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      categories: Language::CategoryResource.new(categories)
    }
  end

  sig { void }
  def show
    category = Language::Category.with_locale.find_by! slug: params.expect(:id)
    landing_pages = category.language_landing_pages.web.where(listed: true).merge(Language.ordered)

    title = t(".header", name: category.header)
    description = t(".meta.description", name: category.header)
    seo_tags = {
      title:,
      description:,
      canonical: language_category_url(category.slug),
      og: {
        title:,
        description:
      },
      twitter: {
        card: "summary",
        site: t("links.hexlet_twitter_handle")
      }
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      courseCategory: Language::CategoryResource.new(category),
      qnaItems: Language::CategoryQnaItemResource.new(category.qna_items),
      categoryLandingPages: Language::LandingPageForListsResource.new(landing_pages),
      lead: LeadCrudResource.new(LeadStruct.new)
    }
  end
end
