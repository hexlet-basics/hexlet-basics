# frozen_string_literal: true

module LocaleConcern
  extend ActiveSupport::Concern

  private

  def setup_locale
    I18n.locale = params[:suffix].presence || I18n.default_locale
  end

  def prepare_locale_settings
    # NOTE: never redirect bots
    setup_locale and return if browser_bot?

    if view_context.current_page?(root_path) && !params[:suffix]
      remembered_locale = session[:locale].presence
      if remembered_locale
        # root page, no subdomain and no default locale -> redirect
        if remembered_locale.to_sym != I18n.default_locale
          redirect_to root_url(suffix: remembered_locale), allow_other_host: true
        end
      else
        # root page, no subdomain, never changed locale
        ru_country_codes = [ "RU" ]
        if locale_from_header == :ru || ru_country_codes.include?(country_by_ip)
          redirect_to root_url(suffix: :ru), allow_other_host: true
        end
      end
    else
      setup_locale
      # not root page or root with subdomain
      session[:locale] = I18n.locale
    end
  end
end
