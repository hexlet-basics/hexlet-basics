# frozen_string_literal: true

class UserMailer < ApplicationMailer
  def reset_password
    @user = params[:user]
    @suffix = params[:suffix]
    @password_reset_token = @user.password_reset_token

    I18n.with_locale(@suffix.presence || I18n.default_locale) do
      mail(to: @user.email)
    end
  end
end
