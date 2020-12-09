# frozen_string_literal: true

class Web::PasswordsController < Web::ApplicationController
  def edit
    @user = User::PasswordForm.find_by!(reset_password_token: params[:reset_password_token])
  end

  def update
    @user = User::PasswordForm.find_by!(reset_password_token: params[:reset_password_token])

    if @user.update(params[:user_password_form])
      redirect_to root_path
    else
      render :edit
    end
  end
end
