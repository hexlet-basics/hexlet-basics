class Ai::ApplicationController < ActionController::Base
  skip_before_action :verify_authenticity_token
  rate_limit to: 10, within: 3.minutes

  # include ActionController::MimeResponds
  # # skip_before_action :verify_authenticity_token
  include AuthConcern
  include LocaleConcern
end
