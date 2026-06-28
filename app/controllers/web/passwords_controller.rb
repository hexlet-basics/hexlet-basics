# typed: strict

class Web::PasswordsController < Web::ApplicationController
  allow_unauthenticated_access
  before_action :assert_token_passed
  before_action :set_user_password_form

  rescue_from ActiveSupport::MessageVerifier::InvalidSignature, with: :handle_invalid_token

  sig { returns(T.untyped) }
  def edit
    seo_tags = {
      title: t(".title")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      userPassword: UserPasswordResource.new(@user),
      token: token
    }
  end

  sig { returns(T.untyped) }
  def update
    struct = ApplicationParamsStruct.from_params!(PasswordStruct, params.require(:data))
    result = UserService.update_password(T.must(@user), struct)

    case result
    when Typed::Success
      f(:success)
      redirect_to root_path
    when Typed::Failure
      f(:error)
      redirect_to edit_password_path(token), inertia: { errors: result.error.errors }
    end
  end

  private

  sig { returns(T.untyped) }
  def assert_token_passed
    return if token.present?

    handle_invalid_token
  end

  sig { returns(T.untyped) }
  def set_user_password_form
    @user = T.let(T.unsafe(User).find_by_password_reset_token!(token), T.nilable(User))
  end

  sig { returns(T.untyped) }
  def handle_invalid_token
    f(:error)
    redirect_to root_path
  end

  sig { returns(T.untyped) }
  def token
    params[:token]
  end
end
