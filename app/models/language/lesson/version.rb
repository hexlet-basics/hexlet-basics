# frozen_string_literal: true

# == Schema Information
#
# Table name: language_lesson_versions
#
#  id                  :integer          not null, primary key
#  natural_order       :integer
#  order               :integer
#  original_code       :string
#  path_to_code        :string
#  prepared_code       :string
#  test_code           :string
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#  language_id         :bigint           not null
#  language_version_id :bigint           not null
#  lesson_id           :bigint           not null
#  module_version_id   :bigint           not null
#
# Indexes
#
#  index_language_lesson_versions_on_language_id          (language_id)
#  index_language_lesson_versions_on_language_version_id  (language_version_id)
#  index_language_lesson_versions_on_lesson_id            (lesson_id)
#  index_language_lesson_versions_on_module_version_id    (module_version_id)
#
# Foreign Keys
#
#  language_id          (language_id => languages.id)
#  language_version_id  (language_version_id => language_versions.id)
#  lesson_id            (lesson_id => language_lessons.id)
#  module_version_id    (module_version_id => language_module_versions.id)
#
class Language::Lesson::Version < ApplicationRecord
  belongs_to :language_version, class_name: "Language::Version"
  belongs_to :lesson
  belongs_to :language
  belongs_to :module_version, class_name: "Language::Module::Version"
  has_one :module, through: :module_version

  has_many :infos, dependent: :destroy

  validates :natural_order, presence: true

  def next_lesson_version
    language_version
      .lesson_versions.order(:natural_order)
      .joins(:infos)
      .merge(Language::Lesson::Version::Info.with_locale)
      .find_by("natural_order > ?", natural_order)&.lesson
  end

  def prev_lesson_version
    language_version
      .lesson_versions.order(natural_order: :desc)
      .joins(:infos)
      .merge(Language::Lesson::Version::Info.with_locale)
      .find_by("natural_order < ?", natural_order)&.lesson
  end

  def to_s
    name
  end
end
