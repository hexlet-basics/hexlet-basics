# typed: false
# frozen_string_literal: true

# Skeleton adapter for SMSC.ru (https://smsc.ru/api/http/).
# NOT wired in yet. To enable once the provider is chosen:
#   1. add `gem "faraday"` to the Gemfile
#   2. add a `sms` block to configus.rb (login/password on ENV)
#   3. switch the production branch in DependencyFactory to this class
#   4. uncomment the body below and bump this file to `typed: strict`
class SmsSenderSmsc < SmsSenderInterface
  ENDPOINT = "https://smsc.ru/sys/send.php"

  def self.deliver(phone:, text:)
    raise NotImplementedError, "Configure SMSC credentials and uncomment the request body"

    # Faraday.post(ENDPOINT) do |req|
    #   req.params = {
    #     login: configus.sms.login,
    #     psw: configus.sms.password,
    #     phones: phone,
    #     mes: text,
    #     fmt: 3,
    #     charset: "utf-8"
    #   }
    # end
  end
end
