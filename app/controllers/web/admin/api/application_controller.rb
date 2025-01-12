# frozen_string_literal: true

class Web::Admin::Api::ApplicationController < ActionController::API
  include AuthConcern
  before_action :require_admin_api_auth!
  respond_to :json
end
