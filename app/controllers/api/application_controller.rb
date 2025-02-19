# frozen_string_literal: true

class Api::ApplicationController < ActionController::API
  # skip_before_action :verify_authenticity_token
  include AuthConcern
  include LocaleConcern
  include EventConcern

  before_action :setup_locale

  respond_to :json
end
