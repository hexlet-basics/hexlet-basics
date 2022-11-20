# frozen_string_literal: true

class Web::ApplicationController < ApplicationController
  include ActionView::Helpers::UrlHelper
  include FlashConcern
  # include TitleConcern
  include EventConcern
  include LocaleConcern

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
                              .with_locale
                              .joins(language_version: :current_language)
                              .includes(:language, :language_version)
    @language_categories = Language::Category.all
  end

  private

  def prepare_locale_settings
    # NOTE: never redirect bots
    if browser.bot?
      I18n.locale = params[:locale] || I18n.default_locale
      return
    end

    if current_page?(root_path) && !params[:locale]
      remembered_locale = session[:locale].presence
      if remembered_locale
        # root page, no subdomain and no default locale -> redirect
        if remembered_locale.to_sym != I18n.default_locale
          redirect_to root_url(locale: remembered_locale), allow_other_host: true
        end
      else
        # root page, no subdomain, never changed locale
        ru_country_codes = ['RU']
        if ru_country_codes.include?(country_by_ip)
          redirect_to root_url(locale: :ru), allow_other_host: true
        end
      end
    else
      # not root page or root with subdomain
      session[:locale] = I18n.locale
    end
  end

  def country_by_ip
    @country_by_ip ||= Geocoder.search(request.remote_ip).first&.country_code || 'EN'
  end
end
