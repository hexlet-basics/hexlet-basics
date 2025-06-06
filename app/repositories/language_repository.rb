# frozen_string_literal: true

module LanguageRepository
  extend ActiveSupport::Concern

  included do
    scope :with_locale, ->(locale = I18n.locale) {
      joins(current_version: :infos)
        .merge(Language::Version::Info.with_locale(locale))
    }
    scope :web, -> { with_locale.completed_progress }
    scope :ordered, -> { order(order: :asc) }
  end
end
