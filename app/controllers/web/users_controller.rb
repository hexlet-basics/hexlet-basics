# frozen_string_literal: true

class Web::UsersController < Web::ApplicationController
  def new
    @user_form = UserForm.new(User.new)
  end

  def create
    @user_form = UserForm.new(User.new)

    if @user_form.validate(user_params)
      @user_form.save
      sign_in @user_form.model

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
