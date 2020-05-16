# frozen_string_literal: true

class Web::SessionsController < Web::ApplicationController
  def new
    @user = User.new
  end

  def create
    @user = User.find_by(email: user_params[:email])

    if @user&.authenticate(user_params[:password])
      sign_in @user
      redirect_to root_path
    else
      redirect_to new_sessions_path
    end
  end

  def destroy
    sign_out
    redirect_to root_path
  end

  private

  def user_params
    params.require(:user).permit(:email, :password)
  end
end
