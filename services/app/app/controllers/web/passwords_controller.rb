class Web::PasswordsController < Web::ApplicationController
  def new
    @password_form = PasswordForm.new(User.new)
  end

  def create
    @password_form = PasswordForm.new(User.new)

    if @password_form.validate(password_params)
      user = @password_form.user
      UserService.reset_password(user)

      redirect_to root_path
    else
      render :new
    end
  end

  private

  def password_params
    params.require(:password).permit(:email)
  end
end
