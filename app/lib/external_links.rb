# frozen_string_literal: true

class ExternalLinks
  HEXLET_CURLS = %w[
    hexlet_rails
    hexlet_java
    hexlet_php
    hexlet_python
    hexlet_frontend
    hexlet_layout_designer
    hexlet_profession
  ]

  HEXLET_CURLS.each do |method_name|
    define_singleton_method("#{method_name}_curl") do
      ENV.fetch('BASE_DOMAIN', 'https://ru.hexlet.io') + I18n.t("links.#{method_name}")
    end
  end

  def self.hexlet_support_email
    'hola@codica.la'
  end

  def self.source_code_curl
    'https://github.com/hexlet-basics'
  end

  def self.maxim_ilyahov_curl
    'https://maximilyahov.ru/hello/'
  end

  def self.textru_curl
    'https://text.ru/'
  end

  def self.content_watch_curl
    'https://content-watch.ru/text/'
  end

  def self.glavred_curl
    'https://glvrd.ru/'
  end

  def self.telegram_community_curl
    'https://t.me/hexletcommunity'
  end

  def self.method_missing(name)
    method_name = name.to_s.delete_suffix('_curl').to_sym
    links = I18n.t('links')
    return links[method_name] if links.key?(method_name)

    super
  end

  def self.respond_to_missing?(method_name, _include_all)
    links = I18n.t('links')
    links.key?(method_name)
  end
end
