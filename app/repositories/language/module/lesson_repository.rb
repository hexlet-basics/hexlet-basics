# frozen_string_literal: true

module Language::Module::LessonRepository
  extend ActiveSupport::Concern

  included do
    scope :web, -> { order(:order) }
  end
end
