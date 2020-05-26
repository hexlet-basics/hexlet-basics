# frozen_string_literal: true

module Language::ModuleRepository
  extend ActiveSupport::Concern

  included do
    scope :web, -> { for_locale.order(:order) }
    scope :for_locale, -> { joins(:descriptions).where(language_module_descriptions: { locale: I18n.locale }) }
  end
end
