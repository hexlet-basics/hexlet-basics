# frozen_string_literal: true

class Web::ApplicationController < ApplicationController
  include AuthManagment
  include FlashConcern
  include TitleConcern

  before_action do
    locale = request.subdomains.first || current_user.locale || :en

    if locale == :ru && request.subdomains.empty?
      redirect_to url_for(params.merge(subdomain: locale, only_path: false).permit!)
    end

    locale = request.subdomains.first || :en

    I18n.locale = locale
  end

  before_action do
    gon.current_user = {
      id: current_user.id,
      email: current_user.email,
      created_at: current_user.created_at,
      is_guest: current_user.guest?
    }

    gon.locale = I18n.locale
  end
end
