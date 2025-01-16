# frozen_string_literal: true

class Web::UsersController < Web::ApplicationController
  def new
    user = User::SignUpForm.new

    seo_tags = {
      title: t(".title"),
      description: t(".meta.description"),
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      user: UserResource.new(user)
    }
  end

  def create
    user = User::SignUpForm.new(params[:user_sign_up_form])

    if user.save
      sign_in user
      # js_event_options = {
      #   user: @user
      # }
      # js_event :signed_up, js_event_options

      f(:success)
      redirect_to root_url
    else
      redirect_to_inertia new_user_url, user
    end
  end
end
