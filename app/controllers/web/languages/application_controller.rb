# frozen_string_literal: true

class Web::Languages::ApplicationController < Web::ApplicationController
  # helper_method :resource_language

  def resource_language
    @resource_language ||= Language.find_by!(slug: params[:language_id])
  end

  def resource_language_landing_page
    @resource_language_landing_page ||= resource_language.landing_pages.find_by!(locale: I18n.locale)
  end
end
