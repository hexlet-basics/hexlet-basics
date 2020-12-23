# frozen_string_literal: true

class Web::LocalesController < Web::ApplicationController

  def switch
    locale = params[:locale]
    redirect_path = request.referrer || root_path

    if !I18n.available_locales.include?(locale&.to_sym)
      redirect_back fallback_location: redirect_path
      return
    end

    if !current_user.guest?
      current_user.locale = locale
      current_user.save!
    end

    session[:locale] = locale

    redirect_to redirect_path
  end

end
