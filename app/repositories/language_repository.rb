# frozen_string_literal: true

module LanguageRepository
  extend ActiveSupport::Concern

  included do
    scope :ordered, -> { order(order: :asc) }
  end
end
