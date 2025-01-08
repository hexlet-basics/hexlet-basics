# frozen_string_literal: true

class Web::RemindPasswordsController < Web::ApplicationController
  def new
    # remind_password_form = RemindPasswordForm.new
    render inertia: true, props: {
      # signInForm: SignInFormResource.new(sign_in_form)
    }
  end

  def create
    remind_password_form = RemindPasswordForm.new(params[:remind_password_form])

    if remind_password_form.valid?
      user = remind_password_form.user
      UserService.reset_password!(user)

      f(:success)
      redirect_to root_path
    else
      # render :new
      redirect_to remind_password_path, inertia: { errors: remind_password_form.errors }
    end
  end
end
