# typed: true

class Web::SessionsController < Web::ApplicationController
  allow_unauthenticated_access only: [ :new, :create ]
  before_action :redirect_if_authenticated, only: [ :new, :create ]

  def new
    seo_tags = {
      title: t(".title"),
      canonical: new_session_url,
      description: t(".meta.description")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {}
  end

  def create
    struct = ApplicationParamsStruct.from_params!(SignInStruct, params.require(:data))

    email = struct.email.to_s.strip.downcase
    user = User.authenticate_by(email:, password: struct.password)&.then do |authenticated_user|
      authenticated_user if authenticated_user.active?
    end

    unless user
      struct.errors.add(:password, :cannot_sign_in)
      f(:error)
      return redirect_to new_session_path, inertia: { errors: struct.errors }
    end

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
  end

  def destroy
    terminate_session
    redirect_to root_path
  end
end
