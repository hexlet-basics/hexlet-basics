class Web::PasswordsController < Web::ApplicationController
  before_action :assert_reset_token_passed

  def edit
    user_password_form = User::PasswordForm.find_by!(reset_password_token: params[:reset_password_token])

    render inertia: true, props: {
      userPassword: UserPasswordResource.new(user_password_form)
    }
  end

  def update
    user_password_form = User::PasswordForm.find_by!(reset_password_token: params[:reset_password_token])

    if user_password_form.update(params[:user_password_form])
      f(:success)
      redirect_to root_path
    else
      # f(:error, now: true)
      redirect_to_inertia edit_password_path(reset_password_token: params[:reset_password_token]), user_password_form
    end
  end

  private

  def assert_reset_token_passed
    return if params[:reset_password_token].present?

    redirect_to root_path
  end
end
