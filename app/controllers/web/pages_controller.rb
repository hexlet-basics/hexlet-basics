# frozen_string_literal: true

class Web::PagesController < Web::ApplicationController
  PAGES = %w[about
             cookie_policy
             privacy
             authors
             tos].freeze

  DISALLOWED_PAGES = %w[cookie_policy privacy tos].freeze

  def show
    page = params[:id]
    unless PAGES.include? page
      raise ActionController::RoutingError, "Page not found"
    end

    title = t(".parts.#{params[:id]}.title")

    description = t(".parts.#{params[:id]}.meta.description")
    seo_tags = {
      title:,
      description:,
      canonical: page_url(params[:id]),
      og: {
        title:,
        description:
      },
      twitter: {
        card: "summary",
        site: "@hexlethq"
      }
    }
    if DISALLOWED_PAGES.exclude? page
      seo_tags[:canonical] = page_url(page)
    end
    set_meta_tags seo_tags
    # set_meta_tags title: t(@page, scope: 'web.pages')

    render inertia: true, props: {
      page:,
      title:
    }
  end
end
