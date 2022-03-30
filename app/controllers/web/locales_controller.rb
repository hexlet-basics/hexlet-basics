# frozen_string_literal: true

class Web::LocalesController < Web::ApplicationController
  skip_before_action :prepare_locale_settings, only: [:switch]

  def switch
    locale = params[:locale]
    # redirect_path = request.referer || root_path

    unless I18n.available_locales.include?(locale&.to_sym)
      redirect_back fallback_location: root_path
      return
    end

    unless current_user.guest?
      current_user.locale = locale
      current_user.save!
    end

    session[:locale] = locale

    sd = locale == 'en' ? nil : locale
    redirect_to root_url(subdomain: sd), allow_other_host: true
  end
end
