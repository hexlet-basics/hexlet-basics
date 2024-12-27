# frozen_string_literal: true

class Web::LocalesController < Web::ApplicationController
  skip_before_action :prepare_locale_settings, only: [ :switch ]

  def switch
    locale = params[:new_locale]
    # redirect_path = requst.referer || root_path

    unless I18n.available_locales.include?(locale&.to_sym)
      redirect_back fallback_location: root_path(locale: AppHost.locale_for_url(I18n.default_locale))
      return
    end

    unless current_user.guest?
      current_user.locale = locale
      current_user.save!
    end

    session[:locale] = locale

    redirect_to root_url(locale: AppHost.locale_for_url(locale)), allow_other_host: true
  end
end
