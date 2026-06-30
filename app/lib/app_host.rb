# typed: strict

module AppHost
  extend T::Sig

  sig { returns(String) }
  def self.canonical
    ENV.fetch("APP_HOST", "code-basics.com")
  end

  sig { returns(T.untyped) }
  def self.port
    ENV.fetch("APP_PORT", 443)
  end

  # NOTE: for en dont use path /en but /
  sig { params(locale: T.untyped).returns(T.untyped) }
  def self.locale_for_url(locale = I18n.locale)
    return nil if locale&.to_s == I18n.default_locale.to_s

    locale
  end
end
