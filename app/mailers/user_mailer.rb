# typed: strict
# frozen_string_literal: true

class UserMailer < ApplicationMailer
  extend T::Sig

  sig { void }
  def reset_password
    user = T.cast(params[:user], User)
    @user = T.let(user, T.nilable(User))
    @suffix = T.let(params[:suffix], T.nilable(String))
    @password_reset_token = T.let(user.password_reset_token, T.untyped)

    I18n.with_locale(@suffix.presence || I18n.default_locale) do
      mail(to: user.email)
    end
  end
end
