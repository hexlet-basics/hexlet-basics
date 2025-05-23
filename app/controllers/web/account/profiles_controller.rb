class Web::Account::ProfilesController < Web::Account::ApplicationController
  def edit
    form = current_user.becomes(User::ProfileForm)

    seo_tags = {
      title: t(".title"),
      description: t(".meta.description")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      form: UserProfileFormResource.new(form)
    }
  end

  def update
    form = current_user.becomes(User::ProfileForm)

    if form.update(params[:user])
      f(:success)

      redirect_to edit_account_profile_path
    else
      f(:error)

      redirect_to edit_account_profile, inertia: { errors: form.errors }
    end
  end

  def destroy
    current_user.mark_as_removed!
    sign_out

    f(:success)

    redirect_to root_path
  end
end
