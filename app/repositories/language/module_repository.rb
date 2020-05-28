# frozen_string_literal: true

module Language::ModuleRepository
  extend ActiveSupport::Concern

  included do
    scope :web, -> { order(:order) }
  end
end
