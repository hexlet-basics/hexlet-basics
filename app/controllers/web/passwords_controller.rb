class Web::PasswordsController < Web::ApplicationController
  allow_unauthenticated_access
  before_action :assert_token_passed
  before_action :set_user_password_form

  rescue_from ActiveSupport::MessageVerifier::InvalidSignature, with: :handle_invalid_token

  def edit
    seo_tags = {
      title: t(".title")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      userPassword: UserPasswordResource.new(@user_password_form),
      token: token
    }
  end

  def update
    if @user_password_form.update(params[:data])
      f(:success)
      redirect_to root_path
    else
      redirect_to_inertia edit_password_path(token), @user_password_form
    end
  end

  private

  def assert_token_passed
    return if token.present?

    handle_invalid_token
  end

  def set_user_password_form
    @user_password_form = User.find_by_password_reset_token!(token).becomes(User::PasswordForm)
  end

  def handle_invalid_token
    f(:error)
    redirect_to root_path
  end

  def token
    params[:token]
  end
end
