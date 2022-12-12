# frozen_string_literal: true

# == Schema Information
#
# Table name: language_module_versions
#
#  id                  :bigint           not null, primary key
#  language_id         :bigint           not null
#  language_version_id :bigint           not null
#  module_id           :bigint           not null
#  order               :integer
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#
class Language::Module::Version < ApplicationRecord
  belongs_to :language_version, class_name: 'Language::Version'
  belongs_to :module
  belongs_to :language

  has_many :lesson_versions, dependent: :destroy,
                             foreign_key: :module_version_id,
                             class_name: 'Language::Lesson::Version',
                             inverse_of: :module_version
  has_many :infos, dependent: :destroy
end
