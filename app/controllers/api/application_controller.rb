class Api::ApplicationController < ActionController::API
  include ActionController::MimeResponds
  # skip_before_action :verify_authenticity_token
  include AuthConcern
  include LocaleConcern
  include EventConcern

  respond_to :json, :xml
end
