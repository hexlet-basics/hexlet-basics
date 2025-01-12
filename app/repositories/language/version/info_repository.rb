# frozen_string_literal: true

module Language::Version::InfoRepository
  extend ActiveSupport::Concern

  included do
    scope :with_locale, ->(locale = I18n.locale) { where(locale: locale) }
    scope :completed, -> { merge(Language.with_progress(:completed)) }
    scope :incompleted, -> { merge(Language.with_progress(:in_development)) }
    scope :ordered, -> { order("language.order asc") }
  end
end
