# frozen_string_literal: true

module RedirectConcern
  extend ActiveSupport::Concern

  def redirect_archived_language
    if landing_page.archived?
      next_landing_page = landing_page.landing_page_to_redirect
      main_landing_page = landing_page.language.landing_pages.find_by(main: true)

      if next_landing_page
        redirect_to view_context.language_path(next_landing_page.slug), status: 301
      elsif main_landing_page
        redirect_to view_context.language_path(main_landing_page.slug), status: 301
      else
        redirect_to view_context.root_path, status: 301
      end
    end
  end
end
