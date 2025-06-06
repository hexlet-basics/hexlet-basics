class Web::Languages::ApplicationController < Web::ApplicationController
  def resource_language
    @resource_language ||= Language.find_by!(slug: params[:language_id])
  end

  def resource_language_landing_page
    @resource_language_landing_page ||= resource_language.landing_pages.web.find_by(main: true)
  end

  # alias_method :resource_language_landing_page, :landing_page
end
