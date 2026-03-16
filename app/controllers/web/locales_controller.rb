# frozen_string_literal: true

class Web::LocalesController < Web::ApplicationController
  skip_before_action :prepare_locale_settings, only: [ :switch ]

  def switch
    locale = params[:new_locale]

    unless I18n.available_locales.include?(locale&.to_sym)
      redirect_back fallback_location: root_path(suffix: AppHost.locale_for_url(I18n.default_locale))
      return
    end

    unless current_user.guest?
      current_user.locale = locale
      current_user.save!
    end

    session[:locale] = locale

    route_params = Rails.application.routes.recognize_path(request.referer)
    redirect_to url_for(**route_params, suffix: AppHost.locale_for_url(locale))
  end
end
