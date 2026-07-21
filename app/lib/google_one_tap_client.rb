# typed: strict

# Production adapter: delegates to the googleauth gem. Exists so the gem module
# (which can't subclass GoogleOneTapInterface) is reachable through the typed
# dependency contract.
class GoogleOneTapClient < GoogleOneTapInterface
  extend T::Sig

  sig do
    override.params(
      token: T.untyped,
      aud: T.nilable(String),
      azp: T.nilable(String),
      iss: T.nilable(String)
    ).returns(T::Hash[String, T.untyped])
  end
  def self.verify_oidc(token, aud: nil, azp: nil, iss: nil)
    Google::Auth::IDTokens.verify_oidc(token, aud:, azp:, iss:)
  end
end
