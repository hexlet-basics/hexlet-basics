# typed: strict

class Web::Account::ProfilesController < Web::Account::ApplicationController
  sig { returns(T.untyped) }
  def edit
    seo_tags = {
      title: t(".title"),
      description: t(".meta.description")
    }
    set_meta_tags seo_tags

    render inertia: true, props: {
      form: UserProfileFormResource.new(T.must(current_user)),
      passkeys: User::CredentialResource.new(T.must(current_user).credentials)
    }
  end

  sig { returns(T.untyped) }
  def update
    struct = ApplicationParamsStruct.from_params(ProfileStruct, params.require(:data))
    result = UserService.update_profile(T.must(current_user), struct)

    case result
    when Typed::Success
      f(:success)
      redirect_to edit_account_profile_path
    when Typed::Failure
      f(:error)
      redirect_to edit_account_profile_path, inertia: { errors: result.error.errors }
    end
  end

  sig { returns(T.untyped) }
  def destroy
    T.must(current_user).mark_as_removed!
    terminate_session

    f(:success)

    redirect_to root_path
  end
end
