# typed: true

# `has_secure_password ..., reset_token: {...}` makes ActiveModel define the
# `find_by_password_reset_token`/`!` class methods (see
# ActiveModel::SecurePassword), which tapioca's DSL compiler doesn't emit.
# Declare them here so token lookups stay typed without `T.unsafe`.
class User
  class << self
    sig { params(token: T.nilable(String)).returns(T.nilable(User)) }
    def find_by_password_reset_token(token); end

    sig { params(token: T.nilable(String)).returns(User) }
    def find_by_password_reset_token!(token); end
  end

  sig { returns(String) }
  def password_reset_token; end
end
