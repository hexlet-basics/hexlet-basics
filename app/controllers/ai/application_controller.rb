class Ai::ApplicationController < ActionController::Base
  skip_before_action :verify_authenticity_token
  # include ActionController::MimeResponds
  # # skip_before_action :verify_authenticity_token
  include AuthConcern
  include LocaleConcern
end
