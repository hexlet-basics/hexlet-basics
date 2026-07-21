# typed: false

Typelizer.configure do |config|
  # config.enable = false
  config.inheritance_strategy = :inheritance
  config.associations_strategy = :active_record
  config.enum_runtime = true
end
