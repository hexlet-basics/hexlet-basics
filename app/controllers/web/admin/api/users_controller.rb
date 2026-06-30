# typed: strict
# frozen_string_literal: true

class Web::Admin::Api::UsersController < Web::Admin::Api::ApplicationController
  sig { returns(T.untyped) }
  def search
    term = params[:term].presence || params[:q].presence

    if term.blank?
      render json: []
      return
    end

    search = User.ransack(first_name_or_last_name_or_email_cont: term)
    users = search.result.limit(20)

    options = users.map do |user|
      { value: user.id.to_s, label: [ user.to_s, user.email ].compact.join(" — ") }
    end

    render json: options
  end
end
