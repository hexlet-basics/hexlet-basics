# typed: strict
# frozen_string_literal: true

class GoogleAuthStub
  class << self
    extend T::Sig

    sig { params(_credential: T.untyped, _options: T.untyped).returns(T::Hash[String, String]) }
    def verify_oidc(_credential, _options)
      {
        "email" => "example@mail.com",
        "sub" => "1144422255589998888"
      }
    end
  end
end
