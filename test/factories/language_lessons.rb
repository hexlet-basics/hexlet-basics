# frozen_string_literal: true

# == Schema Information
#
# Table name: language_lessons
#
#  id            :integer          not null, primary key
#  natural_order :integer
#  order         :integer
#  original_code :text
#  path_to_code  :string(255)
#  prepared_code :text
#  slug          :string(255)
#  state         :string(255)
#  test_code     :text
#  created_at    :datetime         not null
#  updated_at    :datetime         not null
#  language_id   :bigint
#  module_id     :bigint
#  upload_id     :bigint
#
# Indexes
#
#  index_language_lessons_on_language_id_and_slug  (language_id,slug) UNIQUE
#  language_module_lessons_language_id_index       (language_id)
#  language_module_lessons_module_id_index         (module_id)
#  language_module_lessons_upload_id_index         (upload_id)
#
# Foreign Keys
#
#  language_id  (language_id => languages.id)
#  module_id    (module_id => language_modules.id)
#  upload_id    (upload_id => uploads.id)
#
FactoryBot.define do
  factory :language_lesson do
    slug { "MyString" }
    language { nil }
    language_module { nil }
    current_version { nil }
  end
end
