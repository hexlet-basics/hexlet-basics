class ApplicationController < ActionController::Base
  # respond_to :html
  include Pagy::Backend

  include AuthConcern
  # content_security_policy Rails.env.production?

  # def default_url_options
  #   { locale: AppHost.locale_for_url(I18n.locale) }
  # end
end
