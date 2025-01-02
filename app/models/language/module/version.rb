# frozen_string_literal: true

# == Schema Information
#
# Table name: language_module_versions
#
#  id                  :integer          not null, primary key
#  order               :integer
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#  language_id         :bigint           not null
#  language_version_id :bigint           not null
#  module_id           :bigint           not null
#
# Indexes
#
#  index_language_module_versions_on_language_id          (language_id)
#  index_language_module_versions_on_language_version_id  (language_version_id)
#  index_language_module_versions_on_module_id            (module_id)
#
# Foreign Keys
#
#  language_id          (language_id => languages.id)
#  language_version_id  (language_version_id => language_versions.id)
#  module_id            (module_id => language_modules.id)
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
