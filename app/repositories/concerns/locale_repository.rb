# frozen_string_literal: true

module LocaleRepository
  extend ActiveSupport::Concern

  included do
    scope :with_locale, ->(locale = I18n.locale) { where(locale: locale) }
  end
end
