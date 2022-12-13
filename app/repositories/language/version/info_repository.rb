# frozen_string_literal: true

module Language::Version::InfoRepository
  extend ActiveSupport::Concern

  included do
    scope :with_locale, ->(locale = I18n.locale) { where(locale: locale) }
    scope :completed, -> { merge(Language.with_progress(:completed)) }
    scope :incompleted, -> { merge(Language.without_progress(:completed)) }
    scope :ordered, -> { order('language.order asc') }
  end
end
