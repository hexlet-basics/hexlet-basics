# typed: strict
# frozen_string_literal: true

class Web::PhoneAuthController < Web::ApplicationController
  allow_unauthenticated_access
  before_action :redirect_if_authenticated

  sig { returns(T.untyped) }
  def new
    set_meta_tags title: t(".title")
    render inertia: true, props: {}
  end

  sig { returns(T.untyped) }
  def create
    struct = ApplicationParamsStruct.from_params!(PhoneRequestStruct, params.require(:data))
    result = PhoneAuthService.request_code(T.must(struct.phone), ip: request.remote_ip)

    case result
    when Typed::Success
      f(:success)
      redirect_to verify_phone_auth_path(phone: result.payload)
    when Typed::Failure
      f(:error, type: :alert)
      redirect_to new_phone_auth_path, inertia: { errors: { phone: t("phone_auth.errors.#{result.error}") } }
    end
  end

  sig { returns(T.untyped) }
  def verify
    set_meta_tags title: t(".title")
    render inertia: true, props: { phone: params[:phone].to_s }
  end

  sig { returns(T.untyped) }
  def confirm
    struct = ApplicationParamsStruct.from_params!(PhoneVerifyStruct, params.require(:data))
    phone = T.must(struct.phone)
    result = PhoneAuthService.verify_code(phone, T.must(struct.code))

    case result
    when Typed::Success
      user = result.payload
      publish_auth_events(user)
      sign_in(user)
      f(:success)
      redirect_to after_authentication_url
    when Typed::Failure
      f(:error, type: :alert)
      redirect_to verify_phone_auth_path(phone:), inertia: { errors: { code: t("phone_auth.errors.#{result.error}") } }
    end
  end

  private

  sig { params(user: User).void }
  def publish_auth_events(user)
    event = if user.previously_new_record?
      UserSignedUpEvent.new(data: {
        user_id: user.id, email: user.email,
        first_name: user.first_name, last_name: user.last_name, locale: I18n.locale.to_s
      })
    else
      UserSignedInEvent.new(data: {
        user_id: user.id, occurrence_count: -1, email: user.email, locale: I18n.locale.to_s
      })
    end

    publish_event(event, user)
    js_event(event)
  end
end
