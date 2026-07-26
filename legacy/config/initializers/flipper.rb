# typed: false

# Feature flags backed by Active Record (flipper_features / flipper_gates).
# The adapter is wired automatically by flipper-active_record; this file only
# tunes the admin web UI mounted at /admin/flipper (see config/routes.rb).
Flipper::UI.configure do |config|
  config.banner_text = "#{Rails.env} environment"
  config.banner_class = Rails.env.production? ? "danger" : "info"
  config.feature_creation_enabled = true
  config.feature_removal_enabled = true
  config.show_feature_description_in_list = true
end
