unless HexletBasics::Application.assets_precompiling?
  OpenAI.configure do |config|
    config.access_token = ENV.fetch("OPENAI_ACCESS_TOKEN")
    # NOTE: пока закомментировал, если будем использовать то нужно будет добавить секреты
    # config.admin_token = ENV.fetch("OPENAI_ADMIN_TOKEN") # Optional, used for admin endpoints, created here: https://platform.openai.com/settings/organization/admin-keys
    # config.organization_id = ENV.fetch("OPENAI_ORGANIZATION_ID") # Optional
    config.log_errors = Rails.env.development?
  end
end
