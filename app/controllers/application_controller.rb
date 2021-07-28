# frozen_string_literal: true

class ApplicationController < ActionController::Base
  include AuthConcern
  content_security_policy Rails.env.production?
end
