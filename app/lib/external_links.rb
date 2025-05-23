# frozen_string_literal: true

class ExternalLinks
  def self.hexlet_support_email
    "support@hexlet.io"
  end

  def self.source_code_curl
    "https://github.com/hexlet-basics"
  end

  def self.lesson_source_code_curl(lesson_version_info)
    repository_path = ExternalLinks.source_code_curl
    locale = lesson_version_info.locale
    path_to_description = File.join(repository_path, lesson_version_info.version.path_to_code, locale, "README.md")

    path_to_description.sub("modules", "blob/main/modules")
  end


  def self.maxim_ilyahov_curl
    "https://maximilyahov.ru/hello/"
  end

  def self.textru_curl
    "https://text.ru/"
  end

  def self.content_watch_curl
    "https://content-watch.ru/text/"
  end

  def self.glavred_curl
    "https://glvrd.ru/"
  end

  def self.telegram_community_curl
    "https://t.me/HexletLearningBot"
  end

  # def self.method_missing(name)
  #   return 'FIXME'
  #   # method_name = name.to_s.delete_suffix("_curl").to_sym
  #   # links = I18n.t("links")
  #   # return links[method_name] if links.key?(method_name)
  #
  #   super
  # end

  def self.respond_to_missing?(method_name, _include_all)
    links = I18n.t("links")
    links.key?(method_name)
  end
end
