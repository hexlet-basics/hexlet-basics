# frozen_string_literal: true

# == Schema Information
#
# Table name: language_lesson_version_infos
#
#  id                  :bigint           not null, primary key
#  definitions         :string
#  description         :string
#  instructions        :string
#  locale              :string
#  name                :string
#  theory              :string
#  tips                :string
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#  language_id         :bigint           not null
#  language_version_id :bigint           not null
#  version_id          :bigint           not null
#
# Indexes
#
#  index_language_lesson_version_infos_on_language_id          (language_id)
#  index_language_lesson_version_infos_on_language_version_id  (language_version_id)
#
# Foreign Keys
#
#  fk_rails_...  (language_id => languages.id)
#  fk_rails_...  (language_version_id => language_versions.id)
#  fk_rails_...  (version_id => language_lesson_versions.id)
#
class Language::Lesson::Version::Info < ApplicationRecord
  include Language::Lesson::Version::InfoRepository

  serialize :tips, type: Array
  serialize :definitions, type: Array

  belongs_to :language
  belongs_to :version
  belongs_to :language_version, class_name: "Language::Version"

  def to_s
    name
  end
end
