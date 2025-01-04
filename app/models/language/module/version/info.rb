# frozen_string_literal: true

# == Schema Information
#
# Table name: language_module_version_infos
#
#  id                  :integer          not null, primary key
#  description         :string
#  locale              :string
#  name                :string
#  created_at          :datetime         not null
#  updated_at          :datetime         not null
#  language_id         :bigint           not null
#  language_version_id :bigint           not null
#  version_id          :bigint           not null
#
# Indexes
#
#  index_language_module_version_infos_on_language_id          (language_id)
#  index_language_module_version_infos_on_language_version_id  (language_version_id)
#
# Foreign Keys
#
#  language_id          (language_id => languages.id)
#  language_version_id  (language_version_id => language_versions.id)
#  version_id           (version_id => language_module_versions.id)
#
class Language::Module::Version::Info < ApplicationRecord
  include Language::Module::Version::InfoRepository

  belongs_to :language
  belongs_to :language_version, class_name: "Language::Version"
  belongs_to :version

  def to_s
    name
  end
end
