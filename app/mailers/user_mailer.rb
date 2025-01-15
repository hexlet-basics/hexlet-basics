# frozen_string_literal: true

class UserMailer < ApplicationMailer
  def reset_password
    @user = params[:user]
    @suffix = params[:suffix]
    mail(to: @user.email)
  end
end
