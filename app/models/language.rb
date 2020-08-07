# frozen_string_literal: true

class Language < ApplicationRecord
  include StateConcern
  extend Enumerize

  enumerize :slug, in: %i[php javascript python java html css racket elixir ruby go]

  validates :slug, presence: true

  belongs_to :current_version, optional: true, class_name: 'Language::Version'

  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy
  has_many :versions, dependent: :destroy
  has_many :members, dependent: :destroy

  has_many :current_module_infos, through: :current_version, source: :module_infos
  has_many :current_lesson_infos, through: :current_version, source: :lesson_infos
  has_many :current_lesson_versions, through: :current_version, source: :lesson_versions
  has_many :current_module_versions, through: :current_version, source: :module_versions
  has_many :current_lessons, through: :current_version, source: :lessons

  delegate :to_s, to: :current_version

  aasm :state, column: :state do
    state :hidden, initial: true
    state :in_development
    state :published

    event :publish do
      transitions from: %i[hidden in_development], to: :published
    end

    event :develop do
      transitions from: [:hidden], to: :in_development
    end

    event :hide do
      transitions from: %i[developing published], to: :hidden
    end
  end
end
