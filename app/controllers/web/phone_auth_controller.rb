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
    struct = ApplicationParamsStruct.from_params(PhoneRequestStruct, params.require(:data))
    return redirect_to(new_phone_auth_path, inertia: { errors: struct.errors }) unless struct.valid?

    result = PhoneAuthService.request_code(T.must(struct.phone), ip: request.remote_ip)

    case result
    when Typed::Failure
      f(:error, type: :alert, values: { reason: result.error })
      redirect_to new_phone_auth_path, inertia: { errors: { phone: phone_error(result.error) } }
    else
      f(:success)
      redirect_to verify_phone_auth_path(phone: result.payload)
    end
  end

  sig { returns(T.untyped) }
  def verify
    phone = params[:phone].to_s
    return redirect_to(new_phone_auth_path) if phone.blank?

    set_meta_tags title: t(".title")
    render inertia: true, props: { phone: }
  end

  sig { returns(T.untyped) }
  def confirm
    struct = ApplicationParamsStruct.from_params(PhoneVerifyStruct, params.require(:data))
    phone = struct.phone.to_s
    return redirect_to(verify_phone_auth_path(phone:), inertia: { errors: struct.errors }) unless struct.valid?

    result = PhoneAuthService.verify_code(phone, T.must(struct.code))

    case result
    when Typed::Failure
      f(:error, type: :alert)
      return redirect_to verify_phone_auth_path(phone:), inertia: { errors: { code: code_error(result.error) } }
    end

    user = result.payload
    publish_auth_events(user)
    sign_in(user)

    f(:success)
    redirect_to after_authentication_url
  end

  private

  sig { params(user: User).void }
  def publish_auth_events(user)
    event = if user.previously_new_record?
      UserSignedUpEvent.new(data: {
        user_id: user.id, email: user.email,
        first_name: user.first_name, last_name: user.last_name, locale: I18n.locale
      })
    else
      UserSignedInEvent.new(data: {
        user_id: user.id, occurrence_count: -1, email: user.email, locale: I18n.locale
      })
    end

    publish_event(event, user)
    js_event(event)
  end

  sig { params(reason: Symbol).returns(String) }
  def phone_error(reason) = t("phone_auth.errors.#{reason}", default: t("phone_auth.errors.invalid_phone"))

  sig { params(reason: Symbol).returns(String) }
  def code_error(reason) = t("phone_auth.errors.#{reason}", default: t("phone_auth.errors.invalid_code"))
end
