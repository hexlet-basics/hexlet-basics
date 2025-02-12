# frozen_string_literal: true

class Api::ApplicationController < ActionController::API
  # skip_before_action :verify_authenticity_token
  include AuthConcern
  include LocaleConcern
  include EventConcern

  respond_to :json

  def event_store
    Rails.configuration.event_store
  end
end
