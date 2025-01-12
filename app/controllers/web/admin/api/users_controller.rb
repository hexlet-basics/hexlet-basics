# frozen_string_literal: true

class Web::Admin::Api::UsersController < Web::Admin::Api::ApplicationController
  def search
    if !params[:query]
      render json: []
      return
    end

    search = User.ransack(first_name_or_last_name_or_email_cont: params[:query])
    users = search.result.limit(20)

    render json: UserResource.new(users)
  end
end
