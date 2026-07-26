# typed: strict
# frozen_string_literal: true

module BannerRepository
  extend ActiveSupport::Concern
  include LocaleRepository

  included do
    T.bind(self, T.class_of(Banner))

    scope :active, lambda { |locale = I18n.locale, now = Time.current|
      published_state
        .with_locale(locale)
        .where("starts_at IS NULL OR starts_at <= ?", now)
        .where("finishes_at IS NULL OR finishes_at >= ?", now)
        .order(starts_at: :desc, created_at: :desc)
    }
  end
end
