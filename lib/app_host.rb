# frozen_string_literal: true

module AppHost
  def self.subdomain
    I18n.locale == :ru ? 'ru' : ''
  end

  def self.canonical
    ENV.fetch('APP_HOST')
  end
end
