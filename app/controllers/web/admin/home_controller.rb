# frozen_string_literal: true

class Web::Admin::HomeController < Web::Admin::ApplicationController
  def index
    @q = User.admin.ransack(params[:q])
    @q.sorts = 'created_at desc' if @q.sorts.empty?
    @users = @q.result
  end
end
