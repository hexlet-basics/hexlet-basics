# typed: strict
# frozen_string_literal: true

module LocaleConcern
  extend ActiveSupport::Concern
  extend T::Sig
  extend T::Helpers
  requires_ancestor { ApplicationController }

  included do
    T.bind(self, T.class_of(ApplicationController))
    around_action :setup_locale
  end

  sig { void }
  def require_russian_locale
    raise ActionController::RoutingError, "Not Found" unless I18n.locale == :ru
  end

  private

  sig { params(blk: T.proc.void).void }
  def setup_locale(&blk)
    locale = params[:suffix].presence || I18n.default_locale
    I18n.with_locale(locale, &blk)
  end

  sig { void }
  def prepare_locale_settings
    # NOTE: never redirect bots
    return if send(:browser_bot?)

    if view_context.current_page?(root_path) && !params[:suffix]
      remembered_locale = session[:locale].presence
      if remembered_locale
        # root page, no subdomain and no default locale -> redirect
        if remembered_locale != I18n.default_locale.to_s
          redirect_to root_url(suffix: remembered_locale), allow_other_host: true
        end
      else
        # root page, no subdomain, never changed locale
        ru_country_codes = [ "RU" ]
        if locale_from_header == "ru" || ru_country_codes.include?(country_by_ip)
          redirect_to root_url(suffix: :ru), allow_other_host: true
        end
      end
    else
      # not root page or root with subdomain
      session[:locale] = I18n.locale.to_s
    end
  end

  sig { returns(T.nilable(String)) }
  def locale_from_header
    return unless request.env["HTTP_ACCEPT_LANGUAGE"]

    # Extract the preferred locale from the Accept-Language header
    parsed_locales = request.env["HTTP_ACCEPT_LANGUAGE"]
                      .split(",")
                      .map { |l| l.split(";").first }
    # Find the first locale that is available in the app
    parsed_locales.find { |locale| I18n.available_locales.map(&:to_s).include?(locale) }
  end

  sig { returns(String) }
  def country_by_ip
    Geocoder.search(request.remote_ip).first&.country_code || "EN"
  end
end
