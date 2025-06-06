# frozen_string_literal: true

module Language::Version::InfoRepository
  extend ActiveSupport::Concern
  include LocaleRepository

  included do
    scope :completed, -> { merge(Language.completed_progress) }
    scope :incompleted, -> { merge(Language.in_development_progress) }
    scope :ordered, -> { order("language.order asc") }
    scope :current, -> {
      joins(:language)
      .merge(Language.completed_progress)
      .where(Language.arel_table[:current_version_id].eq(Language::Version::Info.arel_table[:language_version_id]))
    }
  end
end
