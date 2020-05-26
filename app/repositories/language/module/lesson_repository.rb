# frozen_string_literal: true

module Language::Module::LessonRepository
  extend ActiveSupport::Concern

  included do
    scope :web, -> { for_locale.order(:order) }
    scope :for_locale, -> { joins(:descriptions).where(language_module_lesson_descriptions: { locale: I18n.locale }) }
  end
end
