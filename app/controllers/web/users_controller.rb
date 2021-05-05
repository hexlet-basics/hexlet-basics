# frozen_string_literal: true

class Web::UsersController < Web::ApplicationController
  before_action do
    title :base
  end

  def new
    @user = User::SignUpForm.new
  end

  def create
    @user = User::SignUpForm.new(params[:user_sign_up_form])

    if @user.save
      sign_in @user
      js_event :signed_up

      f(:success)
      redirect_to root_path
    else
      render :new
    end
  end
end
