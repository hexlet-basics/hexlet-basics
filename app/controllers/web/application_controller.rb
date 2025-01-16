# frozen_string_literal: true

class Web::ApplicationController < ApplicationController
  # include ActionView::Helpers::UrlHelper

  allow_browser versions: :modern
  inertia_share flash: -> { flash.to_hash }

  # https://inertia-rails.dev/guide/error-handling
  # rescue_from StandardError, with: :inertia_error_page
  #
  # def inertia_error_page(exception)
  #   raise exception if Rails.env.local?
  #
  #   status = ActionDispatch::ExceptionWrapper.new(nil, exception).status_code
  #
  #   render inertia: "error", props: { status: }, status:
  # end

  # include ActionView::Helpers::UrlHelper
  include FlashConcern
  # include TitleConcern
  # include EventConcern
  include LocaleConcern

  before_action :prepare_locale_settings

  inertia_share do
    language_categories = Language::Category.all
    languages_infos = Language::Version::Info.with_locale.includes([ language: :current_version ])

    {
      courseCategories: Language::CategoryResource.new(language_categories),
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

  def default_url_options
    { suffix: params[:suffix] }
  end

  def prepare_locale_settings
    # NOTE: never redirect bots
    if browser.bot?
      I18n.locale = params[:suffix] || I18n.default_locale
      return
    end

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
      I18n.locale = params[:suffix].presence || I18n.default_locale
      # not root page or root with subdomain
      session[:locale] = I18n.locale
    end
  end

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
