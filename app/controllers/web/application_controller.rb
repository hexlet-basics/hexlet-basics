# frozen_string_literal: true

class Web::ApplicationController < ApplicationController
  include ActionView::Helpers::UrlHelper
  include FlashConcern
  # include TitleConcern
  include EventConcern

  before_action :prepare_locale_settings

  before_action do
    gon.current_user = {
      id: current_user.id,
      email: current_user.email,
      created_at: current_user.created_at,
      is_guest: current_user.guest?
    }

    gon.locale = I18n.locale
    gon.events = EventsMapping.events
  end

  private

  def prepare_locale_settings
    subdomain = request.subdomains.first
    I18n.locale = subdomain || :en

    # NOTE: боты не должны попадать на автоматический редирект
    unless browser.bot?
      subdomains = {
        ru: 'ru'
      }

      ru_country_codes = ['RU']
      remembered_locale = session[:locale]&.to_sym

      if current_page?(root_path)
        if remembered_locale && remembered_locale != I18n.locale
          url = root_url(subdomain: subdomains.fetch(remembered_locale, ''))
          redirect_to url
        elsif !remembered_locale && !subdomain && ru_country_codes.include?(country_by_ip)
          url = root_url(subdomain: 'ru')
          redirect_to url
        elsif !subdomain && ru_country_codes.exclude?(country_by_ip)
          # Говорим о том, что в английском пока не очень много контента и возможно вы хотели русский
          # f(:, now: true)
        end
      end
    end
  end

  def country_by_ip
    @country_by_ip ||= Geocoder.search(request.remote_ip).first&.country || 'RU'
  end
end
