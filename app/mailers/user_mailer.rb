# frozen_string_literal: true

class UserMailer < ApplicationMailer
  def reset_password(user_id)
    @user = User.find(user_id)

    mail(to: @user.email)
  end
end
