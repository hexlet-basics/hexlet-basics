# typed: strict
# frozen_string_literal: true

class Web::MagicLinksController < Web::ApplicationController
  allow_unauthenticated_access
  before_action :redirect_if_authenticated, only: %i[new create]

  sig { void }
  def show
    user = User.find_by_token_for(:magic_link, params[:id])

    if user
      publish_auth_events(user)
      sign_in(user)
      f(:success)
      redirect_to after_authentication_url
    else
      f(:error, type: :alert)
      redirect_to new_magic_link_path
    end
  end

  sig { void }
  def new
    set_meta_tags title: t(".title")
    render inertia: true, props: {}
  end

  sig { void }
  def create
    struct = ApplicationParamsStruct.from_params!(MagicLinkStruct, params.require(:data))
    MagicLinkService.request!(T.must(struct.email), suffix: params[:suffix])

    # Anti-enumeration: the outcome is the same whether or not the email exists.
    f(:success)
    redirect_to root_path
  end

  private

  sig { params(user: User).void }
  def publish_auth_events(user)
    event = UserSignedInEvent.new(data: {
      user_id: user.id, occurrence_count: -1, email: user.email, locale: I18n.locale.to_s
    })
    publish_event(event, user)
    js_event(event)
  end
end
