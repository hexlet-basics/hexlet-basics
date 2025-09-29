require "omniauth/strategies/oauth2"

module OmniAuth
  module Strategies
    class Vk < OmniAuth::Strategies::OAuth2
      option :name, "vk"

      option :client_options,
             site: "https://id.vk.ru/",
             token_url: "/oauth2/auth",
             auth_scheme: :request_body

      option :pkce, true
      option :provider_ignores_state, true

      def token_params
        super.tap do |params|
          %w[code code_verifier device_id].each do |v|
            if request.params[v]
              params[v.to_sym] = request.params[v]
            end
          end
        end
      end

      def callback_url
        options.callback_url || full_host + callback_path
      end

      uid { raw_info["user"]["user_id"].to_s }

      info do
        { email: raw_info["user"]["email"] }
      end

      extra do
        { raw_info: raw_info }
      end

      def raw_info
        @raw_info ||= access_token.get(
          "/oauth2/user_info",
          params: { client_id: options.client_id }
        ).parsed
      end
    end
  end
end
