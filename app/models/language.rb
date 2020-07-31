# frozen_string_literal: true

class Language < ApplicationRecord
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

  delegate :to_s, to: :current_version

  def next_lesson_for_user(user)
    finished_members = user.finished_members_for_language(self)

    lesson_version = current_lesson_versions.where.not(lesson_id: finished_members.map(&:lesson_id)).order(:natural_order).first
    lesson_version&.lesson
  end
end
