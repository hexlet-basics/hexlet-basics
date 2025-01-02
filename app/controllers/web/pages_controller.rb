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

    render inertia: true, props: {
      page:
    }

    # set_meta_tags title: t(@page, scope: 'web.pages')
  end
end
