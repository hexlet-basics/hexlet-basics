# frozen_string_literal: true

class Web::Admin::Api::UsersController < Web::Admin::Api::ApplicationController
  def search
    search = User.ransack(first_name_or_last_name_or_email_cont: params[:term])
    users = search.result.limit(20).map do |user|
      UserSerializer.to_select2(user)
    end

    respond_with results: users, format: :json
  end
end
