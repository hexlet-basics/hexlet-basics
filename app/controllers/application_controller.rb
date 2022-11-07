# frozen_string_literal: true

require 'application_responder'

# frozen_string_literal: true

class ApplicationController < ActionController::Base
  self.responder = ApplicationResponder
  respond_to :html

  include AuthConcern
  content_security_policy Rails.env.production?

  def default_url_options
    { locale: AppHost.locale_for_url(I18n.locale) }
  end
end
