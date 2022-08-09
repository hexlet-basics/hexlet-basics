# frozen_string_literal: true

module AppHost
  def self.subdomain
    subdomain_for(I18n.locale)
  end

  def self.subdomain_for(locale)
    locale == :ru ? 'ru' : ''
  end

  def self.canonical
    ENV.fetch('APP_HOST')
  end
end
