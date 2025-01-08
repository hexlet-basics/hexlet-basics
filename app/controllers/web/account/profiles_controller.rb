# frozen_string_literal: true

class Web::Account::ProfilesController < Web::Account::ApplicationController
  before_action :authenticate_user!

  def edit
    form = current_user.becomes(User::ProfileForm)
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
