# frozen_string_literal: true

class Web::RemindPasswordsController < Web::ApplicationController
  def new
    @remind_password_form = RemindPasswordForm.new
  end

  def create
    @remind_password_form = RemindPasswordForm.new(params[:remind_password_form])

    if @remind_password_form.valid?
      user = @remind_password_form.user
      UserService.reset_password!(user)

      f(:success)
      redirect_to root_path
    else
      render :new
    end
  end
end
