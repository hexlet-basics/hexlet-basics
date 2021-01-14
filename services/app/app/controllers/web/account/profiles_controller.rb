# frozen_string_literal: true

class Web::Account::ProfilesController < Web::Account::ApplicationController
  before_action :authenticate_user!

  before_action do
    title :base
  end

  def edit
    @form = current_user.becomes(User::ProfileForm)
  end

  def destroy
    current_user.mark_as_removed!
    sign_out

    f(:success)

    redirect_to root_path
  end
end
