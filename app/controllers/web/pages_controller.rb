# frozen_string_literal: true

class Web::PagesController < Web::ApplicationController
  PAGES = %w[about
             cookie_policy
             privacy
             authors
             tos].freeze

  def show
    page = params[:id]
    unless PAGES.include? page
      raise ActionController::RoutingError, "Page not found"
    end

    title = t(".parts.#{params[:id]}.title")

    seo_tags = {
      title:,
      description: t(".parts.#{params[:id]}.meta.description")
    }
    set_meta_tags seo_tags
    # set_meta_tags title: t(@page, scope: 'web.pages')

    render inertia: true, props: {
      page:,
      title:
    }
  end
end
