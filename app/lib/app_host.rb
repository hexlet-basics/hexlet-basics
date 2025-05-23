module AppHost
  def self.canonical
    ENV.fetch("APP_HOST", "code-basics.com")
  end

  def self.port
    ENV.fetch("APP_PORT", 443)
  end

  # NOTE: for en dont use path /en but /
  def self.locale_for_url(locale = I18n.locale)
    return nil if locale&.to_sym == I18n.default_locale

    locale
  end
end
