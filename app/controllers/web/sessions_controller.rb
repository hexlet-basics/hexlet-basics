# typed: true

class Web::SessionsController < Web::ApplicationController
  allow_unauthenticated_access only: [ :new, :create ]
  before_action :redirect_if_authenticated, only: [ :new, :create ]

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
    sign_in_form = SignInForm.new(params[:data])

    if sign_in_form.valid?
      user = sign_in_form.user
      sign_in(user)

      data = {
        user_id: user.id,
        occurrence_count: -1,
        email: T.must(user.email),
        locale: I18n.locale
      }
      event = UserSignedInEvent.new(data:)
      EventSender.publish_event(event, user)

      f(:success)
      redirect_to after_authentication_url
    else
      f(:error)
      redirect_to_inertia new_session_path, sign_in_form
    end
  end

  def destroy
    terminate_session
    redirect_to root_path
  end
end
