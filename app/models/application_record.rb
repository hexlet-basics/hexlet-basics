# typed: strict

class ApplicationRecord < ActiveRecord::Base
  extend T::Sig
  include TypedEnumSupport

  self.abstract_class = true

  # Эагерная предобработка вариантов (`preprocessed:`) гоняет нативную libvips прямо на save.
  # В test/ci/dev она не нужна (варианты не проверяются, а libvips/ruby-vips там не ставится),
  # поэтому предобрабатываем только в staging/production.
  PREPROCESS_VARIANTS = T.let(Rails.env.production? || Rails.env.staging?, T::Boolean)

  sig { params(enum_name: T.untyped).returns(T.untyped) }
  def self.enum_as_hashes(enum_name)
    send(enum_name).map { |key, value| { key: key, value: value } }
  end

  sig { returns(String) }
  def self.form_key
    to_s.underscore.tr("/", "_")
  end
end
