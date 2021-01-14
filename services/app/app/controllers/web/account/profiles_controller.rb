# frozen_string_literal: true

class Web::Account::ProfilesController < Web::Account::ApplicationController
  before_action :authenticate_user!

  before_action do
    title :base
  end

  def edit
    @form = current_user.becomes(User::ProfileForm)
  end
end
