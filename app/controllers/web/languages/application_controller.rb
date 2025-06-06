class Web::Languages::ApplicationController < Web::ApplicationController
  def resource_language
    @resource_language ||= Language.find_by!(slug: params[:language_id])
  end

  def resource_language_landing_page
    @resource_language_landing_page ||= resource_language.landing_pages.published.with_locale.find_by!(main: true)
    # @resource_language_landing_page ||= resource_language.landing_pages.find_by!(locale: I18n.locale)
  end

  # alias_method :resource_language_landing_page, :landing_page
end
