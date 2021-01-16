# frozen_string_literal: true

class Web::Admin::HomeController < Web::Admin::ApplicationController
  def index
    @q = User.admin.ransack(params[:q])
    @users = @q.result
  end
end
