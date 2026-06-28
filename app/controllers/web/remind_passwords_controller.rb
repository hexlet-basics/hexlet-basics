# typed: strict

class Web::RemindPasswordsController < Web::ApplicationController
  allow_unauthenticated_access

  sig { returns(T.untyped) }
  def new
    seo_tags = {
      title: t(".title")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {}
  end

  sig { returns(T.untyped) }
  def create
    struct = ApplicationParamsStruct.from_params!(RemindPasswordStruct, params.require(:data))

    user = User.find_by(email: struct.email)

    unless user
      struct.errors.add(:email, :user_does_not_exist)
      f(:error)
      return redirect_to new_remind_password_url, inertia: { errors: struct.errors }
    end

    UserService.reset_password!(user, params[:suffix])

    f(:success)
    redirect_to root_path
  end
end
