class Web::SessionsController < Web::ApplicationController
  before_action :guests_only!, only: [ :new, :create ]

  def new
    sign_in_form = SignInForm.new

    seo_tags = {
      title: t(".title"),
      canonical: new_session_url,
      description: t(".meta.description")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      signInForm: SignInFormResource.new(sign_in_form)
    }
  end

  def create
    sign_in_form = SignInForm.new(params[:user])

    if sign_in_form.valid?
      sign_in(sign_in_form.user)

      f(:success)
      redirect_to root_path
    else
      f(:error)
      redirect_to_inertia new_session_path, sign_in_form
    end
  end

  def destroy
    sign_out
    redirect_to root_path
  end
end
