# typed: true

module SessionTestHelper
  extend T::Sig

  sig { params(user_or_fixture: T.any(User, Symbol)).returns(User) }
  def sign_in_as(user_or_fixture)
    user = resolve_user(user_or_fixture)

    post session_url, params: { data: { email: T.must(user.email), password: "password" } }
    assert_redirected_to root_path

    user
  end

  sig { void }
  def sign_out
    Current.session&.destroy!
    cookies.delete("session_id")
  end

  sig { returns(T::Boolean) }
  def authenticated?
    cookies["session_id"].present?
  end

  private

    sig { params(user_or_fixture: T.any(User, Symbol)).returns(User) }
    def resolve_user(user_or_fixture)
      return user_or_fixture if user_or_fixture.is_a?(User)

      users(user_or_fixture)
    end
end
