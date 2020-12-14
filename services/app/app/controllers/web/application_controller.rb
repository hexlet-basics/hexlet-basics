# frozen_string_literal: true

class Web::ApplicationController < ApplicationController
  around_action :switch_locale
  include AuthManagment
  include FlashConcern

  def switch_locale(&action)
    locale = extract_locale_from_subdomain || I18n.default_locale
    I18n.with_locale(locale, &action)
  end

  private

  def extract_locale_from_subdomain
    parsed_locale = request.subdomains.first
    I18n.available_locales.map(&:to_s).include?(parsed_locale) ? parsed_locale : nil
  end

  before_action do
    gon.current_user = {
      id: current_user.id,
      email: current_user.email,
      created_at: current_user.created_at,
      isGuest: current_user.guest?
    }
  end
end
