# frozen_string_literal: true

class Web::ApplicationController < ApplicationController
  include ActionView::Helpers::UrlHelper
  include FlashConcern
  # include TitleConcern
  include EventConcern

  before_action :prepare_locale_settings

  before_action do
    gon.push({
               current_user: {
                 id: current_user.id,
                 email: current_user.email,
                 created_at: current_user.created_at,
                 is_guest: current_user.guest?
               },
               locale: I18n.locale,
               events: EventsMapping.events
             })
  end

  before_action do
    @language_version_infos = Language::Version::Info
                              .where({ locale: I18n.locale })
                              .joins(language_version: :current_language)
  end

  private

  def prepare_locale_settings
    subdomain = request.subdomains.first

    # NOTE: never redirect bots
    # NOTE: Temporary test of OG markup
    return I18n.locale = (subdomain || :en) if browser.bot?

    # TODO: write tests
    if current_page?(root_path) && !subdomain
      remembered_locale = session[:locale].presence
      if remembered_locale
        # root page, no subdomain, changed locale
        if remembered_locale == 'en'
          I18n.locale = :en
        else
          redirect_to root_url(subdomain: remembered_locale), allow_other_host: true
        end
      else
        # root page, no subdomain, never changed locale
        ru_country_codes = ['RU']
        if ru_country_codes.include?(country_by_ip)
          redirect_to root_url(subdomain: 'ru'), allow_other_host: true
        else
          I18n.locale = :en
        end
      end
    else
      # not root page or root with subdomain
      I18n.locale = subdomain || :en
      session[:locale] = I18n.locale
    end
  end

  def country_by_ip
    @country_by_ip ||= Geocoder.search(request.remote_ip).first&.country_code || 'EN'
  end
end
