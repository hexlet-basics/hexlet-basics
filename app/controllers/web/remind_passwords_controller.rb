class Web::RemindPasswordsController < Web::ApplicationController
  def new
    remind_password_form = RemindPasswordForm.new

    seo_tags = {
      title: t(".title")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      passwordReminder: PasswordReminderResource.new(remind_password_form)
    }
  end

  def create
    remind_password_form = RemindPasswordForm.new(params[:remind_password_form])

    if remind_password_form.valid?
      user = remind_password_form.user
      UserService.reset_password!(user, params[:suffix])

      f(:success)
      redirect_to root_path
    else
      f(:error)
      redirect_to_inertia new_remind_password_url, remind_password_form
    end
  end
end
