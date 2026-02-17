InertiaRails.configure do |config|
  # Example: force a full-reload if the deployed assets change.
  config.flash_keys = %i[notice alert error success]
  config.ssr_enabled = ViteRuby.config.ssr_build_enabled
  config.always_include_errors_hash = true
  config.use_data_inertia_head_attribute = true
  config.use_script_element_for_initial_page = true
  config.version = ViteRuby.digest
end
