# frozen_string_literal: true

class Web::Admin::ApplicationController < Web::ApplicationController
  include AuthManagment

  before_action :authenticate_admin!
end
