class Web::PasswordsController < Web::ApplicationController
  before_action :assert_reset_token_passed

  def edit
    user_password_form = User::PasswordForm.find_by!(reset_password_token: params[:reset_password_token])

    seo_tags = {
      title: t(".title")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      userPassword: UserPasswordResource.new(user_password_form),
      resetPasswordToken: params[:reset_password_token]
    }
  end

  def update
    user_password_form = User::PasswordForm.find_by!(reset_password_token: params[:reset_password_token])

    if user_password_form.update(params[:user])
      f(:success)
      redirect_to root_path
    else
      # f(:error)
      redirect_to_inertia edit_password_path(reset_password_token: params[:reset_password_token]), user_password_form
    end
  end

  private

  def assert_reset_token_passed
    return if params[:reset_password_token].present?

    redirect_to root_path
  end
end
