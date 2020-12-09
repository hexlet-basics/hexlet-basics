# frozen_string_literal: true

class Web::RemindPasswordsController < Web::ApplicationController
  def new
    @remind_password_form = RemindPasswordForm.new(User.new)
  end

  def create
    @remind_password_form = RemindPasswordForm.new(User.new)

    if @remind_password_form.validate(remind_password_params)
      user = @remind_password_form.user
      UserService.reset_password(user)

      redirect_to root_path
    else
      render :new
    end
  end

  private

  def remind_password_params
    params.require(:remind_password).permit(:email)
  end
end
