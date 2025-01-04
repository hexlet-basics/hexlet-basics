# frozen_string_literal: true

class Web::ApplicationController < ApplicationController
  allow_browser versions: :modern
  inertia_share flash: -> { flash.to_hash }

  include ActionView::Helpers::UrlHelper
  include FlashConcern
  # include TitleConcern
  include EventConcern
  include LocaleConcern

  before_action :prepare_locale_settings

  # before_action do
  #   gon.push({
  #              current_user: {
  #                id: current_user.id,
  #                email: current_user.email,
  #                created_at: current_user.created_at,
  #                is_guest: current_user.guest?
  #              },
  #              locale: I18n.locale,
  #              events: EventsMapping.events
  #            })
  # end

  before_action do
    # @language_menu_data = build_language_menu_data
  end

  inertia_share do
    language_categories = Language::Category.all
    languages_infos = Language::Version::Info.with_locale

    {
      languageCategories: Language::CategoryResource.new(language_categories),
      courses: LanguageResource.new(languages_infos),
      locale: I18n.locale,
      suffix: I18n.locale == :en ? nil : I18n.locale,
      auth: {
        user: UserResource.new(current_user)
      }
    }
  end

  # before_action do
  #   @switching_locales ||= I18n.available_locales.without(I18n.locale)
  #                              .index_with { |locale| switch_locale_path(new_locale: locale) }
  # end

  private

  def build_language_menu_data
    languages_multicolumns_treshold = 10

    scope = Language::Version::Info
            .with_locale
            .ordered
            .includes(:language_version)
            .joins(:language, language_version: :current_language)
            .preload(:language)

    completed_language_version_infos = scope.completed
    incompleted_language_version_infos = scope.incompleted
    languages_count = scope.size

    # NOTE: Если мало языков нет смысла разбивать на несколько колонок
    columns_count = languages_count < languages_multicolumns_treshold ? 1 : 2

    {
      completed: completed_language_version_infos,
      incompleted: incompleted_language_version_infos,
      columns_count: columns_count
    }
  end

  def prepare_locale_settings
    # NOTE: never redirect bots
    if browser.bot?
      I18n.locale = params[:suffix] || I18n.default_locale
      return
    end

    if current_page?(root_path) && !params[:suffix]
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
      I18n.locale = params[:suffix].presence || I18n.default_locale
      # not root page or root with subdomain
      session[:locale] = I18n.locale
    end
  end

  def country_by_ip
    @country_by_ip ||= Geocoder.search(request.remote_ip).first&.country_code || "EN"
  end

  def locale_from_header
    return unless request.env["HTTP_ACCEPT_LANGUAGE"]

    # Extract the preferred locale from the Accept-Language header
    parsed_locales = request.env["HTTP_ACCEPT_LANGUAGE"]
                      .split(",")
                      .map { |l| l.split(";").first }
    # Find the first locale that is available in the app
    parsed_locales.find { |locale| I18n.available_locales.map(&:to_s).include?(locale) }
  end
end
