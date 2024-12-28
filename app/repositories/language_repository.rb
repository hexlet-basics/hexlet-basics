# frozen_string_literal: true

module LanguageRepository
  extend ActiveSupport::Concern

  included do
    scope :with_locale, lambda { |locale = I18n.locale|
      where('exists (?)', Language::Version
                            .where('language_versions.id = languages.current_version_id')
                            .joins(:infos)
                            .merge(Language::Version::Info.with_locale(locale))
                            .select('1')
                            .limit(1))
    }
    scope :web, -> { with_locale.with_progress(:completed) }
    scope :ordered, -> { order(order: :asc) }
  end
end
