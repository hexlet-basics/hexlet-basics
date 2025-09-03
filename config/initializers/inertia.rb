InertiaRails.configure do |config|
  # Example: force a full-reload if the deployed assets change.
  config.ssr_enabled = ViteRuby.config.ssr_build_enabled
  config.always_include_errors_hash = true
  config.version = ViteRuby.digest
end
