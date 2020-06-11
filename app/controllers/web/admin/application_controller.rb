# frozen_string_literal: true

class Web::Admin::ApplicationController < Web::ApplicationController
  before_action :authenticate_admin!

  def authenticate_admin!
    return redirect_to root_path unless current_user.admin?
  end
end
