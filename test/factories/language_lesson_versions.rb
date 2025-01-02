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
FactoryBot.define do
  factory :language_lesson_version do
    # TODO: implement later
  end
end
