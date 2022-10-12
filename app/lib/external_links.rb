# frozen_string_literal: true

class ExternalLinks
  def self.hexlet_support_email
    'support@hexlet.io'
  end

  def self.source_code_curl
    'https://github.com/hexlet-basics'
  end

  def self.language_source_code_curl(language)
    "#{source_code_curl}/exercises-#{language}"
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
