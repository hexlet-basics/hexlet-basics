# typed: strict
# frozen_string_literal: true

module RedirectConcern
  extend ActiveSupport::Concern
  extend T::Sig
  extend T::Helpers
  requires_ancestor { ApplicationController }

  sig { returns(T.untyped) }
  def redirect_archived_language
    # landing_page is defined on the concrete controller (Web::LanguagesController),
    # not ApplicationController, so it isn't visible here — reach it via T.unsafe.
    landing_page = send(:landing_page)
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
