# frozen_string_literal: true

class GoogleAuthStub
  class << self
    def verify_oidc(_credential, _options)
      {
        "email" => "example@mail.com",
        "sub" => "1144422255589998888"
      }
    end
  end
end
