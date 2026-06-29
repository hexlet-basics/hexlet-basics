unless HexletBasics::Application.assets_precompiling?
  RubyLLM.configure do |config|
    config.openai_api_key = ENV.fetch("OPENAI_ACCESS_TOKEN")
    config.default_model = ENV.fetch("AI_DEFAULT_MODEL", "gpt-4.1")
    config.http_proxy = configus.hexlet_proxy.url if configus.hexlet_proxy.url.present?
    config.model_registry_class = "AiModel"
  end
end
