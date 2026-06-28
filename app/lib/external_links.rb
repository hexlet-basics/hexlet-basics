# typed: strict
# frozen_string_literal: true

class ExternalLinks
  extend T::Sig

  sig { returns(String) }
  def self.hexlet_support_email
    "support@hexlet.io"
  end

  sig { returns(String) }
  def self.source_code_curl
    "https://github.com/hexlet-basics"
  end

  sig { params(path_to_code: T.untyped, locale: T.untyped).returns(String) }
  def self.lesson_source_code_curl(path_to_code, locale)
    repository_path = ExternalLinks.source_code_curl
    path_to_description = File.join(repository_path, path_to_code, locale, "README.md")

    path_to_description.sub("modules", "blob/main/modules")
  end


  sig { returns(String) }
  def self.maxim_ilyahov_curl
    "https://maximilyahov.ru/hello/"
  end

  sig { returns(String) }
  def self.textru_curl
    "https://text.ru/"
  end

  sig { returns(String) }
  def self.content_watch_curl
    "https://content-watch.ru/text/"
  end

  sig { returns(String) }
  def self.glavred_curl
    "https://glvrd.ru/"
  end

  sig { returns(String) }
  def self.telegram_community_curl
    I18n.t("common.community_url")
  end

  # def self.method_missing(name)
  #   return 'FIXME'
  #   # method_name = name.to_s.delete_suffix("_curl").to_sym
  #   # links = I18n.t("links")
  #   # return links[method_name] if links.key?(method_name)
  #
  #   super
  # end

  sig { params(method_name: T.untyped, _include_all: T.untyped).returns(T.untyped) }
  def self.respond_to_missing?(method_name, _include_all)
    links = I18n.t("links")
    links.key?(method_name)
  end
end
