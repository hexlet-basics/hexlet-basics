# frozen_string_literal: true

class Web::ApplicationController < ApplicationController
  include BrowserConcern

  allow_modern_browsers
  inertia_share flash: -> { flash.to_hash }

  include FlashConcern
  include EventConcern
  include LocaleConcern

  before_action :prepare_locale_settings

  inertia_share do
    language_categories = Language::Category.all
    languages_infos = Language::Version::Info
      .current
      .with_locale
      .includes([ language: :current_version ])

    {
      courseCategories: Language::CategoryResource.new(language_categories),
      courses: LanguageResource.new(languages_infos),
      locale: I18n.locale,
      suffix: I18n.locale == :en ? nil : I18n.locale,
      auth: {
        user: UserResource.new(current_user)
      },
      mobileBrowser: mobile_browser?
    }
  end

  before_action do
    gon.current_user = UserResource.new(current_user)
  end

  # before_action do
  #   @switching_locales ||= I18n.available_locales.without(I18n.locale)
  #                              .index_with { |locale| switch_locale_path(new_locale: locale) }
  # end

  private

  def ransack_params(defaults)
    raw = params.permit(:sf, :so, fields: {}).with_defaults({ fields: {} }).with_defaults(defaults)
    ransack = raw["fields"]

    if raw.key?("sf")
      ransack["s"] = "#{raw["sf"]} #{raw["so"] == "1" ? 'asc' : 'desc'}"
    end

    ransack.to_unsafe_hash
  end

  def grid_params(pagy = nil)
    result = params.permit(:sf, :so, fields: {}).with_defaults({ page: 1 })
    if pagy
      result[:tr] = pagy.count()
      result[:per] = pagy.limit()
    end

    OpenStruct.new(result)
  end

  def redirect_to_inertia(url, model)
      redirect_to url, inertia: { errors: model.errors }
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
