# typed: strict
# frozen_string_literal: true

module LocaleRepository
  extend ActiveSupport::Concern
  extend T::Helpers
  requires_ancestor { ActiveRecord::Base }

  included do
    T.bind(self, T.class_of(ActiveRecord::Base))
    scope :with_locale, ->(locale = I18n.locale) { where(locale: locale) }
  end
end
