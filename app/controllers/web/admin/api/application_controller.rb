# frozen_string_literal: true

class Web::Admin::Api::ApplicationController < ApplicationController
  before_action :require_admin_api_auth!
  respond_to :json
end
