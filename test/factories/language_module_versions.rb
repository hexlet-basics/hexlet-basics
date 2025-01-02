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
FactoryBot.define do
  factory :language_module_version do
    language_version { nil }
    order { 'MyString' }
  end
end
