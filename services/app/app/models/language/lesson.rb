# frozen_string_literal: true

class Language::Lesson < ApplicationRecord
  include Language::LessonRepository

  # FIXME: add unique index
  validates :slug, presence: true, uniqueness: { scope: :language }

  belongs_to :language
  belongs_to :module

  has_many :versions, dependent: :destroy
  has_many :members, dependent: :destroy

  has_many :infos, through: :versions, class_name: 'Language::Lesson::Version::Info'
end
