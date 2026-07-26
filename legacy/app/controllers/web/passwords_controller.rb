# typed: strict

class Web::PasswordsController < Web::ApplicationController
  allow_unauthenticated_access
  before_action :assert_token_passed

  rescue_from ActiveSupport::MessageVerifier::InvalidSignature, with: :handle_invalid_token

  sig { void }
  def edit
    seo_tags = {
      title: t(".title")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      userPassword: UserPasswordResource.new(user),
      token: token
    }
  end

  sig { void }
  def update
    struct = ApplicationParamsStruct.from_params!(PasswordStruct, params.require(:data))
    result = UserService.update_password(T.must(user), struct)

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

  sig { void }
  def assert_token_passed
    return if token.present?

    handle_invalid_token
  end

  sig { returns(T.nilable(User)) }
  def user
    @user ||= T.let(User.find_by_password_reset_token!(token), T.nilable(User))
  end

  sig { params(_exception: T.untyped).void }
  def handle_invalid_token(_exception = nil)
    f(:error)
    redirect_to root_path
  end

  sig { returns(T.nilable(String)) }
  def token
    params[:token]
  end
end
