# frozen_string_literal: true

class Language::Lesson < ApplicationRecord
  include Language::LessonRepository
  include AASM

  # FIXME: add unique index
  validates :slug, presence: true, uniqueness: { scope: :language }

  belongs_to :language
  belongs_to :module

  has_many :versions, dependent: :destroy
  has_many :members, dependent: :destroy

  has_many :infos, through: :versions, class_name: 'Language::Lesson::Version::Info'

  aasm :state do
    state :created, initial: true
    state :active
    state :archived

    event :activate do
      transitions from: %i[created archived], to: :active
    end

    event :mark_as_archived do
      transitions from: %i[created archived], to: :archived
    end
  end
end
