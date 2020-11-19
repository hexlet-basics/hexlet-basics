# frozen_string_literal: true

class LanguageVersionManager
  def find_or_create_language_with_version(slug)
    language = Language.find_or_create_by!(slug: slug)

    language_version = language.versions.build
    language_version.save!

    language_version
  end
end
