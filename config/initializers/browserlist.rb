Browserslist.configure do |config|
  config.file_path = ".browserslist.json"

  # Enable/disable strict mode
  # When strict mode is enabled, missing browsers hash value will be set to false, which in conjunction
  # with `allow_browser`means they will be forbidden from accessing your application.
  config.strict = false
end
