class Api::ApplicationController < ActionController::API
  include ActionController::MimeResponds
  # skip_before_action :verify_authenticity_token
  include AuthConcern
  include LocaleConcern
  include EventConcern
  include AhoyConcern

  respond_to :json, :xml
end
