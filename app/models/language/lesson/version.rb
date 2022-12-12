# frozen_string_literal: true

# == Schema Information
#
# Table name: language_lesson_versions
#
#  id                  :bigint           not null, primary key
#  order               :integer
#  original_code       :string
#  prepared_code       :string
#  test_code           :string
#  natural_order       :integer
#  path_to_code        :string
#  language_version_id :bigint           not null
#  language_id         :bigint           not null
#  lesson_id           :bigint           not null
#  module_version_id   :bigint           not null
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#
class Language::Lesson::Version < ApplicationRecord
  belongs_to :language_version, class_name: 'Language::Version'
  belongs_to :lesson
  belongs_to :language
  belongs_to :module_version, class_name: 'Language::Module::Version'

  has_many :infos, dependent: :destroy

  def next_lesson
    language_version
      .lesson_versions.order(:natural_order)
      .find_by('natural_order > ?', natural_order)&.lesson
  end

  def prev_lesson
    language_version
      .lesson_versions.order(natural_order: :desc)
      .find_by('natural_order < ?', natural_order)&.lesson
  end

  def to_s
    name
  end
end
