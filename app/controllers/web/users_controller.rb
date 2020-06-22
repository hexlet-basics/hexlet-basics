# frozen_string_literal: true

class Web::UsersController < Web::ApplicationController
  def new
    @user = User::SignUpType.new
  end

  def create
    @user = User::SignUpType.new(user_params)

    if @user.save
      sign_in @user
      f(:success)

      redirect_to root_path
    else
      render :new
    end
  end

  private

  def user_params
    params.require(:user).permit(:email, :password)
  end
end
