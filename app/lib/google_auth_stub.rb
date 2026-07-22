# typed: strict
# frozen_string_literal: true

class GoogleAuthStub < GoogleOneTapInterface
  extend T::Sig

  # rubocop:disable Lint/UnusedMethodArgument
  sig do
    override.params(
      _token: T.untyped,
      aud: T.nilable(String),
      azp: T.nilable(String),
      iss: T.nilable(String)
    ).returns(T::Hash[String, T.untyped])
  end
  def self.verify_oidc(_token, aud: nil, azp: nil, iss: nil)
    {
      "email" => "example@mail.com",
      "sub" => "1144422255589998888"
    }
  end
  # rubocop:enable Lint/UnusedMethodArgument
end
