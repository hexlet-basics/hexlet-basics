# typed: strict

# Common type for the Google One Tap verifier dependency. Both the real gem
# (`Google::Auth::IDTokens`, wrapped by GoogleOneTapClient) and the test double
# (GoogleAuthStub) conform to it, so `Dependencies#google_one_tap` gets a real
# type and call sites can invoke `.verify_oidc` without `T.unsafe`.
class GoogleOneTapInterface
  extend T::Sig
  extend T::Helpers

  abstract!

  sig do
    abstract.params(
      token: T.untyped,
      aud: T.nilable(String),
      azp: T.nilable(String),
      iss: T.nilable(String)
    ).returns(T::Hash[String, T.untyped])
  end
  def self.verify_oidc(token, aud: nil, azp: nil, iss: nil); end
end
