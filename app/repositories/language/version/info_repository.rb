# frozen_string_literal: true

module Language::Version::InfoRepository
  extend ActiveSupport::Concern

  included do
    scope :with_locale, ->(locale = I18n.locale) { where(locale: locale) }
    scope :completed, -> { merge(Language.with_progress(:completed)) }
    scope :incompleted, -> { merge(Language.with_progress(:in_development)) }
    scope :ordered, -> { order("language.order asc") }
    scope :current, -> {
      joins(:language)
      .merge(Language.with_progress(:completed))
      .where(Language.arel_table[:current_version_id].eq(Language::Version::Info.arel_table[:language_version_id]))
    }
  end
end
