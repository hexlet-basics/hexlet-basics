class Web::UsersController < Web::ApplicationController
  before_action :guests_only!

  def new
    user = User::SignUpForm.new

    seo_tags = {
      title: t(".title"),
      canonical: view_context.new_user_url,
      description: t(".meta.description")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      demo: params[:demo],
      user: UserSignUpFormResource.new(user)
    }
  end

  def create
    user = User::SignUpForm.new(params[:user])

    begin
      user.save!

      sign_up(user)
      sign_in(user)
      fill_guests_data(user)

      f(:success)
      redirect_to params[:from].presence || root_path
    rescue ActiveRecord::RecordInvalid, ActiveRecord::RecordNotUnique => _e
      f(:error)
      redirect_to_inertia new_user_url, user
      # raise user.errors.full_messages.inspect
    end
  end
end
