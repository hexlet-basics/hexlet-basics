# typed: true

# `webauthn/fake_client` is an opt-in test helper loaded via an explicit
# require, so tapioca's gem RBI does not capture it. Shim the surface we use.
module WebAuthn
  class FakeClient
    def initialize(origin = T.unsafe(nil), token_binding: T.unsafe(nil), encoding: T.unsafe(nil)); end

    def create(challenge: T.unsafe(nil), rp_id: T.unsafe(nil), user_present: T.unsafe(nil), user_verified: T.unsafe(nil), **opts); end

    def get(challenge: T.unsafe(nil), rp_id: T.unsafe(nil), user_present: T.unsafe(nil), user_verified: T.unsafe(nil), sign_count: T.unsafe(nil), **opts); end
  end
end
