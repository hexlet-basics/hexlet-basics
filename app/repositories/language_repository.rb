# typed: true
# frozen_string_literal: true

module LanguageRepository
  extend ActiveSupport::Concern
  extend T::Helpers
  requires_ancestor { ActiveRecord::Base }

  included do
    T.bind(self, T.class_of(Language))
    scope :with_locale, ->(locale = I18n.locale) {
      joins(current_version: :infos)
        .merge(Language::Version::Info.with_locale(locale))
    }
    scope :web, -> { with_locale.completed_progress }
    scope :ordered, -> { order(order: :asc) }
  end
end
