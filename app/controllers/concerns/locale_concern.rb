# frozen_string_literal: true

module LocaleConcern
  extend ActiveSupport::Concern

  included do
    around_action :switch_locale
  end

  # NOTE: for en locale dont use path /en but /
  def switch_locale(&action)
    extracted_locale = extract_locale_from_subdomain
    if extracted_locale
      locale = AppHost.locale_for_url(extracted_locale)
      redirect_to full_url_for(locale: locale, subdomain: nil), allow_other_host: true, status: :moved_permanently
      return
    end

    if params[:locale]&.to_sym == I18n.default_locale
      redirect_to full_url_for(locale: nil), status: :moved_permanently
      return
    end

    if params[:locale].present?
      I18n.with_locale(params[:locale], &action)
    else
      I18n.with_locale(I18n.default_locale, &action)
    end
  end

  def extract_locale_from_subdomain
    parsed_locale = request.subdomains.first
    I18n.available_locales.map(&:to_s).include?(parsed_locale) ? parsed_locale : nil
  end
end
