# typed: true

class ApplicationRecord < ActiveRecord::Base
  include TypedEnumSupport

  self.abstract_class = true

  # Эагерная предобработка вариантов (`preprocessed:`) гоняет нативную libvips прямо на save.
  # В test/ci/dev она не нужна (варианты не проверяются, а libvips/ruby-vips там не ставится),
  # поэтому предобрабатываем только в staging/production.
  PREPROCESS_VARIANTS = Rails.env.production? || Rails.env.staging?

  def self.enum_as_hashes(enum_name)
    send(enum_name).map { |key, value| { key: key, value: value } }
  end

  def self.form_key
    to_s.underscore.tr("/", "_")
  end
end
