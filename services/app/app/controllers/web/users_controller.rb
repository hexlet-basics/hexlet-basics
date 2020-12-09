# frozen_string_literal: true

class Web::UsersController < Web::ApplicationController
  def new
    @user = User::SignUpForm.new
  end

  def create
    @user = User::SignUpForm.new(params[:user_sign_up_form])

    if @user.save
      sign_in @user

      f(:success)
      redirect_to root_path
    else
      render :new
    end
  end
end
