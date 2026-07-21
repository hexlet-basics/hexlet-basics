# typed: strict
# frozen_string_literal: true

class Web::LocalesController < Web::ApplicationController
  allow_unauthenticated_access
  skip_before_action :prepare_locale_settings, only: [ :switch ]

  sig { void }
  def switch
    locale = params[:new_locale]

    unless I18n.available_locales.map(&:to_s).include?(locale)
      redirect_back fallback_location: root_path(suffix: AppHost.locale_for_url(I18n.default_locale))
      return
    end

    if current_user.present?
      T.must(current_user).locale = locale
      T.must(current_user).save!
    end

    session[:locale] = locale

    redirect_to localized_referer_path(locale)
  end

  private

  # Keep the user on the page they came from, swapping the locale prefix in the
  # path. Falls back to the localized root when there is no usable referer.
  sig { params(locale: String).returns(String) }
  def localized_referer_path(locale)
    suffix = AppHost.locale_for_url(locale)
    fallback = root_path(suffix: suffix)

    referer = request.referer
    return fallback if referer.blank?

    uri =
      begin
        URI.parse(referer)
      rescue URI::InvalidURIError
        return fallback
      end
    return fallback unless uri.host == request.host

    stripped = T.must(uri.path).sub(%r{\A/(?:es|ru)(?=/|\z)}, "")
    return fallback if stripped.blank? || stripped == "/"

    prefix = suffix ? "/#{suffix}" : ""
    query = uri.query.present? ? "?#{uri.query}" : ""
    "#{prefix}#{stripped}#{query}"
  end
end
