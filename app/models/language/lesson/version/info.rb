# frozen_string_literal: true

# == Schema Information
#
# Table name: language_lesson_version_infos
#
#  id                  :bigint           not null, primary key
#  name                :string
#  description         :string
#  locale              :string
#  theory              :string
#  tips                :string
#  definitions         :string
#  instructions        :string
#  language_id         :bigint           not null
#  language_version_id :bigint           not null
#  version_id          :bigint           not null
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#
class Language::Lesson::Version::Info < ApplicationRecord
  include Language::Lesson::Version::InfoRepository

  serialize :tips, Array
  serialize :definitions, Array

  belongs_to :language
  belongs_to :version
  belongs_to :language_version, class_name: 'Language::Version'

  def to_s
    name
  end
end
